# Phase 1: Environment Setup (Docker & Go)

To properly learn horizontal scaling and routing, we need multiple instances of a database running simultaneously. Trying to install and configure 3 separate MySQL servers manually on a Mac is guaranteed to fail due to conflicting ports and data directories. 

Instead, we will use **Docker Compose**. Docker allows us to spin up isolated containers, each running their own MySQL engine, on different ports.

## 1. Setting up the Docker Environment

1. Ensure Docker Desktop is installed and running on your Mac.
2. In your terminal, navigate to the `horizontal_scaling` folder (where this file belongs).
3. Create a file named `docker-compose.yml` and paste the following:

```yaml
version: '3.8'

services:
  db-master:
    image: mysql:8.0
    command: --server-id=1 --log-bin=mysql-bin --binlog-format=ROW
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: social_app
    ports:
      - "3306:3306"

  db-replica-1:
    image: mysql:8.0
    command: --server-id=2
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
    ports:
      - "3307:3306"
    depends_on:
      - db-master

  db-replica-2:
    image: mysql:8.0
    command: --server-id=3
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
    ports:
      - "3308:3306"
    depends_on:
      - db-master
```

### Explaining the Setup:
- **`db-master`** runs on your machine's default port `3306`. It starts with `--log-bin` enabled, tracking every single write.
- **`db-replica-1`** runs on port `3307`. 
- **`db-replica-2`** runs on port `3308`.
- Note the distinct `--server-id` (1, 2, 3). This is mandatory for MySQL replication so servers know who is talking.

**Start the cluster!**
*(Run this in the directory with your `docker-compose.yml`)*
```bash
docker-compose up -d
```
You can check they are running with `docker ps`.

## 2. Setting up the Golang Router Foundation

Our database nodes are running, but who talks to them? In a real architecture, the backend application makes the choice of *where* to send a query. Let's create a minimal Go application.

1. Create a `go-router` directory and initialize the module:
```bash
mkdir go-router && cd go-router
go mod init go-router
go get github.com/go-sql-driver/mysql
```

2. Create a `main.go` file inside `go-router/`:

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

// DBCluster holds connections to our Master and Replicas
type DBCluster struct {
    Master   *sql.DB
    Replicas []*sql.DB
}

func main() {
    // 1. Connect to Master (Port 3306)
    master, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/social_app")
    if err != nil { log.Fatal("Error connecting to Master:", err) }
    defer master.Close()

    // 2. Connect to Replica 1 (Port 3307)
    replica1, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3307)/")
    if err != nil { log.Fatal("Error connecting to Replica 1:", err) }
    defer replica1.Close()

    // 3. Connect to Replica 2 (Port 3308)
    replica2, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3308)/")
    if err != nil { log.Fatal("Error connecting to Replica 2:", err) }
    defer replica2.Close()

    // Create our cluster representation
    cluster := &DBCluster{
        Master:   master,
        Replicas: []*sql.DB{replica1, replica2},
    }

    fmt.Println("✅ Connection pools established to Master and Replicas.")
    
    // We will build the routing logic in Phase 2!
}
```

You can test that it connects to your Docker containers by running:
```bash
go run main.go
```

---
**Next Step:** Head over to `02_phase2_read_replication.md` to actually connect the replication streams and build Go's query routing!
