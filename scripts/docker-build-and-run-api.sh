#!/bin/bash

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

. "$(dirname $0)/./docker-run-api.sh"