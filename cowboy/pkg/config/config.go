package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

// Config : holds all the app config params
type Config struct {
	EsURL       []string `json:"esURL" required:"true"`
	RestAddress string   `json:"restAddress" required:"true"`
	GrpcAddress string   `json:"grpcAddress" required:"true"`
	Regions map[string]string `json:"regions" required:"true"`
	Kubeconfig string   `json:"kubeconfig" required:"true"`
}

var (
	config *Config
)

// Init : initialises and returns config
func Init(configFile string) {
	filePath, _ := filepath.Abs(configFile)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Unable to read config file : [ %v ].  Err : %v", configFile, err)
	}

	conf := &Config{}
	err = json.Unmarshal(file, conf)
	if err != nil {
		log.Fatalf("Unable to parse config from file : [ %v ].  Err : %v", configFile, err)
	}
	config = conf
}

// GetRestAddress : returns es url from config
func GetRestAddress() string {
	return config.RestAddress
}

// GetRegions : returns regions
func GetRegions() []string {
	var r []string
	for k,_ := range config.Regions {
		r = append(r,k);
	}
	return r
}

// GetKubeconfig : returns kubeconfig path
func GetKubeconfig(region string) string {
	return config.Regions[region]
}

// GetGrpcAddress : returns es url from config
func GetGrpcAddress() string {
	return config.GrpcAddress
}

// GetESURL : returns es url from config
func GetESURL() []string {
	return config.EsURL
}
