package storage

type StorageInterface interface {
	PlaceKey(name string) (bool, error)
	RemoveKey(name string) (bool, error)
	GetKey(name string) (bool, error)
	LogAction(action string, user string, repo string, metadata ThroneMetadata) (bool, error)
}

type ThroneMetadata struct {
	Action    string
	User      string
	Repo      string
	Timestamp string
}
