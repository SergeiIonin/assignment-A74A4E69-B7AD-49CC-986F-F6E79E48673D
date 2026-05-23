#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

IMAGE="assignment-a74a4e69-b7ad-49cc-986f-f6e79e48673d"
CONF="${1:-$ROOT_DIR/conf.yaml}"
CONTAINER_ID=""

cleanup() {
  if [ -n "$CONTAINER_ID" ]; then
    echo "Shutting down container $CONTAINER_ID"
    docker stop "$CONTAINER_ID"
  fi
}
trap cleanup EXIT

# Rebuild if image doesn't exist or there are non-README changes
IMAGE_EXISTS=$(docker image inspect "$IMAGE" > /dev/null 2>&1 && echo "yes" || echo "no")
GIT_CHANGES=$(git -C "$ROOT_DIR" diff HEAD -- . ':(exclude)README.md' 2>/dev/null)

if [ "$IMAGE_EXISTS" = "no" ] || [ -n "$GIT_CHANGES" ]; then
  echo "Rebuilding image..."
  docker build -t "$IMAGE" "$ROOT_DIR"
else
  echo "No changes detected, skipping rebuild"
fi

LOG_FILE="$ROOT_DIR/server_$(date '+%Y-%m-%d_%H-%M-%S').log"
CONTAINER_ID=$(docker run -d --rm \
  -p 8080:8080 \
  -v "$CONF":/app/conf.yaml:ro \
  "$IMAGE")
echo "Container started: $CONTAINER_ID, logs -> $LOG_FILE"
docker logs -f "$CONTAINER_ID" > "$LOG_FILE" 2>&1 &

sleep 2

echo "Testing dashboard endpoint for user id 1"
curl http://localhost:8080/dashboard/1 | jq .
echo "--------------"
echo "Testing dashboard endpoint for user id 2"
curl http://localhost:8080/dashboard/2 | jq .
echo "--------------"
echo "Testing dashboard endpoint for user id 3"
curl http://localhost:8080/dashboard/3 | jq .
echo "--------------"
