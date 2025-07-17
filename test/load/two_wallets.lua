-- usage example:
-- wrk -t12 -c120 -d60s -s wallet_load_test_two.lua --latency http://localhost:80

local baseUrl = "/api/v1"
local wallet1 = "11111111-1111-1111-1111-111111111111"
local wallet2 = "22222222-2222-2222-2222-222222222222"
local depositAmount = 1000
local withdrawAmount = 500

local stats = {
    requests = 0,
    errors_5xx = 0,
    wallet1_initial = 0,
    wallet2_initial = 0,
    wallet1_expected = 0,
    wallet2_expected = 0,
    wallet1_final = 0,
    wallet2_final = 0
}

function getBalance(walletId)
    local req = wrk.format("GET", baseUrl .. "/wallets/" .. walletId)
    local status, _, body = wrk.request(req)
    if status == 200 then
        return tonumber(body)
    end
    return nil
end

function init(args)
    local createBody = function(id)
        return string.format('{"walletId":"%s"}', id)
    end

    local res1 = wrk.request(wrk.format("POST", baseUrl .. "/wallet",
                                      {["Content-Type"]="application/json"},
                                      createBody(wallet1)))

    local res2 = wrk.request(wrk.format("POST", baseUrl .. "/wallet",
                                      {["Content-Type"]="application/json"},
                                      createBody(wallet2)))

    stats.wallet1_initial = getBalance(wallet1) or 0
    stats.wallet2_initial = getBalance(wallet2) or 0

    print(string.format("Initial balances - Wallet1: %d, Wallet2: %d",
          stats.wallet1_initial, stats.wallet2_initial))
end

function request()
    stats.requests = stats.requests + 1

    local targetWallet = (stats.requests % 2 == 0) and wallet1 or wallet2
    local isDeposit = (math.random(0, 1) == 1)
    local amount = isDeposit and depositAmount or withdrawAmount

    if targetWallet == wallet1 then
        stats.wallet1_expected = stats.wallet1_expected + (isDeposit and amount or -amount)
    else
        stats.wallet2_expected = stats.wallet2_expected + (isDeposit and amount or -amount)
    end

    local body = string.format(
        '{"walletId":"%s","operationType":"%s","amount":%d}',
        targetWallet,
        isDeposit and "DEPOSIT" or "WITHDRAW",
        amount
    )

    return wrk.format("POST", baseUrl .. "/wallet",
                     {["Content-Type"]="application/json"},
                     body)
end

function response(status, headers, body)
    if status >= 500 and status < 600 then
        stats.errors_5xx = stats.errors_5xx + 1
    end
end

function done(summary, latency, requests)
    stats.wallet1_final = getBalance(wallet1) or 0
    stats.wallet2_final = getBalance(wallet2) or 0

    print("\n=== Test Results ===")
    print(string.format("Total requests: %d", stats.requests))
    print(string.format("5xx errors: %d (%.2f%%)",
          stats.errors_5xx, (stats.errors_5xx/stats.requests)*100))

    print("\nWallet 1:")
    print(string.format("  Initial balance: %d", stats.wallet1_initial))
    print(string.format("  Expected change: %+d", stats.wallet1_expected))
    print(string.format("  Expected final: %d", stats.wallet1_initial + stats.wallet1_expected))
    print(string.format("  Actual final: %d", stats.wallet1_final))
    print(string.format("  Discrepancy: %d",
          stats.wallet1_final - (stats.wallet1_initial + stats.wallet1_expected)))

    print("\nWallet 2:")
    print(string.format("  Initial balance: %d", stats.wallet2_initial))
    print(string.format("  Expected change: %+d", stats.wallet2_expected))
    print(string.format("  Expected final: %d", stats.wallet2_initial + stats.wallet2_expected))
    print(string.format("  Actual final: %d", stats.wallet2_final))
    print(string.format("  Discrepancy: %d",
          stats.wallet2_final - (stats.wallet2_initial + stats.wallet2_expected)))

    print("\nPerformance:")
    print(string.format("  Requests/sec: %.2f", stats.requests/summary.duration))
    print(string.format("  Latency (mean): %.2fms", latency.mean/1000))
    print(string.format("  Latency (p99): %.2fms", latency:percentile(99)/1000))
end
