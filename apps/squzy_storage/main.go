package main

import (
	"squzy/apps/squzy_storage/config"
	"squzy/apps/squzy_storage/executable"
	_ "squzy/apps/squzy_storage/version"
)

// When we're running binary, we want to use config from OS env
func main() {
	cfg := config.New()
	executable.Execute(cfg)
}