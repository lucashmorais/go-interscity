package main

import (
	"fmt"
	"os"
	"strconv"
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

func TestResourceGetAll(num_clients int, num_requests_per_client int) {
	coreTester(num_requests_per_client, num_clients, resource_adaptor.GetResources, nil)
}

func TestResourceGetAllAsync() {
	coreTester(80, 1, resource_adaptor.GetResourcesAsync, nil)
}

func TestResourceGetSingle() {
	coreTester(1, 4, resource_adaptor.GetResource, "0dd0a31b-acdf-4627-857f-fbacc4c86be4")
}

func TestResourceGetSingleAsync() {
	coreTester(800000, 1, resource_adaptor.GetResourceAsync, "0dd0a31b-acdf-4627-857f-fbacc4c86be4")
}

func TestCreateAndDeleteResource() {
	coreTester(1000, 4, resource_adaptor.CreateResourceAndDelete, nil)
}

func TestCreateAndUpdateResource() {
	coreTester(1, 4, resource_adaptor.CreateResourceAndUpdate, nil)
}

func TestCreateAndUpdateResourceAsync() {
	coreTester(100000, 1, resource_adaptor.CreateResourceAndUpdateAsync, nil)
}

func TestCreateAndUpdateAndDeleteResource() {
	coreTester(1000, 4, resource_adaptor.CreateResourceAndUpdateAndDelete, nil)
}

func TestCreateAndGetAndUpdateAndDeleteResource() {
	coreTester(1000, 4, resource_adaptor.CreateResourceAndGetAndUpdateAndDelete, nil)
}

func main() {
	godotenv.Load("./config/config.env")
	resource_adaptor.NumFails = 0

	//TODO: Move this to each of the test drivers
	// measurement.InitMeasurementArrays(1, 800000)

	args := os.Args[1:]
	num_clients, _ := strconv.Atoi(args[0])
	num_requests_per_client, _ := strconv.Atoi(args[1])

	// TestCreateAndDeleteResource()
	// TestCreateAndUpdateResource()
	// TestCreateAndUpdateAndDeleteResource()
	// TestCreateAndGetAndUpdateAndDeleteResource()
	// TestResourceGetSingle()
	// TestResourceGetSingleAsync()
	TestResourceGetAll(num_clients, num_requests_per_client)
	// TestResourceGetAllAsync()
	// TestCreateAndUpdateResource()
	// TestCreateAndUpdateResourceAsync()
}
