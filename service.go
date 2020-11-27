package godnsbl

import (
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Creating a new service with a config file
func NewLookupServiceWithConfig(configFile string) *LookupService {
	lookupService := NewLookupService()

	fileInfo, err := os.Stat(configFile)

	// If the config file doesn't exist or it's a directory just return an empty service
	if os.IsNotExist(err) || fileInfo.IsDir() {
		return lookupService
	}

	configData := readConfig(configFile)

	err = yaml.Unmarshal(configData, lookupService)
	if err != nil {
		panic(err)
	}

	return lookupService
}

// Creating a new service without a config file
func NewLookupService() *LookupService {
	return &LookupService{
		StartTime: time.Now().Unix(),
	}
}

// Reading the actual config file
func readConfig(configFile string) []byte {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	return yamlFile
}
