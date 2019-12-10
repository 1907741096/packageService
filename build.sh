#! /bin/bash

env=$1
if [ -z "${env}" ]; then
    env=debug
fi

ENV_ARR=(debug test release)
if ! echo "${ENV_ARR[@]}" | grep -w ${env} &>/dev/null; then
    echo "please select one in env options:" "${ENV_ARR[@]}"
    exit
fi

#if [ "$env" = release ]; then
#    export GOPATH=/data/www/wwwroot/go
#    cd /data/www/wwwroot/go/src/moon
#fi

output=$2
if [ -z ${output} ]; then
    output=as
fi

if [ "$env" = release -o "$env" = test ]; then
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${output} --tags "${env}" main.go
else
    go build -o ${output} --tags "${env}" main.go
fi
