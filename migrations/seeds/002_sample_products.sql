-- ============================================================
-- MagTrade 種子資料 - 示範商品
-- 類型: Seed（測試資料）
-- 環境: 僅開發環境
-- 說明: 提供基礎商品用於功能測試
-- ============================================================

-- 基礎測試商品（Apple 系列）
INSERT INTO products (name, description, original_price, image_url, status) VALUES
('iPhone 15 Pro Max 256GB', 
 'Apple iPhone 15 Pro Max，256GB，鈦金屬設計', 
 9999.00, 'https://example.com/iphone15.jpg', 1),

('MacBook Pro 14 M3 Pro', 
 'Apple MacBook Pro 14 吋，M3 Pro 晶片，18GB RAM，512GB SSD', 
 16999.00, 'https://example.com/macbook.jpg', 1),

('AirPods Pro 2', 
 'Apple AirPods Pro 第二代，主動降噪', 
 1899.00, 'https://example.com/airpods.jpg', 1)

ON CONFLICT DO NOTHING;
