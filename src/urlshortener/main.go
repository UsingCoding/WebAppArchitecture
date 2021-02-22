package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port        int               `json:"port"`
	DefaultPath string            `json:"default_path"`
	Paths       map[string]string `json:"paths"`
}

func GetConfigPath() string {
	configPath := flag.String("f", "config.json", "A config file path")
	flag.Parse()
	return *configPath
}

func ReadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open file %s error: %s\n", configPath, err))
	}

	stringConfig, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read from config file %s error: %s\n", configPath, err))
	}

	err = file.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("close config file %s error: %s", configPath, err))
	}

	config := Config{}

	err = json.Unmarshal(stringConfig, &config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse config file %s error: %s\n", configPath, err))
	}

	return &config, nil
}

type RedirectProxy struct {
	notFoundRedirectUrl string
	handlers            map[string]string
}

func (proxy RedirectProxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	redirectUrl, ok := proxy.handlers[request.URL.Path]
	if !ok {
		redirectUrl = proxy.notFoundRedirectUrl
	}

	log.Printf("Moved from %s to %s\n", request.URL.Path, redirectUrl)

	http.Redirect(writer, request, redirectUrl, http.StatusTemporaryRedirect)
}

func buildRedirectProxy(config *Config) *RedirectProxy {
	return &RedirectProxy{
		notFoundRedirectUrl: config.DefaultPath,
		handlers:            config.Paths,
	}
}

func validateConfig(config *Config) error {
	if config.Port <= 0 {
		return errors.New("invalid port")
	}

	if config.DefaultPath == "" {
		return errors.New("default path must be not empty")
	}

	return nil
}

func main() {
	configPath := GetConfigPath()

	config, err := ReadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
		return
	}

	err = validateConfig(config)

	if err != nil {
		log.Fatalf("Config invalid: %s", err)
	}

	proxy := buildRedirectProxy(config)

	log.Printf("Starting server on port: %d", config.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), proxy)
	if err != nil {
		log.Fatalf("Failed to start http server: %s", err)
	}
}
