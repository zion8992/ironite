#!/bin/bash

CONTAINER_NAME="cubyzListDB"
MYSQL_ROOT_PASSWORD="H0EeLfLnO,xDEVELOPERSx4c!#%"

# Check if container is running
if [ ! "$(docker ps -q -f name=^/${CONTAINER_NAME}$)" ]; then
    echo "Container is not running. Start it first."
    exit 1
fi

# Connect to MySQL root shell
docker exec -it $CONTAINER_NAME mysql -u root -p$MYSQL_ROOT_PASSWORD
