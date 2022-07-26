package config

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEnv(t *testing.T) {
	os.Setenv("AUTHENTIK_SECRET_KEY", "bar")
	cfg = nil
	Get().fromEnv()
	assert.Equal(t, "bar", Get().SecretKey)
}

func TestConfigEnv_Scheme(t *testing.T) {
	os.Setenv("foo", "bar")
	os.Setenv("AUTHENTIK_SECRET_KEY", "env://foo")
	cfg = nil
	Get().fromEnv()
	assert.Equal(t, "bar", Get().SecretKey)
}

func TestConfigEnv_File(t *testing.T) {
	file, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.Write([]byte("bar"))

	os.Setenv("AUTHENTIK_SECRET_KEY", fmt.Sprintf("file://%s", file.Name()))
	cfg = nil
	Get().fromEnv()
	assert.Equal(t, "bar", Get().SecretKey)
}
