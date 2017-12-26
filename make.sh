#!/usr/bin/env bash
VERSION=$(date +"%y.%m.%d")-dev
VERSION=1.0.0
WORKING_DIR=$(pwd)
DOCKER_IMAGE=go-recommender

echo 'compile application...'
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

echo 'Build docker image'
docker build -t quanpv/$DOCKER_IMAGE:$VERSION .

echo 'Login to Docker hub'
# docker login hub.docker.com

echo 'Push to Docker hub'
docker push  quanpv/$DOCKER_IMAGE:$VERSION