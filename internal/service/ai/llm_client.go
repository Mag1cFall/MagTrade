// LLM 大語言模型客戶端
//
// 本檔案封裝 OpenAI 兼容 API（如 SiliconFlow）的調用
// 支援普通對話和串流對話（SSE）
// 主要用於智慧客服和策略分析
package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"go.uber.org/zap"
)

// LLMClient 大語言模型客戶端
type LLMClient struct {
	cfg        *config.AIConfig
	httpClient *http.Client
	log        *zap.Logger
}

func NewLLMClient(cfg *config.AIConfig, log *zap.Logger) *LLMClient {
	return &LLMClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // LLM 回應可能較慢
		},
		log: log,
	}
}

// ChatMessage 對話消息
type ChatMessage struct {
	Role             string `json:"role"` // system/user/assistant
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"` // 思考過程（部分模型支援）
}

// ChatCompletionRequest OpenAI 格式請求
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream"`
}

// ChatCompletionResponse OpenAI 格式回應
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role             string `json:"role"`
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content,omitempty"`
		} `json:"message"`
		Delta struct { // 串流時的增量內容
			Role             string `json:"role"`
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content,omitempty"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// StreamChunk 串流對話的單個片段
type StreamChunk struct {
	Type    string `json:"type"` // thinking/content/done/error
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// Chat 普通對話（等待完整回應）
func (c *LLMClient) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	reqBody := ChatCompletionRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		MaxTokens:   c.cfg.MaxTokens,
		Temperature: c.cfg.Temperature,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.cfg.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.log.Error("LLM API error",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)),
		)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	c.log.Debug("LLM response",
		zap.Int("prompt_tokens", chatResp.Usage.PromptTokens),
		zap.Int("completion_tokens", chatResp.Usage.CompletionTokens),
	)

	return chatResp.Choices[0].Message.Content, nil
}

// ChatStream 串流對話（逐字回傳）
func (c *LLMClient) ChatStream(ctx context.Context, messages []ChatMessage, chunkHandler func(chunk StreamChunk) error) error {
	reqBody := ChatCompletionRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		MaxTokens:   c.cfg.MaxTokens,
		Temperature: c.cfg.Temperature,
		Stream:      true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.cfg.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// 解析 SSE 串流
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" { // 串流結束
			return chunkHandler(StreamChunk{Type: "done", Done: true})
		}

		var chatResp ChatCompletionResponse
		if err := json.Unmarshal([]byte(data), &chatResp); err != nil {
			continue
		}

		if len(chatResp.Choices) == 0 {
			continue
		}

		delta := chatResp.Choices[0].Delta
		if delta.ReasoningContent != "" { // 思考過程
			if err := chunkHandler(StreamChunk{Type: "thinking", Content: delta.ReasoningContent}); err != nil {
				return err
			}
		}
		if delta.Content != "" { // 正式內容
			if err := chunkHandler(StreamChunk{Type: "content", Content: delta.Content}); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

// ChatWithSystem 帶系統提示詞的對話
func (c *LLMClient) ChatWithSystem(ctx context.Context, systemPrompt string, userMessage string) (string, error) {
	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMessage},
	}
	return c.Chat(ctx, messages)
}

// ChatWithHistory 帶歷史記錄的對話
func (c *LLMClient) ChatWithHistory(ctx context.Context, systemPrompt string, history []ChatMessage, userMessage string) (string, error) {
	messages := make([]ChatMessage, 0, len(history)+2)
	messages = append(messages, ChatMessage{Role: "system", Content: systemPrompt})
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{Role: "user", Content: userMessage})
	return c.Chat(ctx, messages)
}

// ChatStreamWithHistory 帶歷史記錄的串流對話
func (c *LLMClient) ChatStreamWithHistory(ctx context.Context, systemPrompt string, history []ChatMessage, userMessage string, chunkHandler func(chunk StreamChunk) error) error {
	messages := make([]ChatMessage, 0, len(history)+2)
	messages = append(messages, ChatMessage{Role: "system", Content: systemPrompt})
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{Role: "user", Content: userMessage})
	return c.ChatStream(ctx, messages, chunkHandler)
}
