#! /bin/bash

git pull
tag=$(date "+%Y%m%d%H%M%S")
docker build -t myapp:$tag .

docker rm -f myapp
docker run -d --name myapp -p 9101:8888 myapp:$tag