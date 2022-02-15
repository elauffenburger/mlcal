#! /bin/bash

function log() {
    if [ -z $VERBOSE ]; then
        return
    fi

    echo $1
}

while getopts "re:c:p:v" flag; do
    case "${flag}" in
        r) RUN_ONLY=1;;
        e) ML_EMAIL=${OPTARG};;
        c) ML_COURSE=${OPTARG};;
        p) ML_PASSWORD=${OPTARG};;
        v) VERBOSE=1;;
    esac
done

if [ -z "${ML_EMAIL}" ]; then
    ML_EMAIL="elauffenburger@gmail.com"
fi

if [ -z "${ML_COURSE}" ]; then
    ML_COURSE="890"
fi

if [ -z "${ML_PASSWORD}" ]; then
    echo "-p (password) is required"
    exit 1
fi

if [ -z "${RUN_ONLY}" ]; then
    docker build -f Dockerfile.cli -t mlcalc-cli .
else
    log "-r provided; skipping build"
fi

docker run -it \
    --env ML_EMAIL=$ML_EMAIL \
    --env ML_PASSWORD=$ML_PASSWORD \
    --env ML_COURSE=$ML_COURSE \
    mlcalc-cli 
