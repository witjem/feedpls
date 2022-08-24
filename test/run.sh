#!/bin/bash
set -eu

ROOT_DIR=$(git rev-parse --show-toplevel)

stop_feedpls() {
  echo "Stopping feedpls instance"
  docker compose -f $ROOT_DIR/examples/docker-compose.yml down -v
}

wait_for_url () {
    echo "Testing $1"
    max_in_s=$2
    delay_in_s=1
    total_in_s=0
    while [ $total_in_s -le "$max_in_s" ]
    do
        echo "Wait ${total_in_s}s"
        if (echo -e "GET $1\nHTTP/* 200" | hurl > /dev/null 2>&1;) then
            return 0
        fi
        total_in_s=$(( total_in_s +  delay_in_s))
        sleep $delay_in_s
    done
    return 1
}

trap 'stop_feedpls' ERR EXIT

echo "Building docker image: witjem/feedpls:main"
docker build --tag=witjem/feedpls:main --progress=plain $ROOT_DIR

echo "Starting feedpls container"
docker compose -f $ROOT_DIR/examples/docker-compose.yml up -d

echo "Starting feedpls instance to be ready"
wait_for_url 'http://localhost:8080/ping' 60

echo "Running Hurl tests"
hurl $ROOT_DIR/test/*.hurl --test
