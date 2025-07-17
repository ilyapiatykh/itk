-- wrk -t4 -c40 -d10s -s two_wallets.lua --latency http://localhost:80

local baseUrl = "/api/v1"
local wallet1 = "11111111-1111-1111-1111-111111111111"
local wallet2 = "22222222-2222-2222-2222-222222222222"
local depositAmount = 1000
local withdrawAmount = 500

function init(args)
    math.randomseed(os.time() + tonumber(tostring({}):sub(8)))
end

function request()
    local targetWallet = (math.random(0, 1) == 0) and wallet1 or wallet2
    local isDeposit = (math.random(0, 1) == 1)
    local amount = isDeposit and depositAmount or withdrawAmount

    local body = string.format(
        '{"walletId":"%s","operationType":"%s","amount":%d}',
        targetWallet,
        isDeposit and "DEPOSIT" or "WITHDRAW",
        amount
    )

    return wrk.format("POST", baseUrl .. "/wallet",
                      {["Content-Type"] = "application/json"},
                      body)
end
