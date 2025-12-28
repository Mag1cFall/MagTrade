-- ============================================================
-- MagTrade 種子資料 - 管理員帳號
-- 類型: Seed（初始資料）
-- 環境: 開發環境，僅供演示
-- 說明: 建立預設管理員帳號，生產環境需修改密碼
-- 屬冗餘文件，已被postgre.go功能覆蓋
-- ============================================================

-- 預設管理員帳號
-- 密碼: mag1cfall1337 (bcrypt 雜湊)
-- ⚠️ 生產環境部署後請立即修改密碼！
INSERT INTO users (username, email, password_hash, role, status, email_verified)
VALUES ('admin', 'admin@magtrade.com', 
        '$2a$12$auU8/cmx4aJ8mzl47q6ZGelwE8tYIwkG5PdgqK6dLE81k/jSA.k7i', 
        'admin', 1, true)
ON CONFLICT (username) DO NOTHING;  -- 已存在則跳過，避免重複執行報錯
