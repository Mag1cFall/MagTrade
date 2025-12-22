package cache

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

const DistributedLockScript = `
local lock_key = KEYS[1]
local lock_value = ARGV[1]
local expire_time = tonumber(ARGV[2])

if redis.call('SET', lock_key, lock_value, 'NX', 'PX', expire_time) then
    return 1
end
return 0
`

const DistributedUnlockScript = `
local lock_key = KEYS[1]
local lock_value = ARGV[1]

if redis.call('GET', lock_key) == lock_value then
    return redis.call('DEL', lock_key)
end
return 0
`
