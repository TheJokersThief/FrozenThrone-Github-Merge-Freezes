package frozen_throne

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/TheJokersThief/frozen-throne/frozen_throne/config"
	"github.com/stretchr/testify/assert"
)

var wantConfig = config.Config{
	WriteSecret:   "secret",
	StorageMethod: "file",
	AuditLogKey:   "audit_log",
	GCSConfig: config.GCSConfig{
		Bucket:     "",
		BaseFolder: "",
	},
	FileConfig: config.FileConfig{
		BaseFolder: "../storage",
	},
}

func throneSetup(t *testing.T) *FrozenThrone {
	os.Setenv("WRITE_SECRET", wantConfig.WriteSecret)
	os.Setenv("STORAGE_METHOD", wantConfig.StorageMethod)
	return NewFrozenThrone(context.Background())
}

func TestNewFrozenThrone(t *testing.T) {
	throne := throneSetup(t)
	assert.Equal(t, wantConfig, throne.Config)
}

func TestFreezeAndThaw(t *testing.T) {
	repo := "testrepo"
	user := "testuser"

	throne := throneSetup(t)
	freezeErr := throne.Freeze(repo, user)
	if freezeErr != nil {
		fmt.Println(freezeErr)
		t.Fail()
	}

	_, checkErr := throne.Check(repo)
	if checkErr != nil {
		fmt.Println(checkErr)
		t.Fail()
	}

	thawErr := throne.Thaw(repo, user)
	if thawErr != nil {
		fmt.Println(thawErr)
		t.Fail()
	}
}
