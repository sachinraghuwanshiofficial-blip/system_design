# Phase 2: Database & Schema Foundations

Now that you're logged into the `mysql>` prompt, it's time to build the structure of our Social Media App. A database is essentially a collection of tables, and the "Schema" is the blueprint of what columns and data types those tables contain.

## 1. Creating and Selecting the Database

Before creating tables, we need a "bucket" to hold them.

**Command:**
```sql
CREATE DATABASE social_app;
```

**Verifying it exists:**
```sql
SHOW DATABASES;
```

**Using it (Crucial step!):**
```sql
USE social_app;
```
*Note:* You must tell MySQL which database you want to operate on. Once you run `USE social_app;`, all subsequent commands will apply to this database.

## 2. Introduction to Data Types and Constraints

When creating a table, every column must have a defined data type. Here are the most common ones we'll use:
- `INT`: Whole numbers (useful for IDs).
- `VARCHAR(n)`: Variable-length string text, up to `n` characters (useful for usernames, emails).
- `TEXT`: Long-form text (useful for post content).
- `TIMESTAMP`: Dates and times.

**Constraints** are rules applied to columns:
- `PRIMARY KEY`: Uniquely identifies each row in a table. It cannot be NULL. Customarily, this is an auto-incrementing integer.
- `FOREIGN KEY`: A link to a `PRIMARY KEY` in another table. This enforces "Referential Integrity" (e.g., a post cannot belong to a user that doesn't exist).
- `NOT NULL`: Ensures that a column cannot be left blank.
- `UNIQUE`: Ensures all values in a column are entirely distinct from one another.
- `AUTO_INCREMENT`: Automatically generates the next integer in the sequence when a new row is inserted.
- `DEFAULT`: Sets a default value if none is provided.

## 3. Creating the `users` Table

Let's model the user profile.

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**To verify the table was created and see its structure:**
```sql
SHOW TABLES;
DESCRIBE users;
```

## 4. Creating the `posts` Table (One-to-Many Relationship)

A user can have many posts. Thus, the `posts` table needs a way to point back to the `users` table. This is done via a Foreign Key.

```sql
CREATE TABLE posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```
*Note:* `ON DELETE CASCADE` is very powerful. It means if a user is deleted from the `users` table, MySQL will automatically delete all of their posts to prevent orphaned records!

## 5. Creating the `followers` Table (Many-to-Many Relationship)

Users can follow many users, and be followed by many users. To model a many-to-many relationship, we use a "Join Table" (or mapping table).

```sql
CREATE TABLE followers (
    follower_id INT NOT NULL,
    followee_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, followee_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE
);
```
*Note:* 
- `follower_id` is the person who clicks "Follow".
- `followee_id` is the person receiving the follow.
- Notice the `PRIMARY KEY (follower_id, followee_id)`. This is a *composite primary key*. It ensures that User A can only follow User B once. If they try to follow them again, MySQL will throw an error rather than creating a duplicate row.

---
**Next Step (Data Insertion):** The actual exercises to insert dummy data (users, posts, followers) are located in the next phase file: `03_phase3_crud.md`.
