#!/bin/bash
echo "get $1 shards info from http://127.0.0.1:808$2"
http_proxy="" curl "http://127.0.0.1:808$2/shards?$1=true"
echo ""
