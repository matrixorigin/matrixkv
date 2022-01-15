package config

import (
	cube "github.com/matrixorigin/matrixcube/config"
)

// Config config
type Config struct {
	// Addr http api addr
	Addr string `toml:"addr"`
	// CubeConfig cube config
	CubeConfig cube.Config `toml:"cube"`
}
