# Phase 4: Reading and Assembling the "Feed"

In the real world, data is rarely retrieved from a single table. To display a "feed", we need the post content, the author's username, and perhaps how many followers they have. We accomplish this through **JOINs** and **Aggregations**.

## 1. JOIN Operations

A `JOIN` clause combines rows from two or more tables based on a related column between them (usually a Primary Key to Foreign Key relationship).

### The Essential JOIN: INNER JOIN
We want to see the `username` of the person who wrote the post, alongside the post `content`. The raw `posts` table only has the `user_id`. Let's join them!

```sql
SELECT users.username, posts.content, posts.created_at
FROM posts
JOIN users ON posts.user_id = users.id
ORDER BY posts.created_at DESC;
```

**Using Aliases:**
To type less, you can give tables temporary aliases using `AS`.

```sql
SELECT u.username, p.content, p.created_at
FROM posts AS p
JOIN users AS u ON p.user_id = u.id
ORDER BY p.created_at DESC;
```

### LEFT JOIN, RIGHT JOIN, and FULL OUTER JOIN
- `INNER JOIN`: Only returns rows if there is a match in *both* tables. (If a user has no posts, they won't show up).
- `LEFT JOIN`: Returns *all* rows from the left table (`FROM users`), even if there are no matches in the right table (`JOIN posts`). Unmatched posts will be filled with `NULL`.
- `RIGHT JOIN`: The exact opposite of LEFT JOIN. Returns all rows from the right table, and fills unmatched columns from the left table with `NULL`.
- `FULL OUTER JOIN`: Returns all rows when there is a match in either left or right table. MySQL doesn't have a native `FULL OUTER JOIN` keyword, but we can simulate it using `UNION`.

```sql
-- LEFT JOIN: Returns all users, along with their posts if they have any.
SELECT u.username, p.content
FROM users u
LEFT JOIN posts p ON u.id = p.user_id;
```

```sql
-- RIGHT JOIN: Returns all posts, ensuring no child record is missed even if the parent user was somehow deleted (though our CASCADE constraint prevents this).
SELECT u.username, p.content
FROM users u
RIGHT JOIN posts p ON u.id = p.user_id;
```

```sql
-- Emulating FULL OUTER JOIN using UNION:
-- We combine the results of a LEFT JOIN and a RIGHT JOIN.
SELECT u.username, p.content
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
UNION
SELECT u.username, p.content
FROM users u
RIGHT JOIN posts p ON u.id = p.user_id;
```

## 2. Aggregations (GROUP BY & Functions)

Aggregations allow you to calculate summary statistics like `COUNT`, `SUM`, `AVG`, `MAX`, and `MIN`.

**Counting Total Users:**
```sql
SELECT COUNT(*) FROM users;
```

**Counting Posts per User using `GROUP BY`:**
```sql
SELECT u.username, COUNT(p.id) AS post_count
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
GROUP BY u.username;
```

**Who has the most followers?**
We want to count how many times a user's ID appears as a `followee_id` in the `followers` table.

```sql
SELECT u.username, COUNT(f.follower_id) AS follower_count
FROM users u
LEFT JOIN followers f ON u.id = f.followee_id
GROUP BY u.username
ORDER BY follower_count DESC;
```

## Challenge (The Ultimate "Feed" Query):
Can you write a query that shows all recent posts, but ONLY from the people a specific user (let's say Charlie, ID 3) follows?

```sql
SELECT p.content, u.username AS author
FROM posts p
JOIN users u ON p.user_id = u.id
JOIN followers f ON u.id = f.followee_id
WHERE f.follower_id = 3
ORDER BY p.created_at DESC;
```
*Think through this query: We start with posts, join users to get the author's name, then join followers to see if the reader (Charlie, ID 3) is following that author.*

## Complex Query Example: Derived Tables (Subqueries in the FROM clause)
Sometimes we need to aggregate data first, and *then* join it. 
Goal: Get a list of all users, but instead of counting posts on the fly, use a subquery to find active authors (users with more than 1 post) and join that back to the main users list.

```sql
SELECT u.username, active_authors.post_count
FROM users u
JOIN (
    -- This inner query acts as a temporary "Derived Table"
    SELECT user_id, COUNT(id) as post_count
    FROM posts
    GROUP BY user_id
    HAVING COUNT(id) > 1
) AS active_authors ON u.id = active_authors.user_id;
```
*Note the `HAVING` clause. We use `WHERE` to filter raw rows, but we use `HAVING` to filter aggregated data (like checking if a COUNT is greater than 1).*
