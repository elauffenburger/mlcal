#!/usr/bin/env bash

args=(
    "--email $ML_EMAIL"
    "--password $ML_PASSWORD"
    "--course $ML_COURSE"
    "--refresh $ML_REFRESH_INTERVAL"
)

./api ${args[@]}