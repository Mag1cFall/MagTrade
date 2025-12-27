-- ============================================================
-- MagTrade 資料庫結構初始化
-- 版本: 001
-- 建立日期: 2025-12
-- 說明: 建立核心業務表結構，不含初始資料
-- ============================================================

-- 建立資料庫（Docker 環境由 compose 自動建立，此處為手動備份用）
-- CREATE DATABASE magtrade;
-- \c magtrade;

-- ------------------------------------------------------------
-- 使用者表
-- 儲存系統使用者資訊，支援軟刪除
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,                                    -- 主鍵，自動遞增
    username VARCHAR(50) UNIQUE NOT NULL,                        -- 使用者名稱，唯一
    email VARCHAR(100) UNIQUE NOT NULL,                          -- 電子郵件，唯一
    password_hash VARCHAR(255) NOT NULL,                         -- 密碼雜湊值（bcrypt）
    role VARCHAR(20) DEFAULT 'user',                             -- 角色：user/admin
    status SMALLINT DEFAULT 1,                                   -- 狀態：0=停用, 1=啟用
    email_verified BOOLEAN DEFAULT false,                        -- 郵件是否已驗證
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP                                         -- 軟刪除標記
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- ------------------------------------------------------------
-- 商品表
-- 儲存可上架銷售的商品資訊
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,                                  -- 商品名稱
    description TEXT,                                            -- 商品描述
    original_price DECIMAL(10,2) NOT NULL,                       -- 原價（單位：元）
    image_url VARCHAR(500),                                      -- 商品圖片 URL
    status SMALLINT DEFAULT 1,                                   -- 狀態：0=下架, 1=上架
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_deleted_at ON products(deleted_at);

-- ------------------------------------------------------------
-- 秒殺活動表
-- 儲存限時搶購活動資訊
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS flash_sales (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id),          -- 關聯商品
    flash_price DECIMAL(10,2) NOT NULL,                          -- 秒殺價格
    total_stock INT NOT NULL,                                    -- 總庫存數量
    available_stock INT NOT NULL,                                -- 剩餘庫存（Redis 同步）
    per_user_limit INT DEFAULT 1,                                -- 每人限購數量
    start_time TIMESTAMP NOT NULL,                               -- 活動開始時間
    end_time TIMESTAMP NOT NULL,                                 -- 活動結束時間
    status SMALLINT DEFAULT 0,                                   -- 狀態：0=待開始, 1=進行中, 2=已結束
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_flash_sales_product_id ON flash_sales(product_id);
CREATE INDEX idx_flash_sales_start_time ON flash_sales(start_time);
CREATE INDEX idx_flash_sales_end_time ON flash_sales(end_time);
CREATE INDEX idx_flash_sales_status ON flash_sales(status);
CREATE INDEX idx_flash_sales_deleted_at ON flash_sales(deleted_at);

-- ------------------------------------------------------------
-- 訂單表
-- 儲存使用者下單記錄
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(32) UNIQUE NOT NULL,                        -- 訂單編號（雪花演算法生成）
    user_id BIGINT NOT NULL REFERENCES users(id),                -- 下單使用者
    flash_sale_id BIGINT NOT NULL REFERENCES flash_sales(id),    -- 關聯秒殺活動
    amount DECIMAL(10,2) NOT NULL,                               -- 訂單金額
    quantity INT DEFAULT 1,                                      -- 購買數量
    status SMALLINT DEFAULT 0,                                   -- 狀態：0=待付款, 1=已付款, 2=已取消
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMP,                                           -- 付款時間
    deleted_at TIMESTAMP
);

CREATE INDEX idx_orders_order_no ON orders(order_no);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_flash_sale_id ON orders(flash_sale_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_deleted_at ON orders(deleted_at);

-- ------------------------------------------------------------
-- AI 對話歷史表
-- 儲存使用者與智慧客服的對話記錄
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS chat_histories (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    session_id VARCHAR(64) NOT NULL,                             -- 對話 Session ID
    role VARCHAR(20) NOT NULL,                                   -- 角色：user/assistant
    content TEXT NOT NULL,                                       -- 訊息內容
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_chat_histories_user_session ON chat_histories(user_id, session_id);

-- ------------------------------------------------------------
-- AI 推薦記錄表
-- 儲存 AI 針對秒殺活動的分析建議
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ai_recommendations (
    id BIGSERIAL PRIMARY KEY,
    flash_sale_id BIGINT NOT NULL REFERENCES flash_sales(id),
    recommendation_type VARCHAR(50) NOT NULL,                    -- 類型：strategy/risk/timing
    content JSONB NOT NULL,                                      -- JSON 格式建議內容
    confidence_score FLOAT,                                      -- 信心分數 0-1
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ai_recommendations_flash_sale_id ON ai_recommendations(flash_sale_id);
CREATE INDEX idx_ai_recommendations_type ON ai_recommendations(recommendation_type);
