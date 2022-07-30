package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hanchchch/storenv/pkg/configuration"
	"github.com/hanchchch/storenv/pkg/store"
)

type Action string

const (
	STORE   Action = "store"
	LOAD    Action = "load"
	DESTROY Action = "destroy"
)

var configfile string
var action Action

func usage() {
	fmt.Println("")
	fmt.Println("actions")
	fmt.Println("  store")
	fmt.Println("  load")
	fmt.Println("")
	fmt.Println("options")
	flag.PrintDefaults()
}

func parseArgs() {
	flag.StringVar(&configfile, "c", "storenv.yaml", "path of the configuration file")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	action = Action(args[0])
}

func main() {
	parseArgs()

	config, err := configuration.NewConfiguration(configfile)
	if err != nil {
		panic(fmt.Errorf("failed to parse configuration: %v", err))
	}

	var storage store.Storage
	if config.Storage.S3.Prefix != "" {
		storage = store.NewS3Storage(store.S3StorageOptions{
			AccessKeyId:     "",
			SecretAccessKey: "",
			Region:          config.Storage.S3.Region,
			Bucket:          config.Storage.S3.Bucket, // TODO get default from ~/.storenv
			Prefix:          config.Storage.S3.Prefix,
		})
		if err != nil {
			panic(fmt.Errorf("failed to initialize s3 storage: %v", err))
		}
	}

	if storage == nil {
		panic(fmt.Errorf("no specified storage"))
	}

	if action == STORE {
		fmt.Printf("store %d secrets into %v storage\n", len(config.Secrets), storage.Name())
		err := storage.Store(config.Secrets)
		if err != nil {
			panic(err)
		}
	}

	if action == LOAD {
		fmt.Printf("load %d secrets from %v storage\n", len(config.Secrets), storage.Name())
		err := storage.Load(config.Secrets)
		if err != nil {
			panic(err)
		}
	}
}
