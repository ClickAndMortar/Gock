#!/bin/bash

set -e

test -f .env && {
  export $(egrep -v '^#' .env | xargs)
}

docker build -t ${IMAGE}:${TAG} .
docker push ${IMAGE}:${TAG}
