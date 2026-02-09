package main

import (
	"projectwebcurhat/config"
	"projectwebcurhat/config/server"
)

func main() {
	config.Load()
	server.Run()
}

