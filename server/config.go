package server

import "github.com/rs/zerolog"

type Config struct {
	// Port is the port the EStatz server should listen on.
	Port int
	// Workers is the amount of workers that should be spawned initally to handle packets.
	Workers int
	// Logger is the logger used for the EStatz server.
	Logger zerolog.Logger
}
