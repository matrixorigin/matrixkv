#!/bin/bash

curl -H "Content-type: application/json" -X POST -d '{"key":"k1", "value":"v1"}' http://127.0.0.1:8080/set