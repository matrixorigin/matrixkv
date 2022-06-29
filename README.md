# Matrixkv
This is a distributed kv project to demonstrate how to use matrixcube.

### Build docker image
```bash
make docker
```

### Prepare local data dir
```bash
mkdir -p /data/node0
mkdir -p /data/node1
mkdir -p /data/node2
mkdir -p /data/node3
```

### Run with docker-compose
```bash
docker-compose up
```

### Set on any node
```bash
curl -X POST  -H 'Content-Type: application/json' -d '{"key":"k1","value":"v1"}' http://127.0.0.1:8080/set

curl -X POST  -H 'Content-Type: application/json' -d '{"key":"k2","value":"v2"}' http://127.0.0.1:8081/set

curl -X POST  -H 'Content-Type: application/json' -d '{"key":"k3","value":"v3"}' http://127.0.0.1:8082/set
```


### Get on any node
```bash
curl http://127.0.0.1:8080/get?key=k1
curl http://127.0.0.1:8081/get?key=k1
curl http://127.0.0.1:8082/get?key=k1

curl http://127.0.0.1:8080/get?key=k2
curl http://127.0.0.1:8081/get?key=k2
curl http://127.0.0.1:8082/get?key=k2

curl http://127.0.0.1:8080/get?key=k3
curl http://127.0.0.1:8081/get?key=k3
curl http://127.0.0.1:8082/get?key=k3
```

### Delete on any node
```bash
curl -X POST  -H 'Content-Type: application/json' -d '{"key":"k1"}' http://127.0.0.1:8080/delete

curl -X POST  -H 'Content-Type: application/json' -d '{"key":"k2"}' http://127.0.0.1:8081/delete

curl -X POST  -H 'Content-Type: application/json' -d '{"key":"k3"}' http://127.0.0.1:8082/delete
```

