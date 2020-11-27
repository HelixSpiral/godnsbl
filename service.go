package godnsbl

import (
	"io/ioutil"
	"log"
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
		log.Fatal("Could not unmarshal the data:", err)
	}

	lookupService.timeout, err = time.ParseDuration(lookupService.DnsblTimeout)
	if err != nil {
		log.Fatal("Failed to parse dnsbl timeout:", lookupService.DnsblTimeout)
	}

	return lookupService
}

// Creating a new service without a config file
func NewLookupService() *LookupService {
	return &LookupService{
		DnsblTimeout: "30s",
		StartTime:    time.Now().Unix(),
	}
}

// Reading the actual config file
func readConfig(configFile string) []byte {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Could not read the config file:", err)
	}

	return yamlFile
}
