#!/bin/bash

# MySQL 容器名称
CONTAINER_NAME="docker_cm_explorer_db_1"

# 数据库连接信息
DB_ROOT_USER="root"
DB_ROOT_PASSWORD="chainmaker"

# 要创建的数据库和账户信息
DATABASES=("chainmaker_dquery" "chainmaker_explorer_test")
READONLY_USER="readonly"
READONLY_PASSWORD="readonly123"
CHAINMAKER_USER="chainmaker"
CHAINMAKER_PASSWORD="chainmaker123"

# 等待数据库服务启动
until docker exec $CONTAINER_NAME mysql -u"$DB_ROOT_USER" -p"$DB_ROOT_PASSWORD" -e "SELECT 1;" &> /dev/null; do
    echo "Waiting for database connection..."
    sleep 5
done

echo "Connected to database."

# 创建数据库和用户
for DB in "${DATABASES[@]}"; do
    echo "Creating database $DB if not exists..."
    docker exec $CONTAINER_NAME mysql -u"$DB_ROOT_USER" -p"$DB_ROOT_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS $DB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" || {
        echo "Failed to create database $DB"
        exit 1
    }
    
    echo "Creating readonly user for $DB..."
    docker exec $CONTAINER_NAME mysql -u"$DB_ROOT_USER" -p"$DB_ROOT_PASSWORD" -e "CREATE USER IF NOT EXISTS '$READONLY_USER'@'%' IDENTIFIED BY '$READONLY_PASSWORD'; GRANT SELECT ON $DB.* TO '$READONLY_USER'@'%';" || {
        echo "Failed to create readonly user for $DB"
        exit 1
    }
    
    echo "Creating chainmaker user for $DB..."
    docker exec $CONTAINER_NAME mysql -u"$DB_ROOT_USER" -p"$DB_ROOT_PASSWORD" -e "CREATE USER IF NOT EXISTS '$CHAINMAKER_USER'@'%' IDENTIFIED BY '$CHAINMAKER_PASSWORD'; GRANT ALL PRIVILEGES ON $DB.* TO '$CHAINMAKER_USER'@'%';" || {
        echo "Failed to create chainmaker user for $DB"
        exit 1
    }
done

# 确保全局权限
echo "Granting global privileges to chainmaker user..."
docker exec $CONTAINER_NAME mysql -u"$DB_ROOT_USER" -p"$DB_ROOT_PASSWORD" -e "GRANT ALL PRIVILEGES ON *.* TO '$CHAINMAKER_USER'@'%'; FLUSH PRIVILEGES;" || {
    echo "Failed to grant global privileges"
    exit 1
}

echo "Database initialization complete."