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
    echo "-p (password) is required"
    exit 1
fi

if [ -z "${RUN_ONLY}" ]; then
    docker build -f Dockerfile.api -t mlcalc-api .
else
    log "-r provided; skipping build"
fi

docker run -it --rm \
    --env ML_EMAIL=$ML_EMAIL \
    --env ML_PASSWORD=$ML_PASSWORD \
    --env ML_COURSE=$ML_COURSE \
    --env ML_REFRESH_INTERVAL=$ML_REFRESH_INTERVAL \
    --name mlcalc-api \
    mlcalc-api