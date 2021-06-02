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
	coreTester(2000, 4, resource_adaptor.GetResources, nil)
}

func TestResourceGetSingle() {
	coreTester(2000, 4, resource_adaptor.GetResource, "bd5329d9-1594-4b1a-881e-c6cabc5a3002")
}

func TestResourceGetSingleAsync() {
	coreTester(00000, 1, resource_adaptor.GetResourceAsync, "bd5329d9-1594-4b1a-881e-c6cabc5a3002")
}

func TestCreateAndDeleteResource() {
	coreTester(1000, 4, resource_adaptor.CreateResourceAndDelete, nil)
}

func TestCreateAndUpdateResource() {
	coreTester(10, 4, resource_adaptor.CreateResourceAndUpdate, nil)
}

func TestCreateAndUpdateAndDeleteResource() {
	coreTester(1000, 4, resource_adaptor.CreateResourceAndUpdateAndDelete, nil)
}

func TestCreateAndGetAndUpdateAndDeleteResource() {
	coreTester(1000, 4, resource_adaptor.CreateResourceAndGetAndUpdateAndDelete, nil)
}

func main() {
	godotenv.Load("./config/config.env")
	// TestResourceGetAll()
	// TestCreateAndDeleteResource()
	// TestCreateAndUpdateResource()
	// TestCreateAndUpdateAndDeleteResource()
	// TestCreateAndGetAndUpdateAndDeleteResource()
	// TestResourceGetSingle()
	TestResourceGetSingleAsync()
}
