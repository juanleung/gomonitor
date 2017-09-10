package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

// Configuration is a key value store to load a json config file
type Configuration struct {
	store map[string]interface{}
}

// GetValue return the value store for a given key
func (c *Configuration) GetValue(key string) (string, error) {
	_, ok := c.store[key]

	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}

	t := reflect.TypeOf(c.store[key]).String()
	if t != "string" {
		return "", fmt.Errorf("the value is not a string")
	}

	return c.store[key].(string), nil
}

// LoadConfigurationJSON load a json into a Configuration struct
func LoadConfigurationJSON(filename string) (Configuration, error) {
	file, err := os.Open(filename)
	config := Configuration{}

	if err != nil {
		return config, fmt.Errorf("an error ocurred reading the json file: %v", err)
	}
	var cs interface{}
	err = json.NewDecoder(file).Decode(&cs)
	if err != nil {
		log.Panic(
			fmt.Errorf(
				"Ocurrió un error decodificando el json de la configuración: %v",
				err))
	}
	config.store = cs.(map[string]interface{})
	return config, nil
}
