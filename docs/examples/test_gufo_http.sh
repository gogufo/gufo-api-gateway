#!/usr/bin/env bash
set -e

echo ">>> Starting Gufo + test-ms stack..."
docker compose up -d

echo ">>> Waiting for services to warm up..."
sleep 5

echo ">>> Checking Gufo health..."
curl -sf http://localhost:8090/api/v1/health || {
  echo "Gufo health failed"
  exit 1
}

echo ""
echo ">>> Calling test microservice via Gufo:"
URL="http://localhost:8090/api/v1/test/info"

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}\n" "$URL")
echo "$RESPONSE"

HTTP_CODE=$(echo "$RESPONSE" | grep HTTP_CODE | cut -d: -f2)

if [ "$HTTP_CODE" != "200" ]; then
  echo "TEST FAILED: non-200 status code"
  exit 1
fi

echo ">>> TEST PASSED"
