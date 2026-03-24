# 🗺 MySQL Learning Roadmap: Social Media Backend

Welcome to your databases journey! This guide is designed to be a structured refresh using only the terminal, building a conceptual **Social Media App** database. 

## Phase 1: Environment Setup
1. **Install MySQL**: We will use Homebrew (`brew install mysql`).
2. **Start the Service**: Learn how to start and stop the database engine in the background.
3. **Connect via CLI**: Log in to the MySQL monitor (`mysql -u root`).

## Phase 2: Database & Schema Foundations
1. **Create the Database**: `CREATE DATABASE social_app;`
2. **Design the Schema**:
   - `users`: id, username, email, created_at.
   - `posts`: id, user_id, content, created_at.
   - `followers`: follower_id, followee_id, created_at.
3. **Create Tables**: Learn about Data Types (VARCHAR, INT, TIMESTAMP), Primary Keys, and Foreign Keys.

## Phase 3: Writing Data (CRUD Basics)
1. **Insert (C)**: Add dummy users, posts, and follower relationships.
2. **Select (R)**: Fetch specific users or posts.
3. **Update (U)**: Change a user's email or a post's text.
4. **Delete (D)**: Remove a follower or a post.

## Phase 4: Reading and Assembling the "Feed"
1. **JOINs**: Combine `users` and `posts` to see *who* posted *what*.
2. **Aggregations**: Count the number of followers a user has (`GROUP BY`, `COUNT`).

## Phase 5: Transactions (Advanced)
A transaction ensures that multiple operations succeed or fail together (ACID properties). 
- **Scenario**: A user deletes their account. We must delete their posts, their followers, and their profile in one atomic step.
- **Commands**: `START TRANSACTION`, `COMMIT`, `ROLLBACK`.
