#!/bin/bash
set -e

function log() {
    if [ -z $VERBOSE ]; then
        return
    fi

    echo $1
}

while getopts "re:c:p:v:l" flag; do
    case "${flag}" in
        r) RUN_ONLY=1;;
        e) ML_EMAIL=${OPTARG};;
        c) ML_COURSE=${OPTARG};;
        p) ML_PASSWORD=${OPTARG};;
        v) VERBOSE=1;;
        l) USE_LOCAL_VOLUME=1;;
    esac
done

SCRIPT_DIR=$(dirname $0)

. "${SCRIPT_DIR}/docker-run-api.sh"