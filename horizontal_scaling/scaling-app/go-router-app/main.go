package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Round-robin counter for load balancing Reads
var nextReplicaIndex = 0

// DBCluster holds our connection pools
type DBCluster struct {
	Master   *sql.DB
	Replicas []*sql.DB
}

// Write queries ALWAYS go to the Master
func (cluster *DBCluster) WriteQuery(query string, args ...interface{}) (sql.Result, error) {
	fmt.Println("   -> 🚀 Routing Write to Master...")
	return cluster.Master.Exec(query, args...)
}

// Read queries are load-balanced across all Replicas
func (cluster *DBCluster) ReadQuery(query string, args ...interface{}) (*sql.Rows, error) {
	// Simple Round-Robin Logic
	replicaIndex := nextReplicaIndex % len(cluster.Replicas)
	nextReplicaIndex++ // Increment for the next read

	fmt.Printf("   -> 📖 Routing Read to Replica %d...\n", replicaIndex+1)
	return cluster.Replicas[replicaIndex].Query(query, args...)
}

// Phase 3: A single Shard is a Mini-Cluster (1 Master, N Replicas)
type Shard struct {
	Name    string
	Cluster *DBCluster
}

// Our entire globally distributed Data Tier
type DataTier struct {
	Shards []*Shard
}

// The core brain of horizontal scaling for Writers!
func (dt *DataTier) GetShardForUser(userID int) *Shard {
	// Modulo Hash Routing Logic
	totalShards := len(dt.Shards)
	shardIndex := userID % totalShards

	fmt.Printf("\n🔀 User %d routes to ====> %s\n", userID, dt.Shards[shardIndex].Name)
	return dt.Shards[shardIndex]
}

func main() {
	// 1. Setup Database Connections
	master, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/social_app")
	if err != nil {
		log.Fatal("Error connecting to the master: ", err)
	}
	defer master.Close()

	replica1, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3307)/social_app")
	if err != nil {
		log.Fatal("Error connecting to replica1: ", err)
	}
	defer replica1.Close()

	replica2, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3308)/social_app")
	if err != nil {
		log.Fatal("Error connecting to replica2: ", err)
	}
	defer replica2.Close()

	clusterConnections := &DBCluster{
		Master:   master,
		Replicas: []*sql.DB{replica1, replica2},
	}

	fmt.Println("✅ All database connections established.")

	// 2. Initialize the Sharded Architecture
	shardA := &Shard{Name: "🇺🇸 Shard-A (US-East)", Cluster: clusterConnections}
	shardB := &Shard{Name: "🇪🇺 Shard-B (EU-West)", Cluster: clusterConnections}

	system := &DataTier{
		Shards: []*Shard{shardA, shardB},
	}

	fmt.Println("🌍 Sharded Data Tier Online. Total Shards:", len(system.Shards))

	// --- Simulating Application Traffic ---
	fmt.Println("\n--- Initiating Phase 3 Traffic Simulation ---")

	userSachinID := 301 // Odd number -> 301 % 2 = 1 (Shard-B)
	userIshaID := 302   // Even number -> 302 % 2 = 0 (Shard-A)

	// User 301 writes a record
	targetShard1 := system.GetShardForUser(userSachinID)
	_, err = targetShard1.Cluster.WriteQuery("INSERT INTO social_app.users (id, username, email) VALUES (?, ?, ?)", userSachinID, "sachin_sharded", "sachin@shard.com")
	if err != nil {
		fmt.Printf("   -> ❌ DB Error: %v\n", err)
	}

	// User 302 writes a record
	targetShard2 := system.GetShardForUser(userIshaID)
	_, err = targetShard2.Cluster.WriteQuery("INSERT INTO social_app.users (id, username, email) VALUES (?, ?, ?)", userIshaID, "isha_sharded", "isha@shard.com")
	if err != nil {
		fmt.Printf("   -> ❌ DB Error: %v\n", err)
	}

	// Wait 100 milliseconds for Master Binlog to stream to Replicas
	time.Sleep(100 * time.Millisecond)

	// User 301 reads their own profile
	targetShard1Again := system.GetShardForUser(userSachinID)
	rows, err := targetShard1Again.Cluster.ReadQuery("SELECT username FROM social_app.users WHERE id = ?", userSachinID)
	if err == nil {
		var username string
		if rows.Next() {
			rows.Scan(&username)
			fmt.Printf("   -> ✅ Read Success! Found username: %s\n", username)
		} else {
			fmt.Println("   -> ⚠️ No user found in replica!")
		}
		rows.Close()
	} else {
		fmt.Printf("   -> ❌ Read Error: %v\n", err)
	}
}
