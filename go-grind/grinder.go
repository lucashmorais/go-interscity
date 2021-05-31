package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/lucashmorais/go-interscity/go-grind/resource_adaptor"
)

var NUM_REQUESTS int

func helloGoogle() {
	agent := fiber.Get("https://www.google.com")
	code, s, errors := agent.String()
	fmt.Println("Status code: ", code)
	fmt.Println("Erros: ", errors)
	fmt.Println(s)
}

func main() {
	// Load enviromental variables from config
	godotenv.Load("./config/config.env")

	NUM_REQUESTS = 2000
	NUM_PARALLEL_WORKERS := 4

	resource_adaptor.NUM_REQUESTS = NUM_REQUESTS

	var wg sync.WaitGroup

	ref_time := time.Now()

	for i := 0; i < NUM_PARALLEL_WORKERS; i++ {
		wg.Add(1)
		go resource_adaptor.GetResources(&wg)
	}

	wg.Wait()

	duration := time.Since(ref_time)
	println("Total time: ", duration.Seconds())
	println("Rate: ", float64(NUM_PARALLEL_WORKERS*NUM_REQUESTS)/float64(duration.Seconds()))
}
