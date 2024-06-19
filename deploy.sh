#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

CMD=$1

APP_NAME="${APP_NAME}"

### docker hub
export DOCKER_HUB="hub.oneflow.dev"
export DOCKER_IMAGE="$DOCKER_HUB/app/${APP_NAME}"

# decide App/Service version
VERSION="${GO_PIPELINE_LABEL:default}"

if [ "${CMD}" = "build-image" ]; then
  echo ">>>build push image ${DOCKER_IMAGE}:${VERSION}"
  cat Dockerfile
  docker build -t "$DOCKER_IMAGE:$VERSION" .
#  docker push "$DOCKER_IMAGE:$VERSION"
#  docker rmi "$DOCKER_IMAGE:$VERSION"
elif [ "${CMD}" = "deploy" ]; then
  docker rm -f ${APP_NAME}
  docker run -it --name ${APP_NAME} -p 9101:8888 -d ${DOCKER_IMAGE}:${VERSION}
fi