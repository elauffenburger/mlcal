#!/bin/bash
set -e

if [ -z "${ML_EMAIL}" ]; then
    ML_EMAIL="elauffenburger@gmail.com"
fi

if [ -z "${ML_COURSE}" ]; then
    ML_COURSE="890"
fi

if [ -z "${ML_REFRESH_INTERVAL}" ]; then
    ML_REFRESH_INTERVAL="1m"
fi

if [ -z "${ML_PASSWORD}" ]; then
    echo "ML_PASSWORD is required"
    exit 1
fi

if [ -z "${RUN_ONLY}" ]; then
    docker build -f Dockerfile.api -t mlcal-api .
else
    log "-r provided; skipping build"
fi

if [ -n "${REDIS_HOST}" ]; then
    REDIS_HOST="mlcal-api-redis"
fi

if [ -n "${REDIS_PORT}" ]; then
    REDIS_PORT=6379
fi

EXTRA_API_DOCKER_ARGS=''
if [ -n "${USE_LOCAL_VOLUME}" ]; then
    EXTRA_API_DOCKER_ARGS="$EXTRA_API_DOCKER_ARGS -v $(pwd):/api"
    EXTRA_API_DOCKER_ARGS="$EXTRA_API_DOCKER_ARGS --entrypoint sh"
fi

# Create a network for services to share.
if ! docker network ls | grep mlcal-api-net; then
    docker network create mlcal-api-net
fi

# Clean up the existing redis container.
if docker ps | grep mlcal-api-redis; then
    docker stop mlcal-api-redis
    docker rm mlcal-api-redis
fi

# Start up redis.
docker run -d --name mlcal-api-redis --network mlcal-api-net redis

# Start up the api.
docker run -it --rm --network mlcal-api-net \
    --env ML_EMAIL=$ML_EMAIL \
    --env ML_PASSWORD=$ML_PASSWORD \
    --env ML_COURSE=$ML_COURSE \
    --env ML_REFRESH_INTERVAL=$ML_REFRESH_INTERVAL \
    --env REDIS_HOST=$REDIS_HOST \
    --env REDIS_PORT=$REDIS_PORT \
    --name mlcal-api \
    $EXTRA_API_DOCKER_ARGS \
    mlcal-api

# Clean up.
docker rm mlcal-api-redis