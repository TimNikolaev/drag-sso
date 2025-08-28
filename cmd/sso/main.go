package main

import (
	"fmt"

	"github.com/TimNikolaev/drag-sso/internal/config"
)

func main() {
	// Init config
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// Init logger

	// Init app

	// Run app's gRPC-server
}
