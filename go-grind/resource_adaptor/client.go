package resource_adaptor

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/lucashmorais/go-interscity/go-grind/utils"
)

var NUM_REQUESTS int

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
		println(response.StatusCode())
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

	for {
		agent.Do(request, response)
		if response.StatusCode() == 200 {
			break
		}
		println(response.StatusCode())
	}

	var result PostResponse

	err := json.Unmarshal(response.Body(), &result)

	if err != nil {
		// fmt.Println("[CorePostResource] error:", err)
		// fmt.Println("[CorePostResource]: Body of failed response: " + string(response.Body()))
		return ""
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

	for {
		response := fiber.AcquireResponse()
		agent.Do(request, response)
		if response.StatusCode() == 200 {
			break
		}
		println(response.StatusCode())
	}

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

func GetResources(wg *sync.WaitGroup, extras interface{}) {
	var agent *fiber.Agent
	query_string := utils.GetServerURL() + "/resources"

	defer wg.Done()

	// This DOES NOT perform any request before the following loop starts
	agent = fiber.Get(query_string).Reuse()
	dst := make([]byte, 1000000)

	for i := 0; i < NUM_REQUESTS; i++ {
		agent.Get(dst, query_string)
		// agent.Do(agent.Request(), nil)
	}
}
