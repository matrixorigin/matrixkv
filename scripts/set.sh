#!/bin/bash
key=$1
value=$2
json="{\"key\":\"$key\", \"value\":\"$value\"}"
echo "set $key=$value to http://127.0.0.1:808$3"
http_proxy="" curl -H "Content-type: application/json" -X POST -d "$json" http://127.0.0.1:808$3/set
echo ""
