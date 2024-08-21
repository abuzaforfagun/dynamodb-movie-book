package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	ConnectionStrings DynamoDbConnection `json:"connection_strings"`
}

type DynamoDbConnection struct {
	Region    string `json:"region"`
	TableName string `json:"table_name"`
}

func (c *Configuration) LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&c)

	return nil
}
