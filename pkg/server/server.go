package server

import (
	"net/http"
	"path/filepath"
	"time"

	cpebble "github.com/cockroachdb/pebble"
	"github.com/gin-gonic/gin"
	"github.com/lni/vfs"
	"github.com/matrixorigin/matrixcube/components/log"
	"github.com/matrixorigin/matrixcube/raftstore"
	cube "github.com/matrixorigin/matrixcube/server"
	"github.com/matrixorigin/matrixcube/storage"
	cubeKV "github.com/matrixorigin/matrixcube/storage/kv"
	"github.com/matrixorigin/matrixcube/storage/kv/pebble"
	"github.com/matrixorigin/tinykv/pkg/config"
	"github.com/matrixorigin/tinykv/pkg/kv"
	"github.com/matrixorigin/tinykv/pkg/metadata"
	"go.uber.org/zap"
)

var (
	defaultTimeout = time.Second * 10
)

// Server tinykv server. The server support set, get and delete operation based on http.
type Server struct {
	cfg config.Config
	eng *gin.Engine
	app *cube.Application
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
	kvs, err := pebble.NewStorage(filepath.Join(cfg.CubeConfig.DataPath, "kv-data"), logger, &cpebble.Options{})
	if err != nil {
		panic(err)
	}
	kvCommandExecutor := kv.NewSimpleKVExecutor(kvs)
	kvDataStorage := cubeKV.NewKVDataStorage(cubeKV.NewBaseStorage(kvs, vfs.Default), kvCommandExecutor)

	// we only have a kv-based data storage
	cfg.CubeConfig.Storage.DataStorageFactory = func(group uint64) storage.DataStorage {
		return kvDataStorage
	}
	cfg.CubeConfig.Storage.ForeachDataStorageFunc = func(cb func(storage.DataStorage)) {
		cb(kvDataStorage)
	}

	return &Server{
		cfg: cfg,
		eng: gin.Default(),
		app: cube.NewApplication(cube.Cfg{
			Store: raftstore.NewStore(&cfg.CubeConfig),
		}),
	}
}

// Start start a tiny kv server
func (s *Server) Start() error {
	if err := s.app.Start(); err != nil {
		return err
	}

	s.eng.POST("/set", s.handleSet())
	s.eng.POST("/delete", s.handleDelete())
	s.eng.GET("/get", s.handleGet())

	return s.eng.Run(s.cfg.Addr)
}

func (s *Server) handleSet() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &metadata.SetRequest{}
		c.BindJSON(req)

		_, err := s.app.Exec(cube.CustomRequest{
			Key:        []byte(req.Key),
			Cmd:        []byte(req.Value),
			CustomType: metadata.SetType,
			Write:      true, // write is write operation
		}, defaultTimeout)

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

		_, err := s.app.Exec(cube.CustomRequest{
			Key:        []byte(req.Key),
			CustomType: metadata.DeleteType,
			Write:      true, // delete is write operation
		}, defaultTimeout)

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
		value, err := s.app.Exec(cube.CustomRequest{
			Key:        []byte(key),
			CustomType: metadata.GetType,
			Read:       true, // get is read operation
		}, defaultTimeout)

		resp := &metadata.GetResponse{
			Key: key,
		}
		if err != nil {
			resp.Error = err.Error()
		} else {
			resp.Value = string(value)
		}

		c.JSON(http.StatusOK, resp)
	}
}
