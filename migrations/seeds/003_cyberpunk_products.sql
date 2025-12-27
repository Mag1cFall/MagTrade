-- ============================================================
-- MagTrade 種子資料 - 賽博龐克商品（擴展）
-- 類型: Seed（展示資料）
-- 環境: 僅開發環境
-- 說明: 50件虛構科技商品，用於壓測和 UI 展示
-- ============================================================

INSERT INTO products (name, description, original_price, image_url, status) VALUES
-- 神經介面類
('Neuro-Link Pro v4', '神經介面設備，超低延遲直連資料串流', 12999.00, 'local:phone', 1),
('Neural-Imprint Kit', '數位意識備份與還原套件', 35000.00, 'local:box', 1),
('Neural-Stress Monitor', '認知負載監測器，追蹤高壓事件壓力', 1600.00, 'local:phone', 1),
('Neural-Link Cable', '舊式硬體直連介面線材', 180.00, 'local:box', 1),

-- 義肢強化類
('Cybernetic Arm Mk.II', '工業級義肢，內建電漿切割器', 45000.00, 'local:box', 1),
('Cyber-Eye Mk.III', '光學義眼，100倍變焦與夜視功能', 9500.00, 'local:phone', 1),
('High-Torque Cyber-Hand', '高扭力義手，針對體力勞動優化', 18000.00, 'local:box', 1),

-- 運算設備類
('Onyx Workstation', '量子冷卻桌機，大規模平行模擬運算', 89999.00, 'local:laptop', 1),
('Viper GPU X-100', 'Teraflop 級圖形加速器，即時渲染', 12000.00, 'local:laptop', 1),
('Glitch-Core CPU', '非線性邏輯執行優化處理器', 7800.00, 'local:laptop', 1),
('Matrix-Key Deck', '全息輸入介面，快速程式碼部署', 2100.00, 'local:laptop', 1),

-- 安全防護類
('Thermal Sight Visor', '戰術抬頭顯示器，多頻譜分析', 3499.00, 'local:audio', 1),
('Stealth Cloak Skin', '光線折射穿戴織物，視覺隱匿', 5600.00, 'local:audio', 1),
('Shadow-Step Boots', '消音鞋履，靜音行動', 2400.00, 'local:box', 1),
('Carbon-Fiber Mask', '輕量呼吸器，危險環境適用', 1100.00, 'local:box', 1),
('Digital-Hazmat Suit', '全身防護裝，抵禦輻射與邏輯炸彈', 15000.00, 'local:box', 1),

-- 資料儲存類
('Data-Shard 2TB', '高速實體儲存，加密協定日誌', 899.00, 'local:box', 1),
('Encrypted Journal', '化學自毀墨水實體筆記本', 320.00, 'local:box', 1),
('Smart-Contract Hub', '安全協定執行硬體錢包', 1200.00, 'local:box', 1),

-- 通訊設備類
('Encrypted Comms Unit', '直連衛星上鏈，4096位元加密', 5500.00, 'local:phone', 1),
('Signal Jammer Mini', '口袋型射頻干擾器，隱私保護', 750.00, 'local:phone', 1),
('Encrypted-Voice Mic', '硬體級語音遮罩，清晰通訊', 950.00, 'local:audio', 1),

-- 網路設備類
('Edge-Node Router', '分散式網路中樞，零延遲交易', 4500.00, 'local:box', 1),
('Hydra Server Rack', '多租戶高可用運算叢集', 120000.00, 'local:box', 1),
('Sub-Grid Mapper', '暗網節點即時拓撲掃描器', 6700.00, 'local:laptop', 1),
('Satellite-Dish Pro', '高增益天線，遠端上鏈', 8500.00, 'local:box', 1),
('Fiber-Optic Cable 10km', '超高速實體連線，資料中心用', 5000.00, 'local:box', 1),

-- 電源供應類
('Titan Power-Cell', '高密度電池，義體增強供電', 1200.00, 'local:box', 1),
('Mini-Fusion Cell', '長效電源，攜帶式主機用', 4500.00, 'local:box', 1),

-- 音訊設備類
('Void-Ear Headset', '零洩漏聽覺沉浸，神經隔離', 3200.00, 'local:audio', 1),
('Sonic-Pulse Audio', '聲波防禦與分析工具', 4200.00, 'local:audio', 1),

-- 駭客工具類
('Protocol-Z Bypass', '標準安全檢查點軟體繞過', 1500.00, 'local:box', 1),
('Tactical Glove v2', '觸覺反饋手套，內建 NFC 駭入工具', 1800.00, 'local:box', 1),
('Bio-Metric Scrambler', '偽造生理特徵訊號裝置', 2800.00, 'local:phone', 1),
('Drone-Controller Kit', '自主偵察機遠程控制上鏈', 3200.00, 'local:laptop', 1),

-- 感測器類
('Bio-Monitor Patch', '皮下感測器，即時健康資料上傳', 299.00, 'local:phone', 1),
('Augmented Reality Glasses', 'AR 眼鏡，疊加數位資訊於現實', 4500.00, 'local:phone', 1),

-- 維修工具類
('Nano-Repair Paste', '即時硬體修復，輕微外殼損傷', 350.00, 'local:box', 1),
('Micro-Soldering Pen', '精密焊接工具，電路維修', 250.00, 'local:box', 1),
('Liquid-Metal Paste', '超高導熱率散熱介質', 150.00, 'local:box', 1),
('Static-Guard Spray', '防靜電噴霧，乾燥環境電子保護', 80.00, 'local:box', 1),

-- 加密安全類
('Quantum-Crypto Key', '一次性金鑰，不可破解資料傳輸', 999.00, 'local:box', 1),
('Protocol-Zero Chip', '零知識證明專用矽晶片', 8900.00, 'local:phone', 1),

-- 特殊裝備類
('Plasma Blade XP', '濃縮能量刀刃，快速資產清算', 22000.00, 'local:box', 1),
('Pulse-Grenade (Mock)', '非致命 EMP 模擬器，系統測試用', 450.00, 'local:box', 1),
('Emergency Kill-Switch', '危急系統故障實體中斷器', 500.00, 'local:box', 1),
('Radar-Absorbent Paint', '電子訊號特徵降低塗料', 2200.00, 'local:box', 1),
('Heat-Shield Blanket', '敏感硬體火災防護毯', 600.00, 'local:box', 1),

-- 交通載具類
('Neon-Cycle Frame', '超輕碳纖維車架，高速移動', 15000.00, 'local:box', 1),

-- 可程式硬體類
('Universal Logic Gate', '可程式閘陣列，自訂硬體', 1400.00, 'local:box', 1);
