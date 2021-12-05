package storage

import "time"

type StorageInterface interface {
	PlaceKey(name string) error
	RemoveKey(name string) error
	GetKey(name string) (string, error)
	LogAction(action string, user string, repo string, metadata ThroneMetadata) error
}

type ThroneMetadata struct {
	Action    string    `json:"action"`
	User      string    `json:"user"`
	Repo      string    `json:"repo"`
	Timestamp time.Time `json:"timestamp"`
}
