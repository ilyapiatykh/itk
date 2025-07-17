# itk

Test task for position golang june+

## How can it be improved?

- middleware for authentication processing
- metrics for requests
- database replication for reading only by slave node
- sharding by walletID
- native use of pgx driver
- prepared statements
- cors
- tls

## 📊 `wrk` Load Test Results

### General Parameters

| Parameter             | Value               |
|------------------------|---------------------|
| Duration              | 10.02 sec           |
| Threads               | 12                  |
| Connections           | 120                 |
| Total Requests        | 51,175              |
| Requests/sec          | 5,108.55            |
| Transfer Rate         | 0.95 MB/sec         |
| Total Data Transferred| 9.55 MB             |

---

### Latency

| Metric                 | Value       |
|------------------------|-------------|
| Average (`Avg`)        | 28.49 ms    |
| Standard Deviation     | 29.77 ms    |
| Maximum (`Max`)        | 402.82 ms   |
| Within ±1σ             | 89.51%      |

---

### Requests per Second per Thread

| Metric                 | Value       |
|------------------------|-------------|
| Average (`Avg`)        | 428.59      |
| Standard Deviation     | 99.74       |
| Maximum (`Max`)        | 696.00      |
| Within ±1σ             | 78.30%      |

---

### Latency Distribution

| Percentile | Latency     |
|------------|-------------|
| 50% (Median) | 21.12 ms   |
| 75%          | 37.43 ms   |
| 90%          | 59.48 ms   |
| 99%          | 143.94 ms  |
