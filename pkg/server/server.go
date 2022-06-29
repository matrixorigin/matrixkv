package server

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	cpebble "github.com/cockroachdb/pebble"
	"github.com/gin-gonic/gin"
	"github.com/lni/vfs"
	"github.com/matrixorigin/matrixcube/client"
	"github.com/matrixorigin/matrixcube/components/log"
	"github.com/matrixorigin/matrixcube/pb/rpcpb"
	"github.com/matrixorigin/matrixcube/raftstore"
	"github.com/matrixorigin/matrixcube/storage"
	"github.com/matrixorigin/matrixcube/storage/executor"
	"github.com/matrixorigin/matrixcube/storage/kv"
	"github.com/matrixorigin/matrixcube/storage/kv/pebble"
	"github.com/matrixorigin/matrixkv/pkg/config"
	"github.com/matrixorigin/matrixkv/pkg/metadata"
	"go.uber.org/zap"
)

var (
	defaultTimeout = time.Second * 30
)

// Server matrixkv server. The server support set, get and delete operation based on http.
type Server struct {
	cfg      config.Config
	eng      *gin.Engine
	client   client.Client
	kvClient client.KVClient
	store    raftstore.Store
}

// New create a tiny cube by config
func New(cfg config.Config) *Server {
	logger := log.GetDefaultZapLoggerWithLevel(zap.InfoLevel)

	// init logger
	cfg.CubeConfig.Logger = logger

	// init cube data storage
	// 1. create pebble as local db
	// 2. create executor to execute custom get/set/delete command
	// 3. create kv data storage
	// 4. setup datastorage to cube config
	kvs, err := pebble.NewStorage(filepath.Join(cfg.CubeConfig.DataPath, "kv-data"),
		logger, &cpebble.Options{})
	if err != nil {
		panic(err)
	}
	kvCommandExecutor := executor.NewKVExecutor(kvs)
	kvDataStorage := kv.NewKVDataStorage(kv.NewBaseStorage(kvs, vfs.Default),
		kvCommandExecutor,
		kv.WithFeature(cfg.Feature))

	// we only have a kv-based data storage
	cfg.CubeConfig.Storage.DataStorageFactory = func(group uint64) storage.DataStorage {
		return kvDataStorage
	}
	cfg.CubeConfig.Storage.ForeachDataStorageFunc = func(cb func(uint64, storage.DataStorage)) {
		cb(0, kvDataStorage)
	}

	store := raftstore.NewStore(&cfg.CubeConfig)
	store.Start()

	c := client.NewClient(client.Cfg{Store: store})
	kc := client.NewKVClient(c, 0, rpcpb.SelectLeader)
	return &Server{
		cfg:      cfg,
		eng:      gin.New(),
		store:    store,
		client:   c,
		kvClient: kc,
	}
}

// Start start a tiny kv server
func (s *Server) Start() error {
	if err := s.client.Start(); err != nil {
		return err
	}

	s.eng.POST("/set", s.handleSet())
	s.eng.POST("/delete", s.handleDelete())
	s.eng.GET("/get", s.handleGet())
	s.eng.GET("/shards", s.handleShards())

	return s.eng.Run(s.cfg.Addr)
}

func (s *Server) handleSet() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &metadata.SetRequest{}
		c.BindJSON(req)

		ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
		defer cancel()

		f := s.kvClient.Set(ctx, []byte(req.Key), []byte(req.Value))
		defer f.Close()

		err := f.GetError()
		resp := &metadata.SetResponse{
			Key: req.Key,
		}
		if err != nil {
			resp.Error = err.Error()
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (s *Server) handleDelete() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &metadata.DeleteRequest{}
		c.BindJSON(req)

		ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
		defer cancel()

		f := s.kvClient.Delete(ctx, []byte(req.Key))
		defer f.Close()

		err := f.GetError()
		resp := &metadata.DeleteResponse{
			Key: req.Key,
		}
		if err != nil {
			resp.Error = err.Error()
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (s *Server) handleGet() func(c *gin.Context) {
	return func(c *gin.Context) {
		key := c.Query("key")

		ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
		defer cancel()

		f := s.kvClient.Get(ctx, []byte(key))
		defer f.Close()

		r, err := f.GetKVGetResponse()
		resp := &metadata.GetResponse{
			Key: key,
		}
		if err != nil {
			resp.Error = err.Error()
		} else {
			resp.Value = string(r.Value)
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (s *Server) handleShards() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := uint64(0)
		local := c.Query("local")
		if len(local) > 0 {
			id = s.store.Meta().ID
		}

		var shards []raftstore.Shard
		s.store.GetRouter().AscendRangeWithoutSelectReplica(0, nil, nil, func(shard raftstore.Shard) bool {
			if id == 0 {
				shards = append(shards, shard)
			} else {
				for _, r := range shard.Replicas {
					if r.StoreID == id {
						shards = append(shards, shard)
					}
				}
			}
			return true
		})

		c.JSON(http.StatusOK, shards)
	}
}
