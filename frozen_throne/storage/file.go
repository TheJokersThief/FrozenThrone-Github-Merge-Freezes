package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/TheJokersThief/frozen-throne/frozen_throne/config"
)

type FileStorage struct {
	ServiceConfig *config.Config
	BaseFolder    string
}

// NewFileStorage returns a GCS storage engine
func NewFileStorage(config config.Config) *FileStorage {

	baseFolder := config.FileConfig.BaseFolder
	_, storageFolderErr := os.Stat(baseFolder)
	if os.IsNotExist(storageFolderErr) {
		os.Mkdir(baseFolder, 0744)
	}

	auditFolder := fmt.Sprintf("%s/%s", baseFolder, config.AuditLogKey)
	_, auditFolderErr := os.Stat(auditFolder)
	if os.IsNotExist(auditFolderErr) {
		os.Mkdir(auditFolder, 0744)
	}

	return &FileStorage{
		ServiceConfig: &config,
		BaseFolder:    config.FileConfig.BaseFolder,
	}
}

// PlaceKey writes an object with the name supplied to a GCS bucket
func (s *FileStorage) PlaceKey(name string) error {
	return s.writeObject(name, "")
}

// RemoveKey deletes an object with the supplied name from a GCS bucket
func (s *FileStorage) RemoveKey(name string) error {
	objName := fmt.Sprintf("%v/%v", s.BaseFolder, name)
	err := os.Remove(objName)
	if err != nil {
		return err
	}

	return nil
}

// GetKey retrieves the bucket object and returns its contents
func (s *FileStorage) GetKey(name string) (string, error) {
	objName := fmt.Sprintf("%v/%v", s.BaseFolder, name)
	_, err := os.Stat(objName)
	if os.IsNotExist(err) {
		return "", err
	}

	content, err := ioutil.ReadFile(objName)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// LogAction creates an object with audit data. The intention is to use lifecycle rules to clean up old, unneeded logs.
func (s *FileStorage) LogAction(action string, user string, repo string, metadata ThroneMetadata) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	timestampPrefix := metadata.Timestamp.Format("2006-01-02_15:04:05")
	objName := fmt.Sprintf("%s/%s_%s.json", s.ServiceConfig.AuditLogKey, timestampPrefix, repo)
	return s.writeObject(objName, string(data))
}

// writeObject is a helper function to write an object to a GCS bucket in the config's base folder
func (s *FileStorage) writeObject(name string, content string) error {
	objName := fmt.Sprintf("%s/%s", s.BaseFolder, name)
	fmt.Println(objName)
	err := ioutil.WriteFile(objName, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
