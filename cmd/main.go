package main

func main() {
	cfg := config{
		addr: ":8000",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}
}
