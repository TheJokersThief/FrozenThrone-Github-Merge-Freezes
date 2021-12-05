package main

import (
	"context"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"

	"frozen_throne/frozen_throne/config"
	"frozen_throne/frozen_throne/storage"
)

type FrozenThrone struct {
	Name    string
	Config  interface{}
	Storage storage.StorageInterface
}

func NewFrozenThrone(name string, ctx context.Context) *FrozenThrone {
	config := config.Config{}

	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	var throneStorage storage.StorageInterface
	switch config.StorageMethod {
	case "gcs":
		throneStorage = storage.NewGCSStorage(config, ctx)
	case "redis":
		throneStorage = storage.NewRedisStorage()
	default:
		panic("Storage method chosen does not match the methods available")
	}

	return &FrozenThrone{
		Name:    name,
		Config:  config,
		Storage: throneStorage,
	}
}

func (f *FrozenThrone) Freeze(name string) error {
	return f.Storage.PlaceKey(name)
}

func (f *FrozenThrone) Unfreeze(name string) error {
	return f.Storage.RemoveKey(name)
}

func (f *FrozenThrone) Check(name string) (string, error) {
	return f.Storage.GetKey(name)
}

func IngestHttp(w http.ResponseWriter, r *http.Request) {
	// throne := NewFrozenThrone("test")
}
