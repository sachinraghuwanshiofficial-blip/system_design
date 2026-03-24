# 🗺 Horizontal Scaling Roadmap: Database Routing & Sharding

Welcome to the Horizontal Scaling course! This guide builds upon the foundational database concepts by teaching you how to architect databases to handle massive read and write loads. We will use **Docker Compose** to easily orchestrate multiple MySQL instances locally, and **Golang** to build a fake API server that handles our database routing logic.

## Phase 1: Environment Setup
1. **Docker Layout**: Spinning up a Primary (Master) database and two Read Replicas.
2. **Go Environment**: Initializing a Go module for our application router.
3. **Connectivity Check**: Ensuring our Go app can talk to the default database.
*📖 **Deep Dive**: `01a_setup_deep_dive.md` explains the Docker networking, volume mounts, server IDs, and Golang dependency commands.*

## Phase 2: Master-Slave Read Replication
1. **Configuring the Master**: Enabling the binary log (binlog) and creating a replication user.
2. **Connecting Replicas**: Instructing the Replicas to stream updates from the Master.
3. **Intelligent API Routing (Go)**: 
   - Modifying our Go application to hold multiple database connection pools.
   - Routing all `INSERT`, `UPDATE`, `DELETE` queries to the Master.
   - Routing all `SELECT` queries to the Replicas using a Round-Robin strategy.
*📖 **Deep Dive**: `02a_replication_deep_dive.md` covers Replication Lag, Binlog Formats, and Read-After-Write Consistency.*

## Phase 3: Sharding (Write Scaling)
1. **The Write Bottleneck**: Understanding when Read Replication is no longer enough (e.g., millions of new inserts per second).
2. **Adding Shard B**: Spinning up a second Master node (and its Replica) to handle half the traffic.
3. **Hash-Based Routing (Go)**: 
   - Updating the Go router to determine *which* Shard a user belongs to based on their `user_id` (e.g., `user_id % 2`).
*📖 **Deep Dive**: `03a_sharding_deep_dive.md` covers Consistent Hashing, Distributed ID Generation (Snowflake), and Resharding Migrations.*

## Phase 4: Advanced Concepts (Theory)
1. **Consistent Hashing**: Why simple modulo (`%`) breaks down when adding/removing shards, and how consistent hashing fixes it.
2. **Cross-Shard Joins**: The massive architectural headache of trying to `JOIN` data when User A lives on Shard 1, but User B lives on Shard 2.

---
**Let's begin! Head over to `01_phase1_setup.md` to start.**
