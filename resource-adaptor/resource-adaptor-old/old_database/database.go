package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// This is list of external variables that are
// accessible from packages importing this one
var (
	DBConnResources     *gorm.DB
	DBConnSubscriptions *gorm.DB
)

// > curl -H "Content-Type: application/json" -X POST -d '{"data":{"lat":-23.559616,"lon":-46.731386,"description":"An equipped with various sensors in São Paulo","capabilities":["environment_monitoring"], "status":"active"}}' http://localhost:3002/components | json_pp

type DBResource struct {
	gorm.Model
	Data string `json: "data"`
	UUID string
	// Data map[string]interface{} `json: "data"`
}

type Resource struct {
	Data map[string]interface{} `json: "data"`
}

type Subscription struct {
	gorm.Model
	DBResource   DBResource
	Capabilities []string
}

type Capability struct {
	gorm.Model
	Name string `json: "name"`
}
