#!/usr/bin/env bash

export AWS_ACCESS_KEY_ID=foo
export AWS_SECRET_ACCESS_KEY=bar
export AWS_DEFAULT_REGION=us-east-1

docker compose up -d && \
aws --endpoint=http://127.0.0.1:5000 \
  ssm put-parameter \
    --name 'testtest' \
    --type 'String' \
    --value '12345678901234567890'
