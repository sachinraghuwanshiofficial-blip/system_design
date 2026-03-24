# Phase 2: Master-Slave Read Replication

Now that we have three running MySQL nodes, they are entirely independent. To make them a "Cluster", we must tell `db-replica-1` and `db-replica-2` to listen to everything `db-master` does and copy it perfectly.

## 1. Configuring the Master

Connect to the Master node (running on port `3306`). The easiest way using the MySQL CLI is:
```bash
mysql -h 127.0.0.1 -P 3306 -u root -prootpassword
```

1. Create a special user just for replication:
```sql
CREATE USER 'replica_user'@'%' IDENTIFIED BY 'replica_password';
GRANT REPLICATION SLAVE ON *.* TO 'replica_user'@'%';
FLUSH PRIVILEGES;
```

2. Check the Master Status. **Note down the `File` and `Position`!**
```sql
SHOW MASTER STATUS;
```
*(Example Output: File = `mysql-bin.000001`, Position = `156`)*

## 2. Configuring the Replicas

Connect to `db-replica-1` (port `3307`):
```bash
mysql -h 127.0.0.1 -P 3307 -u root -prootpassword
```

Tell it exactly where to find the Master and where to start reading the binlog (use the File and Position you noted from above):
```sql
CHANGE MASTER TO 
  MASTER_HOST='db-master',        -- It uses the Docker service name!
  MASTER_USER='replica_user',
  MASTER_PASSWORD='replica_password',
  MASTER_LOG_FILE='mysql-bin.000001', 
  MASTER_LOG_POS=156;

START SLAVE;
SHOW SLAVE STATUS\G;  -- Look for Slave_IO_Running: Yes and Slave_SQL_Running: Yes
```

Repeat this exact same process for `db-replica-2` (port `3308`). 

Now, if you `CREATE TABLE users ...` on the Master, it will instantly appear on both Replicas!

## 3. Intelligent API Architecture (Golang)

In a traditional setup, you have one `DB_URL` connection string. Here, your application must be smart enough to split Reads and Writes.

Let's expand our `main.go`. First, let's create a global variable to track which replica we used last (for Round-Robin load balancing).

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "math/rand"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// Round-robin counter
var nextReplicaIndex = 0

// DBCluster holds connections to our Master and Replicas
type DBCluster struct {
    Master   *sql.DB
    Replicas []*sql.DB
}

// Write queries ALWAYS go to the Master
func (cluster *DBCluster) WriteQuery(query string, args ...interface{}) (sql.Result, error) {
    fmt.Println("🚀 Routing Write to Master (Port 3306)")
    return cluster.Master.Exec(query, args...)
}

// Read queries should be load-balanced across Replicas
func (cluster *DBCluster) ReadQuery(query string, args ...interface{}) (*sql.Rows, error) {
    // Round-Robin Logic
    replicaIndex := nextReplicaIndex % len(cluster.Replicas)
    nextReplicaIndex++ // Increment for the next read

    fmt.Printf("📖 Routing Read to Replica %d\n", replicaIndex+1)
    
    // Fallback: If a replica is down, the routing logic gets more complex here!
    return cluster.Replicas[replicaIndex].Query(query, args...)
}

func main() {
    // ... [Previous connections code from Phase 1] ...
    
    // Create the cluster
    cluster := &DBCluster{
        Master:   master,
        Replicas: []*sql.DB{replica1, replica2},
    }

    // --- Simulating App Traffic ---
    
    // 1. A User signs up (WRITE)
    _, err := cluster.WriteQuery("INSERT INTO social_app.users (username, email) VALUES (?, ?)", "sachin", "sachin@example.com")
    if err != nil { fmt.Println("Failed Write:", err) }

    // 2. Someone Views a Profile (READ)
    // The first read will go to Replica 1
    rows, _ := cluster.ReadQuery("SELECT * FROM social_app.users WHERE username = ?", "sachin")
    defer rows.Close()

    // 3. Another person views a profile (READ)
    // This read will automatically go to Replica 2!
    cluster.ReadQuery("SELECT * FROM social_app.users")

    // 4. Yet another read (READ)
    // This read loops back around and goes to Replica 1!
    cluster.ReadQuery("SELECT * FROM social_app.users")
}
```

This is how massive applications like Instagram handle millions of reads per second: by spinning up hundreds of Replicas and routing `SELECT` queries perfectly across them!

---
**Next Step:** Head to `03_phase3_sharding.md` to learn what happens when even the *Master* gets too much Write traffic.
