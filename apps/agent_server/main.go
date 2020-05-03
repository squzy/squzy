package main

import (
	"log"
	"squzy/apps/agent_server/application"
	"squzy/apps/agent_server/server"
)

func main() {
	app := application.New(server.New())
	log.Fatal(app.Run(10001))
}
