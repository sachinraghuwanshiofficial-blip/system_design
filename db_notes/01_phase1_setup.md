# Phase 1: Environment Setup

This phase is all about getting MySQL running on your local machine and learning how to interact with it via the command-line interface (CLI). By the end, you should have a running database server and be logged into the MySQL prompt.

## 1. Installation

Since you're on a Mac and want to use the terminal, **Homebrew** is the easiest package manager for this.

**Command:**
```bash
brew install mysql
```
*What this does:* It downloads the MySQL server software and client tools and places them in their appropriate directories on your Mac.

## 2. Managing the MySQL Service

A database is not just a command; it's an always-on "service" or "daemon" that runs in the background, listening for connections on a specific port (default is `3306`).

**To start MySQL as a background service:**
```bash
brew services start mysql
```
```
*Note:* The database will now run silently. It will even restart if you reboot your Mac.

**To stop MySQL:**
```bash
brew services stop mysql
```

**To restart MySQL (useful if you change config files):**
```bash
brew services restart mysql
```

## 3. Connecting via the CLI

Now that the engine is running, you need a way to talk to it. The `mysql` command-line client is exactly for this.

**Command:**
```bash
mysql -u root
```
*Breakdown:*
- `mysql`: Invokes the client program.
- `-u root`: Tells the client you want to log in as the user named `root`. Root is the default superadmin account created during installation. By default on Homebrew installations, root doesn't have a password.

*If you are prompted for a password or want to set one up, use: `mysql -u root -p`*

## 4. Basic CLI Navigation

Once you see the `mysql>` prompt, you are no longer in your standard bash/zsh terminal. You are inside the MySQL monitor. Every command you type here must end with a semicolon (`;`).

**Useful initial commands:**
- `SHOW DATABASES;` -> Lists all databases currently on the server.
- `SELECT VERSION();` -> Shows which version of MySQL you are running.
- `SELECT CURRENT_DATE;` -> Just a test to see if it can return data.
- `exit` or `quit` -> Leaves the MySQL monitor and brings you back to your regular terminal. (No semicolon needed for exit).

## Troubleshooting:
If you see an error like `ERROR 2002 (HY000): Can't connect to local MySQL server through socket`, it means the background service is *not* running. Check Step 2 (Start MySQL).
