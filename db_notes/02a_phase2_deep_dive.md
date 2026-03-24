# Deep Dive: Relational Theory and Schema Design

Defining your schema is the most consequential step in database development. If you get the schema wrong, your application code becomes complex, queries run slowly, and data inconsistencies creep in over time.

## 1. The Relational Model

Invented by Edgar F. Codd at IBM in 1970, the relational model revolutionized databases. Before this, data was often stored in rigid, hierarchical trees.

The core idea: Data is organized into tables (Relations). Each table represents one *entity* (e.g., Users, Posts). Rows represent *instances* of that entity. Columns represent *attributes*.

Most importantly, relationships between entities are formed not by physical links on the hard drive, but by **logical values** matching across tables (Primary Keys and Foreign Keys).

## 2. Normalization: The Art of Storing Data Once

**Normalization** is the process of structuring a database to reduce data redundancy and improve data integrity. The goal is simple: **Every non-primary-key column should describe the primary key, the whole primary key, and nothing but the primary key.**

Imagine a bad, "un-normalized" schema for a Post:

| post_id | content | author_name | author_email | created_at |
| :--- | :--- | :--- | :--- | :--- |
| 1 | "Hello!" | Alice | alice@example.com | 2023-10-01 |
| 2 | "My 2nd post" | Alice | alice@example.com | 2023-10-02 |

**The Problems (Update Anomalies):**
1. **Redundancy:** Alice's name and email are stored multiple times. Wasted disk space.
2. **Inconsistency Risk:** What if Alice changes her email? You have to update *every* row she ever posted. If you miss one, your data is corrupted (she has two different emails depending on which post you look at).

**The Normalized Solution (Our Schema):**
We split this into `users` and `posts`. We store Alice's email exactly *once* in the `users` table. The `posts` table only stores her `user_id`. When we need the email, we `JOIN` the tables.

## 3. Data Types and Space Efficiency

Choosing the right data type isn't just about correctness; it's about performance. Databases load data from slow HDDs/SSDs into fast RAM to query them. The smaller the row size, the more rows fit into RAM, meaning drastically faster queries.

- `VARCHAR(255)` vs `TEXT`: If a username will never exceed 50 characters, use `VARCHAR(50)`. `VARCHAR` calculates and allocates exactly the string length + 1 byte for storage. `TEXT` fields are stored separately from the main row data and carry overhead; only use them when content size is large and unpredictable (like a blog post body).
- `INT` vs `BIGINT`: An `INT` uses 4 bytes of storage. It can hold values up to ~2.1 billion. A `BIGINT` uses 8 bytes. Unless you expect more than 2 billion users, use `INT` for primary keys to save space.

## 4. Understanding Keys

### Primary Keys (PK)
A PK is the absolute truth for identifying a row. It typically has no business meaning (it's not an email or a username, which might change). We use surrogate keys (auto-incrementing integers) because they are fast to index and never change.

### Foreign Keys (FK) and Referential Integrity
When we defined `FOREIGN KEY (user_id) REFERENCES users(id)`, we told MySQL to act as a bouncer.

If application code tries to insert a post with `user_id = 999`, but no user with ID 999 exists, MySQL refuses the insertion and throws an error. This is **Referential Integrity**. It prevents "orphaned" records that point to thin air, which would cause application crashes later on.
