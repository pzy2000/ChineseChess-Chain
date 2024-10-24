#!/bin/sh

# Start the MySQL container
#docker-compose up -d

# Wait for the MySQL container to be ready
#echo "Waiting for MySQL container to be ready..."
#sleep 10

# Insert data into the MySQL container
echo "Inserting data into the MySQL container..."
cat insert_pg.sql | PGPASSWORD=123456 psql -h 192.168.3.170 -U SYSTEM -p 54321 -d explorer_format_test

echo "Data insertion complete."