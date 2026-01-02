package main

import (
	"log"
	"os"
	"time"
)

func main() {
	cfg := config{
		addr:         ":8000",
		db:           dbConfig{},
		writeTimeout: time.Second * 30,
	}

	api := application{
		config: cfg,
	}

	if err := api.run(api.mount()); err != nil {
		log.Printf("server has failed to start, err %s", err)
		os.Exit(1)
	}
}
