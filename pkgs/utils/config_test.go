package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFile(t *testing.T) {
	file := getConfigFile()
	assert.NotEqual(t, len(file), 0)
}

func TestGetConfig(t *testing.T) {
	assert.NotEqual(t, 0, len(GetConfig("mysql")["user"].(string)))
}
