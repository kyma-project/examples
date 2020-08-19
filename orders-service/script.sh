#!/usr/bin/env bash

curl -v -X POST -H "Content-Type: application/json" https://orders-webhook.34.90.177.145.xip.io/orders -k -d \
  '{
    "consignmentCode": "76272727",
    "orderCode": "76272725",
    "consignmentStatus": "PICKUP_COMPLETE"
  }' 