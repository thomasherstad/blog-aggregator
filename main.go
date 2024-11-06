package main

import (
	"fmt"
	"log"

	"github.com/thomasherstad/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	err = cfg.SetUser("Thomas")
	if err != nil {
		log.Printf("error setting user: %v", err)
	}

	fmt.Printf("DB url: %s\n", cfg.DBUrl)
	fmt.Printf("Current username: %s\n", cfg.CurrentUsername)
	fmt.Println(cfg)
}
