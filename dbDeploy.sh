#!/bin/bash

# Define variables
MYSQL_ROOT_PASSWORD="123456"
MYSQL_USER="apikiller"
MYSQL_PASSWORD="password"
MYSQL_DATABASE="apikiller"

# ANSI color codes
bold=$(tput bold)
normal=$(tput sgr0)
green=$(tput setaf 2)

# Pull mysql:5.7 image
docker pull mysql

# Start mysql container
docker run --name mysql-server -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} -p 3306:3306 -d mysql

# Wait for mysql service to start
until docker exec -it mysql-server mysqladmin ping --silent &> /dev/null; do
    echo "waiting for mysql to start ..."
    sleep 10
done


# Create new user and import data
docker cp ./apikiller.sql mysql-server:/tmp/apikiller.sql
docker exec -i mysql-server mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e "CREATE USER '${MYSQL_USER}'@'%' IDENTIFIED BY '${MYSQL_PASSWORD}'; GRANT ALL PRIVILEGES ON ${MYSQL_DATABASE}.* TO '${MYSQL_USER}'@'%'; FLUSH PRIVILEGES; source /tmp/apikiller.sql;"

echo "MySQL setup complete!"

# Print username and password in green color
echo "\r\n${bold}${green}\r\n======================================================================"
echo "Login MySQL server IP address (default returns info of the first network interface): `hostname -I | awk '{print $1}'`"
echo "Login MySQL server username: ${MYSQL_USER}"
echo "Login MySQL server password: ${MYSQL_PASSWORD}"
echo "Project database: ${MYSQL_DATABASE}${normal}"
echo "======================================================================"