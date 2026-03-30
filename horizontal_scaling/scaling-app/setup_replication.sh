#!/bin/bash
set -e

echo "🧹 1. Cleaning up existing Docker database states..."
mysql -h 127.0.0.1 -P 3306 -u root -prootpassword -e "DROP DATABASE IF EXISTS social_app; RESET MASTER;"
mysql -h 127.0.0.1 -P 3307 -u root -prootpassword -e "STOP SLAVE; RESET SLAVE ALL; DROP DATABASE IF EXISTS social_app;"
mysql -h 127.0.0.1 -P 3308 -u root -prootpassword -e "STOP SLAVE; RESET SLAVE ALL; DROP DATABASE IF EXISTS social_app;"

echo "🚀 2. Setting up Replication User on Master..."
mysql -h 127.0.0.1 -P 3306 -u root -prootpassword -e "
CREATE USER IF NOT EXISTS 'replica_user'@'%' IDENTIFIED BY 'replica_password';
GRANT REPLICATION SLAVE ON *.* TO 'replica_user'@'%';
FLUSH PRIVILEGES;
"

echo "📖 3. Fetching Fresh Master Status..."
MASTER_STATUS=$(mysql -h 127.0.0.1 -P 3306 -u root -prootpassword -e "SHOW MASTER STATUS;" | tail -n 1)
FILE=$(echo "$MASTER_STATUS" | awk '{print $1}')
POS=$(echo "$MASTER_STATUS" | awk '{print $2}')
echo "✅ Master is cleanly at File: $FILE, Position: $POS"

echo "🔗 4. Linking Replicas to Master..."
for PORT in 3307 3308; do
  mysql -h 127.0.0.1 -P $PORT -u root -prootpassword -e "
  CHANGE MASTER TO 
    MASTER_HOST='db-master',
    MASTER_USER='replica_user',
    MASTER_PASSWORD='replica_password',
    MASTER_LOG_FILE='$FILE', 
    MASTER_LOG_POS=$POS;
  START SLAVE;
  "
done

echo "🏗️ 5. Creating Schema on Master (Replicas will auto-copy this!)..."
mysql -h 127.0.0.1 -P 3306 -u root -prootpassword -e "
CREATE DATABASE social_app;
CREATE TABLE social_app.users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);
"

echo "✅ True Master-Slave Replication has been entirely initialized!"
