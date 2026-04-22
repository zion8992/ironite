#!/bin/bash

CONTAINER_NAME="cubyzListDB"

if [ "$(docker ps -q -f name=^/${CONTAINER_NAME}$)" ]; then
    echo "Stopping container..."
    docker stop $CONTAINER_NAME
else
    echo "Container is not running."
fi
