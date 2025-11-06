package main

import (
	"encoding/json"
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/rakibulbanna/go-fiber-postgres/models"
)

func main() {
	// Load all models
	modelsList := []interface{}{
		&models.User{},
		&models.Book{},
	}

	// Extract schema using Atlas GORM provider
	provider := gormschema.New("postgres")
	schema, err := provider.Load(modelsList...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load schema: %v\n", err)
		os.Exit(1)
	}

	// Output the schema as JSON (Atlas format)
	result := map[string]string{
		"url": schema,
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode result: %v\n", err)
		os.Exit(1)
	}
}
