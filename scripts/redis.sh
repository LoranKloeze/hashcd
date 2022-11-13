#!/bin/bash

echo 'Redis will listen on port 8910'
docker build -t redis-big ./redis/.
docker run --name finalcd-redis -p 127.0.0.1:8910:6379 redis-big