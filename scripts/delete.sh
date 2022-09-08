#!/bin/bash
key=$1
json="{\"key\":\"$key\"}"
echo "set $key=$value to http://127.0.0.1:808$2"
http_proxy="" curl -H "Content-type: application/json" -X POST -d "$json" http://127.0.0.1:808$2/delete
echo ""
