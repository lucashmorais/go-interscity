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
	agent.Do(agent.Request(), nil)
}

type PostResponse struct {
	UUID string `json: UUID`
}

func CorePostResource() string {
	query_string := utils.GetServerURL() + "/resources"

	agent := fiber.Post(query_string)
	response := fiber.AcquireResponse()

	agent.Do(agent.Request(), response)

	var result PostResponse

	err := json.Unmarshal(response.Body(), &result)

	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	return result.UUID
}

func CreateResourceAndDeleteIt(wg *sync.WaitGroup, extras interface{}) {
	defer wg.Done()

	for i := 0; i < NUM_REQUESTS; i++ {
		uuid := CorePostResource()
		CoreDeleteResource(uuid)
	}
}

func GetResource(wg *sync.WaitGroup, uuid interface{}) {
	var agent *fiber.Agent
	query_string := utils.GetServerURL() + "/resources/" + uuid.(string)

	defer wg.Done()

	// This DOES NOT perform any request before the following loop starts
	agent = fiber.Get(query_string).Reuse()
	dst := make([]byte, 1000000)

	for i := 0; i < NUM_REQUESTS; i++ {
		agent.Get(dst, query_string)
	}
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
