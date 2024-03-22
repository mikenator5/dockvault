package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	AWS   aws   `json:"aws"`
	Azure azure `json:"azure"`
}

type aws struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
}

type azure struct {
	Account   string `json:"account"`
	Container string `json:"container"`
}

func NewConfig() (Config, error) {
	filePath := filepath.Join(os.Getenv("HOME"), ".dockvault", "config")

	c := Config{}
	_, err := os.Stat(filePath)
	if err == nil {
		// File exists
		cData, err := os.ReadFile(filePath)
		if err != nil {
			return Config{}, err
		}

		if err = json.Unmarshal(cData, &c); err != nil {
			return Config{}, err
		}

	} else if os.IsNotExist(err) {
		// Create directory if it doesn't exist
		dir := filepath.Dir(filePath)
		fmt.Println(dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return Config{}, err
		}

		// Create config file
		file, err := os.Create(filePath)
		if err != nil {
			return Config{}, err
		}

		// Write basic data to config file
		cData, err := json.Marshal(c)
		if err != nil {
			return Config{}, err
		}

		file.Write(cData)
		defer file.Close()

	} else {
		// Some other error
		return Config{}, err
	}

	return c, nil
}
