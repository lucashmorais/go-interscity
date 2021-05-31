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
	fmt.Println("Errors: ", errors)
	fmt.Println(s)
}

func coreTester(numRequests int, numParallelWorkers int, goRoutine func(*sync.WaitGroup, interface{}), extras interface{}) {
	NUM_REQUESTS = numRequests
	NUM_PARALLEL_WORKERS := numParallelWorkers

	resource_adaptor.NUM_REQUESTS = NUM_REQUESTS

	var wg sync.WaitGroup

	ref_time := time.Now()

	for i := 0; i < NUM_PARALLEL_WORKERS; i++ {
		wg.Add(1)
		go goRoutine(&wg, extras)
	}

	wg.Wait()

	duration := time.Since(ref_time)

	println("Total time: ", duration.Seconds())
	// println("Rate: ", float64(NUM_PARALLEL_WORKERS*NUM_REQUESTS)/float64(duration.Seconds()))
	fmt.Printf("Rate: %2.2f\n", float64(NUM_PARALLEL_WORKERS*NUM_REQUESTS)/float64(duration.Seconds()))
}

func TestResourceGetAll() {
	coreTester(10000, 4, resource_adaptor.GetResources, nil)
}

func TestResourceGetSingle() {
	coreTester(10000, 4, resource_adaptor.GetResource, "df53c0bc-fdc0-475c-a875-0fef83cd15f1")
}

func main() {
	godotenv.Load("./config/config.env")
	TestResourceGetAll()
	TestResourceGetSingle()
}
