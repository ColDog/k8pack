package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/coldog/k8pack/kubesetup/pkg/config"
	"github.com/coldog/k8pack/kubesetup/pkg/getter"
)

var configURI string

func handleErr(wrap string, err error) {
	println(wrap + ": " + err.Error())
	os.Exit(1)
}

func main() {
	flag.StringVar(&configURI, "config-uri", "", "Config uri, also can be specified at /etc/kubernetes/.config-uri")
	flag.Parse()

	if configURI == "" {
		data, err := ioutil.ReadFile("/etc/kubernetes/.config-uri")
		if err != nil {
			handleErr("failed to find config", err)
		}
		configURI = strings.TrimSpace(string(data))
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

	err = config.Run()
	if err != nil {
		handleErr("failed to run config", err)
	}
}
