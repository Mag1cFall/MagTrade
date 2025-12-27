-- ============================================================
-- MagTrade 資料庫結構更新 - 審計日誌
-- 版本: 002
-- 建立日期: 2025-12
-- 說明: 新增審計日誌表，記錄使用者關鍵操作
-- ============================================================

-- ------------------------------------------------------------
-- 審計日誌表
-- 記錄使用者登入、下單、修改等關鍵操作軌跡
-- 用於安全審計和問題追蹤
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),                         -- 操作者（可為空，如匿名請求）
    action VARCHAR(50) NOT NULL,                                 -- 操作類型：login/logout/order/update
    resource VARCHAR(100),                                       -- 資源類型：user/order/flash_sale
    resource_id VARCHAR(50),                                     -- 資源 ID
    ip VARCHAR(45),                                              -- 來源 IP（支援 IPv6）
    user_agent VARCHAR(255),                                     -- 瀏覽器標識
    details TEXT,                                                -- 額外詳情（JSON 格式）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP                                         -- 軟刪除標記
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_deleted_at ON audit_logs(deleted_at);
