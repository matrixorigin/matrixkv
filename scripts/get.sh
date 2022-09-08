#!/bin/bash
key=$1
echo "get $key from http://127.0.0.1:808$2"
http_proxy="" curl "http://127.0.0.1:808$2/get?key=$key"
echo ""
