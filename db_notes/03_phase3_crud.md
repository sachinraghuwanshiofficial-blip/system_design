# Phase 3: Writing Data (CRUD Basics)

Now that we have our `users`, `posts`, and `followers` tables, it's time to populate them and learn the four basic database operations: **C**reate, **R**ead, **U**pdate, and **D**elete.

## 1. CREATE (Insert Data)

The `INSERT INTO` statement adds new rows to a table.

**Inserting Users:**
```sql
INSERT INTO users (username, email) 
VALUES ('alice', 'alice@example.com');

INSERT INTO users (username, email) 
VALUES ('bob', 'bob@example.com'), 
       ('charlie', 'charlie@example.com');
```
*Notice we didn't specify the `id` or `created_at` fields. The database automatically generates those for us thanks to `AUTO_INCREMENT` and `DEFAULT CURRENT_TIMESTAMP`.*

**Inserting Posts:**
Assuming Alice's ID is 1, and Bob's ID is 2:
```sql
INSERT INTO posts (user_id, content) 
VALUES (1, 'Hello World! This is my first post.');

INSERT INTO posts (user_id, content) 
VALUES (2, 'Bob here. Excited to join this network.'),
       (1, 'Its Alice again! Loving this application.');
```

**Inserting Followers:**
Let's have Bob (2) and Charlie (3) follow Alice (1).
```sql
INSERT INTO followers (follower_id, followee_id) 
VALUES (2, 1), 
       (3, 1);
```

## 2. READ (Select Data)

The `SELECT` statement retrieves data.

**Fetch all data from a table:**
```sql
SELECT * FROM users;
```

**Fetch specific columns:**
```sql
SELECT username, email FROM users;
```

**Filtering with `WHERE` clauses:**
```sql
SELECT * FROM posts WHERE user_id = 1;
```

**Sorting with `ORDER BY`:**
```sql
SELECT * FROM posts ORDER BY created_at DESC;
```
*(Gets the newest posts first)*

**Limiting results:**
```sql
SELECT * FROM users LIMIT 2;
```

## 3. UPDATE (Modify Data)

The `UPDATE` statement modifies existing records. **Always use a `WHERE` clause with UPDATE, otherwise you will modify every row in the table!**

**Changing a user's email:**
```sql
UPDATE users 
SET email = 'alice.new@example.com' 
WHERE username = 'alice';
```

## 4. DELETE (Remove Data)

The `DELETE` statement removes rows from a table. Like `UPDATE`, **always use a `WHERE` clause!**

**Deleting a specific post:**
```sql
DELETE FROM posts WHERE id = 1;
```

**Testing `ON DELETE CASCADE`:**
Try deleting a user:
```sql
DELETE FROM users WHERE id = 1;
```
*Because we added `ON DELETE CASCADE` when creating our schema in Phase 2, MySQL will automatically delete Alice's posts and follower records, ensuring data integrity.*
