package main

import (
	"context"
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/TheJokersThief/frozen-throne/frozen_throne/config"
	"github.com/TheJokersThief/frozen-throne/frozen_throne/storage"
)

type FrozenThrone struct {
	Config  config.Config
	Storage storage.StorageInterface
}

func NewFrozenThrone(ctx context.Context) *FrozenThrone {
	config := config.Config{}

	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	var throneStorage storage.StorageInterface
	switch config.StorageMethod {
	case "gcs":
		throneStorage = storage.NewGCSStorage(config, ctx)
	default:
		panic("Storage method chosen does not match the methods available")
	}

	return &FrozenThrone{
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
