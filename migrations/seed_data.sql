-- Seed data for MagTrade
-- Run this after GORM has created the tables

-- Insert admin user (password: admin123)
INSERT INTO users (username, email, password_hash, role, status, email_verified) VALUES
('admin', 'admin@magtrade.com', '$2a$10$bKmO1qYASj5loeB5oT4nBO2NXISkf2E7vaOPQQ5/pBXy0FSQpQO0m', 'admin', 1, true)
ON CONFLICT (username) DO NOTHING;

-- Insert sample products
INSERT INTO products (name, description, original_price, image_url, status) VALUES
('iPhone 15 Pro Max 256GB', 'Apple iPhone 15 Pro Max, 256GB, Titanium Design', 9999.00, 'https://example.com/iphone15.jpg', 1),
('MacBook Pro 14 M3 Pro', 'Apple MacBook Pro 14, M3 Pro chip, 18GB RAM, 512GB SSD', 16999.00, 'https://example.com/macbook.jpg', 1),
('AirPods Pro 2', 'Apple AirPods Pro 2nd Gen, Active Noise Cancellation', 1899.00, 'https://example.com/airpods.jpg', 1)
ON CONFLICT DO NOTHING;
