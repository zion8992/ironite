#!/bin/bash

CONTAINER_NAME="cubyzListDB"
VOLUME_NAME="cubyzListDB"
MYSQL_ROOT_PASSWORD="H0EeLfLnO,xDEVELOPERSx4c!#%"

# Check if container already exists
if [ "$(docker ps -a -q -f name=^/${CONTAINER_NAME}$)" ]; then
    echo "Container exists. Starting it..."
    docker start $CONTAINER_NAME
else
    echo "Creating volume (if not exists)..."
    docker volume create $VOLUME_NAME

    echo "Creating and starting MySQL container..."
    docker run -d \
        --name $CONTAINER_NAME \
        -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
        -v $VOLUME_NAME:/var/lib/mysql \
        -p 3306:3306 \
        mysql:8.0
fi
