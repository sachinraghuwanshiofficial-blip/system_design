# Deep Dive: Understanding the Phase 1 Setup

When building distributed systems locally, it is critical to understand *why* we are configuring our environment a certain way. Let's break down every line of `docker-compose.yml` and the Golang initialization commands from `01_phase1_setup.md`.

## 1. Deconstructing `docker-compose.yml`

Docker Compose is a tool for defining and running multi-container Docker applications. Instead of running three massively long `docker run ...` commands in your terminal, we declare the desired state in a YAML file.

### The Service Definition
```yaml
services:
  db-master:
    image: mysql:8.0
```
- **`services:`**: This defines the different "computers" (containers) we want to run.
- **`db-master:`**: This is the internal network name of the container. Inside the Docker network, if another container pings `db-master`, Docker's internal DNS automatically resolves it to this container's IP address.
- **`image: mysql:8.0`**: Tells Docker to download the official MySQL version 8.0 blueprint from Docker Hub.

### The `command` Override
```yaml
    command: --server-id=1 --log-bin=mysql-bin --binlog-format=ROW
```
By default, the `mysql:8.0` image just starts the database. We are overriding the startup command to pass crucial replication arguments:
- **`--server-id=1`**: In a MySQL replication cluster, absolutely every node must have a unique integer ID (1 to 4294967295). If two nodes have the same ID, replication will crash immediately because the nodes cannot differentiate between their own events and the other node's events.
- **`--log-bin=mysql-bin`**: This turns on the Binary Log (which is off by default to save disk space). The database will now record every single `INSERT/UPDATE/DELETE` into files prefixed with `mysql-bin` (e.g., `mysql-bin.000001`). This is what the Replicas will read.
  - **What exactly is the Binary Log?** It is a sequential, append-only file stored on the database's hard drive. When a user runs an `UPDATE`, MySQL first changes the data in memory/disk, then instantly writes an entry to the end of the binlog saying "I just did this update". 
  - **Why `mysql-bin`?** The string `mysql-bin` is just the base filename. MySQL will automatically append sequence numbers to it. As the log gets too big (usually 1GB), MySQL creates `mysql-bin.000002`, `mysql-bin.000003`, and so on. This prevents a single file from consuming all the storage.
- **`--binlog-format=ROW`**: Specifically instructs MySQL to use Row-Based Replication (as discussed in `02a_replication_deep_dive.md`), ensuring perfect data accuracy.

### Environment Variables
```yaml
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: social_app
```
When the container boots up for the very first time, the MySQL initialization script reads these variables:
- **`MYSQL_ROOT_PASSWORD`**: Sets the master password.
- **`MYSQL_DATABASE`**: A massive convenience! It automatically runs `CREATE DATABASE social_app;` for us on boot so we don't have to log in and do it manually.

### Port Mapping
```yaml
    ports:
      - "3306:3306"  # Format is "HostPort : ContainerPort"
```
MySQL always runs on port `3306` *inside* its container.
- For `db-master`, we map Mac's port `3306` -> Container's `3306`.
- For `db-replica-1`, we map Mac's port `3307` -> Container's `3306`.
- For `db-replica-2`, we map Mac's port `3308` -> Container's `3306`.
This is why your Golang app can connect to three different databases dynamically just by changing the port number!

### Dependency Management
```yaml
    depends_on:
      - db-master
```
This tells Docker Compose to strictly boot `db-master` *before* attempting to boot the replicas. 

### The CLI Command: `docker-compose up -d`

After defining the YAML, you ran this command in the terminal. Here is exactly what is happening:
- **`docker-compose`**: The executable command that parses the `docker-compose.yml` file in your current directory.
- **`up`**: The instruction to create and start the containers. Docker looks at the `services` defined in the file. It will first create an isolated virtual network so the containers can talk to each other. Then, it checks `depends_on`, realizes it must start `db-master` first, and creates it. Once `db-master` is created, it parallel-boots `db-replica-1` and `db-replica-2`.
- **`-d` (Detached Mode)**: Without this flag, Docker would hijack your terminal window, constantly printing out the live server logs (startup messages, errors, etc.) of all 3 databases combined. If you hit `Ctrl+C`, it would instantly kill all the databases. By passing `-d`, you tell Docker to start the databases in the background (detached), allowing you to continue using your current terminal window to type new commands.

---

## 2. Deconstructing the Golang Commands

We ran a few seemingly magical commands to set up the router.

### `go mod init go-router`
This command creates a `go.mod` file in your directory. 
Unlike Python where dependencies are installed globally via `pip`, or Node.js where `package.json` needs to be manually created, `go mod init` establishes this folder as an isolated Go Project named `go-router`. It tracks exactly which versions of external libraries you are using so that if someone downloads your code 5 years from now, it still compiles perfectly.

### `go get github.com/go-sql-driver/mysql`
This downloads the official MySQL driver for Go from GitHub.
- It places the downloaded source code in your `$GOPATH` (Go's global cache).
- It updates the `go.mod` file to register that your project securely depends on this specific version of the driver.

### The Blank Import in `main.go`
```go
import (
    _ "github.com/go-sql-driver/mysql"
)
```
In Go, if you import a package and don't use it, the compiler throws an error and refuses to build.
However, we never explicitly call `mysql.DoSomething()`. We just use Go's built-in `database/sql` package to do all the work.
The underscore `_` is a **Blank Import**. It tells the Go compiler: "Download and compile this package, and run its internal `init()` function to secretly register itself as the handler for 'mysql', but don't complain that I'm not directly referencing it in my code."

---
## Summary
You are now running a Software Defined Network via Docker where three isolated computers are communicating, all exposed to your Mac via different port tunnels, and you have configured a compiled backend language to dynamically route connections between them!
