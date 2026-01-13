# Discussion: High-Performance Architecture (RAM + Custom WAL)

## 1. The Core Concept: "No-SQLite" for Hot Path

The user proposes a **Custom In-Memory Engine** backed by a **Write-Ahead Log (WAL)**. This bypasses SQLite entirely for the high-concurrency "Drop" logic.

- **Architecture:** `In-Memory Go Structs` (Primary State) + `Append-Only File` (Durability).
- **Performance:** Limited only by how fast we can write bytes to the end of a file (OS Page Cache makes this near-instant) and Go mutex speed.
- **Why:** SQLite has overhead (SQL parsing, B-Tree locks, transaction management). For a "Limited Drop", we only need to atomically decrement a number and record a winner. We don't need a full Relational DB engine for the hot path.

---

## 2. Detailed Design

### A. The State (RAM)

We maintain the entire "Drop" state in memory using Go data structures.

```go
type DropState struct {
    TotalStock uint32
    Sold       uint32
    Orders     []OrderRecord // or a map for quick lookup
    mu         sync.RWMutex  // Or use a single-threaded actor (channel)
}

var GlobalDrops = map[uint64]*DropState{}
```

### B. The Persistence (WAL)

We do **not** use a database for durability. We use a simple **Append-Only Log File** (e.g., `events.log`).
Before updating the RAM state, we modify the log.

**The Flow:**

1.  **Request:** User buys Drop #1.
2.  **Serialize:** Create a binary/JSON entry: `{"op": "BUY", "drop_id": 1, "user": "namdt", "ts": 123456}`.
3.  **WAL Append:** Write this entry to `drop_1.wal`.
    - _Critical:_ Requires `fsync` (or OS cache if we trust OS stability) to ensure data is on disk.
4.  **Memory Update:** `DropState.Sold++`.
5.  **Response:** Return "Success" to user.

### C. Recovery (Startup)

When the server restarts:

1.  Read `drop_1.wal` from beginning to end.
2.  Replay every "BUY" event against the in-memory counter.
3.  Rebuild the `DropState` to match exactly where it left off.

---

## 3. Comparison

| Feature        | Standard SQLite   | Custom RAM + WAL                                          |
| :------------- | :---------------- | :-------------------------------------------------------- |
| **Speed**      | 100 - 1,000 IOPS  | **100,000+ IOPS** (Limited by sequential write speed/RAM) |
| **Latency**    | Milliseconds      | **Microseconds**                                          |
| **Durability** | High (ACID)       | **High** (If WAL is synced)                               |
| **Complexity** | Low (SQL queries) | **High** (Must write custom recovery & snapshot logic)    |
| **Querying**   | Flexible SQL      | **Hard** (Can only query what matches in-memory structs)  |

---

## 4. Implementation Plan (The "LMAX" Lite)

1.  **Event Definition:** Define strict Event structs (`EventOrderPlaced`, `EventStockAdjusted`).
2.  **WAL Manager:** A service that handles opening files and appending bytes thread-safely.
3.  **Snapshotting (Optional but recommended):**
    - If the WAL gets too big (1GB+), we dump the current `DropState` to a `snapshot.bin` and clear the WAL.
    - Recovery = Load `snapshot.bin` + Replay remaining WAL.
4.  **Integration:**
    - `PurchaseDrop` handler -> Sends event to WAL Log -> Updates RAM -> Returns.
    - Background Job -> Reads RAM state periodically and syncs to "Main" SQLite (for Admin/Reporting tools to see).

## 5. Conclusion

This architecture transforms the backend into a specialized **Trading Engine**. It is the absolute fastest way to handle "Flash Sales" without dedicated hardware.
