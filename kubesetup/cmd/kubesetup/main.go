package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/coldog/k8pack/kubesetup/pkg/config"
	"github.com/coldog/k8pack/kubesetup/pkg/getter"
)

var (
	configURI string
	override  string
)

func handleErr(wrap string, err error) {
	fmt.Println(wrap + ": " + err.Error())
	os.Exit(1)
}

func main() {
	flag.StringVar(&configURI, "config-uri", "", "Config uri (required)")
	flag.StringVar(&override, "override", "", "Optional json map to override the configuration")
	flag.Parse()

	overrideData := map[string]string{}

	if override != "" {
		err := json.Unmarshal([]byte(override), &overrideData)
		if err != nil {
			handleErr("failed to parse override", err)
		}
	}

	data, err := getter.Get(configURI)
	if err != nil {
		handleErr("failed to get config", err)
	}

	config := &config.Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		handleErr("failed to unmarshal config", err)
	}
	config.Override = overrideData

	err = config.Run()
	if err != nil {
		handleErr("failed to run config", err)
	}
}
