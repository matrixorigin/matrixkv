package config

import (
	cube "github.com/matrixorigin/matrixcube/config"
	"github.com/matrixorigin/matrixcube/storage"
)

// Config config
type Config struct {
	// Addr http api addr
	Addr string `toml:"addr"`
	// CubeConfig cube config
	CubeConfig cube.Config `toml:"cube"`
	// Feature storage feature
	Feature storage.Feature
}
