package main

import (
	"fmt"
	"log"

	"github.com/Pizzu/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("Luca")

	if err != nil {
		log.Fatalf("error writing from config: %v", err)
	}

	cfg, err = config.Read()

	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg)

}
