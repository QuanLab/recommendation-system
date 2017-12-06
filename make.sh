#!/usr/bin/env bash
VERSION=$(date +"%y.%m.%d")-dev
WORKING_DIR=$(pwd)
DOCKER_IMAGE=go-recommend

echo 'compile application...'
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

echo 'Build docker image'
docker build -t quanpv/$DOCKER_IMAGE:$VERSION .

echo 'Login to Docker hub'
# docker login hub.docker.com

echo 'Push to Docker hub'
docker push  quanpv/$DOCKER_IMAGE:$VERSION