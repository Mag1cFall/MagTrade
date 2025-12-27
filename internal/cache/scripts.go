// Redis Lua 腳本定義
//
// 本檔案包含秒殺系統核心的 Lua 腳本
// Lua 腳本在 Redis 中原子執行，保證併發安全
// 主要用於：庫存扣減、庫存恢復、分散式鎖
package cache

// DeductStockScript 庫存扣減腳本
// KEYS[1]: 庫存 Key (flash:stock:{id})
// KEYS[2]: 使用者已購數量 Key (flash:bought:{id}:{uid})
// ARGV[1]: 購買數量
// ARGV[2]: 限購數量
// 返回值: [1,"success"] 成功 / [-1,"库存不足"] / [-2,"超出限购数量"]
const DeductStockScript = `
local stock_key = KEYS[1]
local bought_key = KEYS[2]
local quantity = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])

local stock = tonumber(redis.call('GET', stock_key) or 0)
local user_bought = tonumber(redis.call('GET', bought_key) or 0)

if stock < quantity then
    return {-1, "库存不足"}
end

if user_bought + quantity > limit then
    return {-2, "超出限购数量"}
end

redis.call('DECRBY', stock_key, quantity)
redis.call('INCRBY', bought_key, quantity)
redis.call('EXPIRE', bought_key, 86400)

return {1, "success"}
`

// RestoreStockScript 庫存恢復腳本（訂單取消或 Kafka 發送失敗時回滾）
// KEYS[1]: 庫存 Key
// KEYS[2]: 使用者已購數量 Key
// ARGV[1]: 恢復數量
const RestoreStockScript = `
local stock_key = KEYS[1]
local bought_key = KEYS[2]
local quantity = tonumber(ARGV[1])

redis.call('INCRBY', stock_key, quantity)
local bought = tonumber(redis.call('GET', bought_key) or 0)
if bought >= quantity then
    redis.call('DECRBY', bought_key, quantity)
end

return {1, "success"}
`

// DistributedLockScript 分散式鎖加鎖腳本
// 使用 SET NX PX 命令，若 Key 不存在則設定並返回 1，否則返回 0
// KEYS[1]: 鎖 Key (flash:lock:{id}:{uid})
// ARGV[1]: 鎖值（UUID，解鎖時驗證身份）
// ARGV[2]: 過期時間（毫秒）
const DistributedLockScript = `
local lock_key = KEYS[1]
local lock_value = ARGV[1]
local expire_time = tonumber(ARGV[2])

if redis.call('SET', lock_key, lock_value, 'NX', 'PX', expire_time) then
    return 1
end
return 0
`

// DistributedUnlockScript 分散式鎖解鎖腳本
// 只有鎖的持有者（value 匹配）才能解鎖，防止誤刪其他人的鎖
// KEYS[1]: 鎖 Key
// ARGV[1]: 鎖值
const DistributedUnlockScript = `
local lock_key = KEYS[1]
local lock_value = ARGV[1]

if redis.call('GET', lock_key) == lock_value then
    return redis.call('DEL', lock_key)
end
return 0
`
