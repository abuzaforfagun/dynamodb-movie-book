package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ConnectionStrings ConnectionString `json:"connection_strings"`
}

type ConnectionString struct {
	DynamoDb string `json:"dynamoDb"`
}

func (c *Config) LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&c)

	return nil
}
