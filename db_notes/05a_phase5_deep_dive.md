# Deep Dive: ACID, Concurrency Control, and Locks

Transactions are essential for data integrity when operations require multiple steps. However, the real complexity of databases arises when *hundreds of users* are executing transactions concurrently.

## 1. Defining Concurrency Problems

When applications scale, the database handles thousands of connections simultaneously. Without strict controls, transactions will collide and corrupt data. Consider the classic bank balance problem: Let's assume an account has $100.

**The "Lost Update" Problem:**
- **Tx A** (Transaction A) reads the balance ($100).
- **Tx B** reads the balance ($100).
- **Tx A** adds $50, writing back $150.
- **Tx B** subtracts $20, writing back $80.
Tx A's deposit is lost forever because Tx B overwrote it based on stale data!

## 2. Isolation Levels

The 'I' in ACID stands for *Isolation*. The database provides different "levels" of isolation depending on how much concurrency vs. data safety your application needs.

1. **Read Uncommitted (Highest speed,Lowest safety):** Tx A can see data that Tx B changed *but hasn't committed yet* (a "Dirty Read"). Extremely dangerous for financial data.
2. **Read Committed:** Tx A only sees data that Tx B has fully committed. However, if Tx A reads the same row twice, the value might change between reads if another transaction committed an update.
3. **Repeatable Read (MySQL/InnoDB Default):** Once Tx A reads a row, it takes a snapshot. The row's data is guaranteed not to change for the duration of Tx A, even if Tx B updates and commits it. This resolves most anomalies.
4. **Serializable (Lowest speed, Highest safety):** The database forces transactions to execute as if they were running sequentially (one after another). No concurrency is allowed for contested rows. This prevents all anomalies but creates severe performance bottlenecks.

## 3. How Databases Manage Concurrency: Locks

To implement Isolation, MySQL's InnoDB engine uses **Locks**.

### Row-Level vs. Table-Level Locks
- **Table Lock:** The crudest method. If Tx A wants to update user ID 1, it locks the entire `users` table. No other transaction can update *any* user until Tx A finishes. This kills performance.
- **Row-Level Lock (InnoDB's specialty):** If Tx A updates user ID 1, InnoDB locks *only* that specific row. Tx B can simultaneously update user ID 2 without waiting.

### Pessimistic vs. Optimistic Locking

**Pessimistic Locking (Database handled):**
Assuming conflicts *will* happen. If you need to read a row and plan to update it within a transaction, you can lock it immediately when reading:
```sql
START TRANSACTION;
SELECT balance FROM accounts WHERE id = 1 FOR UPDATE;
-- Any other Tx trying to read or write ID 1 will PAUSE and wait here until this Tx finishes.
UPDATE accounts SET balance = balance + 50 WHERE id = 1;
COMMIT;
```

**Optimistic Locking (Application handled):**
Assuming conflicts are *rare*. We add a `version` column to the table.
The application reads the row: `balance=100, version=1`.
The application calculates the new balance (150) and updates:
```sql
UPDATE accounts 
SET balance = 150, version = 2 
WHERE id = 1 AND version = 1;
```
If Tx B sneaked in and updated the row first, the `version` the database holds would now be 2. The `UPDATE` statement above would affect 0 rows. The application detects this and retries the operation from scratch. This is highly scalable because it never holds database locks!

## 4. Deadlocks

A deadlock occurs when two transactions hold locks the other needs, and neither can proceed.

1. **Tx A** locks Row X.
2. **Tx B** locks Row Y.
3. **Tx A** tries to update Row Y (and waits for Tx B's lock).
4. **Tx B** tries to update Row X (and waits for Tx A's lock).

Both transactions are frozen forever.
*How MySQL handles it:* InnoDB instantly detects deadlocks. It proactively aborts ("kills") the transaction that has done the least amount of work, returning an error to the application code, and allowing the other transaction to succeed. Your application code must be built to catch deadlock errors and retry the transaction!
