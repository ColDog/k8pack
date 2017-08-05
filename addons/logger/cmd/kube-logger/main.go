package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/coldog/k8pack/addons/logger/pkg/pipeline"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", "logger.json", "Config file path")
	flag.Parse()

	conf, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("failed to open file %s: %v", configFile, err)
	}

	spec := &pipeline.Spec{}

	err = json.NewDecoder(conf).Decode(spec)
	if err != nil {
		log.Fatalf("failed to read file %s: %v", configFile, err)
	}

	pipe, err := spec.GetPipeline()
	if err != nil {
		log.Fatalf("failed to init: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("starting")
	pipe.Run(ctx)
}
