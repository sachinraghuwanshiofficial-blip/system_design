# Phase 5: Transactions (Advanced)

Transactions are one of the most critical concepts for maintaining data integrity in production applications. They ensure that a sequence of database operations is treated as a single, atomic unit.

## What is a Transaction?
Imagine a bank transfer. 
1. Deduct $100 from Account A.
2. Add $100 to Account B.

If the server crashes after step 1 but before step 2, Account A lost $100, but Account B never got it. The database state is corrupt. 
A **Transaction** ensures that either *all* steps execute successfully, or *none* of them do.

This guarantees the **ACID** properties:
- **A**tomicity: All or nothing.
- **C**onsistency: Data remains valid according to rules/constraints.
- **I**solation: Concurrent transactions don't interfere with each other.
- **D**urability: Once committed, it's saved permanently.

## Basic Syntax

```sql
-- Begin a new transaction
START TRANSACTION;

-- Execute multiple SQL statements
INSERT INTO ...
UPDATE ...
DELETE ...

-- If everything went well, save the changes permanently
COMMIT;

-- OR: If an error occurred, revert all changes since START TRANSACTION
ROLLBACK;
```

## A Social Media Example

Let's say a user reports a post for severe policy violations. As an admin, our application needs to immediately:
1. Delete the offending post.
2. Flag the user's account by updating a `status` column.
3. Add a record to an `audit_log` table tracking the admin action.

We wouldn't want the post deleted if the logging failed. We wrap it in a transaction.

*Setup for this example:*
```sql
ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'ACTIVE';
CREATE TABLE audit_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    action_type VARCHAR(50),
    details TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

*The Transaction:*
```sql
START TRANSACTION;

-- Step 1
DELETE FROM posts WHERE id = 5;

-- Step 2
UPDATE users SET status = 'BANNED' WHERE id = 2; -- assuming post 5 belonged to user 2

-- Step 3
INSERT INTO audit_log (action_type, details) 
VALUES ('CRITICAL_BAN', 'User ID 2 banned due to Post 5 violation.');

COMMIT;
```

In your terminal, you can manually type these out. You'll notice that before you type `COMMIT;`, if you open a *second* terminal window and log into mysql and run `SELECT * FROM users;`, you won't see the 'BANNED' status yet. This demonstrates **Isolation**. Once you run `COMMIT;` in the first window, the changes become visible everywhere.

If you make a mistake while typing out the commands in the transaction block, simply type `ROLLBACK;` and the database will revert to the state it was in before `START TRANSACTION`!
