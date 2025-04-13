package main

import (
	"api-service/internal/app"
	"api-service/internal/config"
	"flag"
	"fmt"
	"log"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "/etc/api-service/config.yaml", "path to config file")
	flag.Parse()

	cfg := config.NewConfig()
	err := cfg.Load(cfgPath)
	if err != nil {
		log.Fatalf("Error load config: %v", err)
		return
	}
	fmt.Println(cfg)

	a := app.NewApp()
	a.Run(*cfg)
}
