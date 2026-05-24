#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

SERVER_PID=""

cleanup() {
  if [ -n "$SERVER_PID" ]; then
    echo "Shutting down server (PID $SERVER_PID)"
    kill -- -"$SERVER_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT

set -m  # job control: background processes get their own process group
LOG_FILE="$ROOT_DIR/server_$(date '+%Y-%m-%d_%H-%M-%S').log"
cd "$ROOT_DIR" && go run cmd/main.go > "$LOG_FILE" 2>&1 &
SERVER_PID=$!
echo "Server started (PID $SERVER_PID), logs -> $LOG_FILE"

echo "Waiting for server to be ready..."
RETRIES=20
until curl -s http://localhost:8080/dashboard/0 > /dev/null 2>&1; do
  RETRIES=$((RETRIES - 1))
  if [ "$RETRIES" -eq 0 ]; then
    echo "Server did not become ready in time"
    exit 1
  fi
  sleep 0.5
done
echo "Server is ready"

for id in {1..10}; do
    echo "Testing dashboard endpoint for user id $id"
    curl http://localhost:8080/dashboard/$id | jq .
    echo "--------------"
done

echo "Calculating average response time"
sleep 0.5  # let the log flush
source "$SCRIPT_DIR/calculateAvgResponse.sh" "$LOG_FILE"
