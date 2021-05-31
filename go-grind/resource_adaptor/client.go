package resource_adaptor

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/lucashmorais/go-interscity/go-grind/utils"
)

var NUM_REQUESTS int

func BasicGetResources() {
	println(utils.GetServerURL())
	agent := fiber.Get(utils.GetServerURL() + "/components")

	_, response, _ := agent.String()
	println(response)
}

func GetResource(wg *sync.WaitGroup, uuid interface{}) {
	var agent *fiber.Agent
	query_string := utils.GetServerURL() + "/components/" + uuid.(string)

	defer wg.Done()

	agent = fiber.Get(query_string).Reuse()
	dst := make([]byte, 1000000)

	for i := 0; i < NUM_REQUESTS; i++ {
		agent.Get(dst, query_string)
	}
}

func GetResources(wg *sync.WaitGroup, extras interface{}) {
	var agent *fiber.Agent
	query_string := utils.GetServerURL() + "/components"

	defer wg.Done()

	agent = fiber.Get(query_string).Reuse()
	dst := make([]byte, 1000000)

	for i := 0; i < NUM_REQUESTS; i++ {
		agent.Get(dst, query_string)
	}
}
