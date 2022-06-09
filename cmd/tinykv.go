package main

import (
	"flag"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/matrixorigin/tinykv/pkg/config"
	"github.com/matrixorigin/tinykv/pkg/server"
)

var (
	addr          = flag.String("addr", "127.0.0.1:8080", "tinykv server address")
	shardCapacity = flag.Uint64("shard-capacity", 1024*1024*64, "Data shard capaticy bytes")
	cfg           = flag.String("cfg", "./cube.toml", "cube toml config file")
)

func main() {
	flag.Parse()

	data, err := ioutil.ReadFile(*cfg)
	if err != nil {
		panic(err)
	}

	cfg := config.Config{
		Addr: *addr,
	}
	if _, err = toml.Decode(string(data), &cfg.CubeConfig); err != nil {
		panic(err)
	}

	cfg.Feature.ShardCapacityBytes = *shardCapacity
	svr := server.New(cfg)
	if err := svr.Start(); err != nil {
		panic(err)
	}
}
