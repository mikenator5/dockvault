package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	homeDirName = ".dockvault"
	configName  = "config.json"
)

type Config struct {
	AWS   *AWS   `json:"aws"`
	Azure *Azure `json:"azure"`
}

type AWS struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
}

type Azure struct {
	StorageAccount string `json:"storageAccount"`
	Container      string `json:"container"`
}

func writeConfig(cfg Config, filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	fmt.Println(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create config file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// Write basic data to config file
	cData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if _, err = file.Write(cData); err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func updateConfig(cfg Config) error {
	filePath := filepath.Join(os.Getenv("HOME"), homeDirName, configName)

	_, err := os.Stat(filePath)
	if err == nil {
		// File exists => update it
		cData, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		c := Config{}
		if err = json.Unmarshal(cData, &c); err != nil {
			return err
		}

		if cfg.AWS != nil {
			c.AWS = cfg.AWS
		}
		if cfg.Azure != nil {
			c.Azure = cfg.Azure
		}

		return writeConfig(c, filePath)

	} else if os.IsNotExist(err) {
		// File does not exist => create it
		return writeConfig(cfg, filePath)
	}

	return nil
}

func NewAWSConfig(aws *AWS) error {
	cfg := Config{
		AWS: aws,
	}
	return updateConfig(cfg)
}

func NewAzureConfig(az *Azure) error {
	cfg := Config{
		Azure: az,
	}
	return updateConfig(cfg)
}

func GetConfig() (Config, error) {
	filePath := filepath.Join(os.Getenv("HOME"), homeDirName, configName)

	_, err := os.Stat(filePath)
	if err == nil {
		// File exists
		cData, err := os.ReadFile(filePath)
		if err != nil {
			return Config{}, err
		}

		c := Config{}
		if err = json.Unmarshal(cData, &c); err != nil {
			return Config{}, err
		}
		return c, nil
	}

	return Config{}, errors.New("Config file was not found. Create one with: dockvault configure")
}
