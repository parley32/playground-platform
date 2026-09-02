#!/usr/bin/env bash

curl -X POST \
  -H "Content-Type: application/json" \
  -d "{\"description\":\"$*\"}" \
  http://localhost:8080/task/blbla/