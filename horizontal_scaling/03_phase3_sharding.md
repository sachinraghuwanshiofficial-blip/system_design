# Phase 3: Sharding (Write Scaling)

Read Replication is amazing for heavy "Read" applications like News websites or Twitter feeds. But what if you are building WhatsApp, where almost every action is a Write (`INSERT` message)? 

A single Master database can only handle so many Writes before the hard drive/CPU maxes out (typically around 10,000-50,000 writes per second).

When your Master hits its limit, you must split your data across *two* Masters. This is called **Sharding**.

## 1. The Sharding Architecture

Imagine `Master A` and `Master B`. 
- Users 1, 3, 5, 7 live on `Master A`.
- Users 2, 4, 6, 8 live on `Master B`.

If `Master A` crashes, half your users are down, but the *other half* can still use the app perfectly fine! This reduces your "Blast Radius".

*(Note: We won't spin up another 2 Docker containers right now to save your laptop's memory, but let's implement the routing logic as if we had them!)*

## 2. Hash-Based Routing Strategy (Golang)

When a request comes in (e.g., "Add a post for User 45"), how does your backend know which Database to talk to?

The most common strategy is **Modulo Hashing**: `Shard ID = User_ID % Total_Shards`.

Let's refactor our `main.go`:

```go
package main

import (
    "database/sql"
    "fmt"
    // ... imports ...
)

// A single Shard is a Mini-Cluster (1 Master, N Replicas)
type Shard struct {
    Name     string
    Master   *sql.DB
    Replicas []*sql.DB
}

// Our entire Data Tier
type DataTier struct {
    Shards []*Shard
}

// The core brain of horizontal scaling!
func (dt *DataTier) GetShardForUser(userID int) *Shard {
    // Hash Routing Logic
    totalShards := len(dt.Shards)
    shardIndex := userID % totalShards
    
    fmt.Printf("🔀 User %d routes to %s\n", userID, dt.Shards[shardIndex].Name)
    return dt.Shards[shardIndex]
}

func main() {
    // 1. Establish connections to Shard A (e.g. US-East)
    shardA := &Shard{
        Name: "Shard-A",
        // Master: ..., Replicas: ...
    }

    // 2. Establish connections to Shard B (e.g. EU-West)
    shardB := &Shard{
        Name: "Shard-B",
        // Master: ..., Replicas: ...
    }

    // 3. Assemble the Fleet
    system := &DataTier{
        Shards: []*Shard{shardA, shardB},
    }

    // --- Simulating App Traffic ---
    
    userSachinID := 101 // Odd number
    userIshaID := 102   // Even number

    // 1. Sachin posts a message
    targetShard := system.GetShardForUser(userSachinID)
    // targetShard.Master.Exec("INSERT INTO posts ...")  // Goes to Shard-B (101 % 2 = 1)

    // 2. Isha posts a message
    targetShard2 := system.GetShardForUser(userIshaID)
    // targetShard2.Master.Exec("INSERT INTO posts...")  // Goes to Shard-A (102 % 2 = 0)
}
```

## 3. The Nightmare of Sharding: Joins & Rebalancing

Sharding solves the Write Bottleneck, but it introduces massive pain:

1. **Cross-Shard Joins are Impossible:** 
   If User `101` (Shard A) wants to see posts from User `102` (Shard B), you cannot write a `JOIN` query. Your Golang application must query Shard A, then independently query Shard B, and stitch the data together in application memory!
   
2. **Rebalancing & Resharding:**
   What happens when you need to add a 3rd Shard? 
   Previously `101 % 2 = 1` (Shard B). 
   Now `101 % 3 = 2` (Shard C).
   Suddenly, your application routing sends requests for User 101 to Shard C, but their data is still stuck on Shard B! 
   Moving terabytes of data live without downtime is one of the hardest problems in database engineering. (The solution involves *Consistent Hashing*).

---
**Congratulations!** 
You have completed the database scaling course. You've gone from a single SQLite file to a globally distributed, multi-master sharded architecture!
