package resource_adaptor

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lucashmorais/go-interscity/go-grind/utils"
	"golang.org/x/sync/semaphore"
)

var NUM_REQUESTS int
var NumFails int
var count int = 0

func BasicGetResources() {
	println(utils.GetServerURL())
	agent := fiber.Get(utils.GetServerURL() + "/resources")

	_, response, _ := agent.String()
	println(response)
}

func CoreDeleteResource(uuid string) {
	query_string := utils.GetServerURL() + "/resources/" + uuid

	agent := fiber.Delete(query_string)
	for {
		response := fiber.AcquireResponse()
		agent.Do(agent.Request(), response)
		if response.StatusCode() == 200 {
			break
		}
		println(response.StatusCode())
	}
}

type PostResponse struct {
	UUID string `json: UUID`
}

func CorePutResource(uuid string, marshalled_body []byte) {
	query_string := utils.GetServerURL() + "/resources/" + uuid

	agent := fiber.Put(query_string)
	request := agent.Request()

	request.SetBody(marshalled_body)
	request.Header.Add("Content-type", "application/json")

	for {
		response := fiber.AcquireResponse()
		agent.Do(request, response)
		if response.StatusCode() == 200 {
			break
		}
		println("Failing at endpoint: " + query_string)
		time.Sleep(time.Millisecond * 100)
		// println(response.StatusCode())
	}
}

func CreateResourceAndGetAndUpdateAndDelete(wg *sync.WaitGroup, extras interface{}) {
	defer wg.Done()

	examplePostBody := GetExampleResourceBody()

	for i := 0; i < NUM_REQUESTS; i++ {
		var uuid string
		for {
			uuid = CorePostResource(examplePostBody)
			if uuid != "" {
				break
			}
		}
		CoreGetResource(uuid, nil)
		CorePutResource(uuid, examplePostBody)
		CoreDeleteResource(uuid)
	}
}

func CreateResourceAndUpdateAndDelete(wg *sync.WaitGroup, extras interface{}) {
	defer wg.Done()

	examplePostBody := GetExampleResourceBody()

	for i := 0; i < NUM_REQUESTS; i++ {
		var uuid string
		for {
			uuid = CorePostResource(examplePostBody)
			if uuid != "" {
				break
			}
		}
		CorePutResource(uuid, examplePostBody)
		CoreDeleteResource(uuid)
	}
}

func CreateResourceAndUpdate(wg *sync.WaitGroup, extras interface{}) {
	defer wg.Done()

	examplePostBody := GetExampleResourceBody()

	for i := 0; i < NUM_REQUESTS; i++ {
		uuid := CorePostResource(examplePostBody)
		CorePutResource(uuid, examplePostBody)
	}
}

func CoreCreateResourceAndUpdate(wgInner *sync.WaitGroup, examplePostBody []byte) {
	defer wgInner.Done()

	uuid := CorePostResource(examplePostBody)
	CorePutResource(uuid, examplePostBody)
}

func CreateResourceAndUpdateAsync(wg *sync.WaitGroup, extras interface{}) {
	wgInner := sync.WaitGroup{}
	defer wg.Done()

	examplePostBody := GetExampleResourceBody()

	for i := 0; i < NUM_REQUESTS; i++ {
		wgInner.Add(1)
		go CoreCreateResourceAndUpdate(&wgInner, examplePostBody)
	}

	println("Waiting for requests to be completed.")
	wgInner.Wait()
}

// Fields are not encodable by json.Marshal if they
// are not exported by having their names starting
// with a capital letter.
type Resource struct {
	Description  string   `json:"description"`
	Capabilities []string `json:"capabilities"`
	Status       string   `json:"status"`
	Lat          float64  `json:"lat"`
	Lon          float64  `json:"lon"`
}

func GetExampleResourceBody() []byte {
	body, err := json.Marshal(map[string]interface{}{"data": Resource{
		Description: "A public bus",
		Capabilities: []string{
			"temperature",
			"humidity",
			"illuminate",
		},
		Status: "active",
		Lat:    -23.559616,
		Lon:    -46.731386,
	}})

	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	println(string(body))
	return body
}

func CorePostResource(body []byte) string {
	query_string := utils.GetServerURL() + "/resources"

	agent := fiber.Post(query_string)
	request := agent.Request()
	response := fiber.AcquireResponse()

	request.SetBody(body)
	request.Header.Add("Content-type", "application/json")

	numLocalFailures := 0

	var result PostResponse
	for {
		for {
			agent.Do(request, response)
			if response.StatusCode() == 200 {
				break
			}
			println(response.StatusCode())
		}

		err := json.Unmarshal(response.Body(), &result)

		if err == nil {
			break
		}

		NumFails++
		numLocalFailures++
		if NumFails%1000 == 0 {
			println("[CorePostResource]: Number of POST failures so far: ", NumFails)
		}
		time.Sleep(time.Millisecond * time.Duration(numLocalFailures*rand.Intn(200)))
	}

	return result.UUID
}

func CreateResourceAndDelete(wg *sync.WaitGroup, extras interface{}) {
	defer wg.Done()

	examplePostBody := GetExampleResourceBody()

	for i := 0; i < NUM_REQUESTS; i++ {
		uuid := CorePostResource(examplePostBody)
		CoreDeleteResource(uuid)
	}
}

func CoreGetResource(uuid string, wgInner *sync.WaitGroup) {
	query_string := utils.GetServerURL() + "/resources/" + uuid
	if wgInner != nil {
		defer wgInner.Done()
	}

	agent := fiber.Get(query_string)
	request := agent.Request()

	request.Header.Add("Content-type", "application/json")

	// start := measurement.StartMeasuringRequest()
	for {
		response := fiber.AcquireResponse()
		agent.Do(request, response)
		if response.StatusCode() == 200 {
			break
		}
		println(response.StatusCode())
	}
	// measurement.StopMeasuringRequest(start, 0)

	agent.CloseIdleConnections()
}

func GetResource(wg *sync.WaitGroup, uuid interface{}) {
	defer wg.Done()

	// This DOES NOT perform any request before the following loop starts
	for i := 0; i < NUM_REQUESTS; i++ {
		CoreGetResource(uuid.(string), nil)
	}
}

func GetResourceAsync(wg *sync.WaitGroup, uuid interface{}) {
	defer wg.Done()

	wgInner := sync.WaitGroup{}

	for i := 0; i < NUM_REQUESTS; i++ {
		wgInner.Add(1)
		go CoreGetResource(uuid.(string), &wgInner)

		// if i%200 == 0 {
		// 	wgInner.Wait()
		// }
	}

	println("Waiting for requests to be completed.")
	wgInner.Wait()
}

func CoreGetResources(wgInner *sync.WaitGroup, ctx *context.Context, concurrencyLimiter *semaphore.Weighted) {
	query_string := utils.GetServerURL() + "/resources"
	if wgInner != nil {
		defer wgInner.Done()
	}

	defer concurrencyLimiter.Release(1)

	agent := fiber.Get(query_string)
	request := agent.Request()

	request.Header.Add("Content-type", "application/json")

	backoffFactor := 12
	for {
		response := fiber.AcquireResponse()
		agent.Do(request, response)
		if response.StatusCode() == 200 && len(response.Body()) != 0 {
			break
		}

		backoffFactor++
		println(backoffFactor)
		// time.Sleep(time.Duration(rand.Int31n(int32(math.Pow(2, float64(backoffFactor))))) * time.Millisecond)
	}

	count++

	if count%1000 == 0 {
		println(count)
	}

	// agent.CloseIdleConnections()
	// println(string(response.Body()))
}

func GetResources(wg *sync.WaitGroup, extras interface{}) {
	var agent *fiber.Agent
	query_string := utils.GetServerURL() + "/resources"

	defer wg.Done()

	// This DOES NOT perform any request before the following loop starts
	agent = fiber.Get(query_string).Reuse()
	request := agent.Request()
	request.Header.Add("Content-type", "application/json")
	// dst := make([]byte, 1000000)

	// println("testing1")
	for i := 0; i < NUM_REQUESTS; i++ {
		// agent.Get(dst, query_string)

		// println("testing2")
		backoffFactor := 12
		for {
			// println("testing3")
			response := fiber.AcquireResponse()
			// println("testing4")
			agent.Do(request, response)
			// println("testing5")
			if response.StatusCode() == 200 && len(response.Body()) != 0 {
				break
			}

			backoffFactor++
			println(backoffFactor)
			// time.Sleep(time.Duration(rand.Int31n(int32(math.Pow(2, float64(backoffFactor))))) * time.Millisecond)
		}
	}
}

func GetResourcesAsync(wg *sync.WaitGroup, extras interface{}) {
	wgInner := sync.WaitGroup{}

	defer wg.Done()

	var concurrencyLimiter = semaphore.NewWeighted(int64(1))
	// ctx := context.Background()
	ctx := context.TODO()

	for i := 0; i < NUM_REQUESTS; i++ {
		wgInner.Add(1)
		concurrencyLimiter.Acquire(ctx, 1)
		go CoreGetResources(&wgInner, &ctx, concurrencyLimiter)
	}

	println("Waiting for requests to be completed.")
	wgInner.Wait()
}
