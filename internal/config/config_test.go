package config

import (
	"math/rand"
	"os"
	"strconv"
	"testing"

	"mta2/commons/maps"
	"mta2/internal/constants"
	"mta2/pkg/ipconfig"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigThreshold(t *testing.T) {
	t.Run("Test with default threshold value", func(t *testing.T) {
		os.Setenv(constants.MTA_THRESHOLD, "")
		defer os.Unsetenv(constants.MTA_THRESHOLD)
		result := LoadConfigThreshold()
		assert.Equal(t, 1, result)
	})
	t.Run("Test with random threshold value", func(t *testing.T) {
		testvalue := strconv.Itoa(rand.Int())
		os.Setenv(constants.MTA_THRESHOLD, testvalue)
		defer os.Unsetenv(constants.MTA_THRESHOLD)
		testthreshold, _ := strconv.Atoi(testvalue)
		result := LoadConfigThreshold()
		assert.Equal(t, testthreshold, result)
	})

}
func TestLoadConfigIPConfiguration(t *testing.T) {
	// Create mock ipconfig.Configuration and ipconfig.IPList objects
	mockConfig := maps.NewMap()
	mockIPList := ipconfig.NewIPConfigList()

	t.Run("Test with valid JSON file path", func(t *testing.T) {
		os.Setenv(constants.PATH, "/Users/b0268986/mta2/mock/data/ipconfig.json")
		defer os.Unsetenv(constants.PATH)
		err := LoadConfigIPConfiguration(mockConfig, mockIPList)
		assert.NoError(t, err)
	})

	t.Run("Test with invalid JSON file", func(t *testing.T) {
		os.Setenv(constants.PATH, "random/random")
		defer os.Unsetenv(constants.PATH)
		err := LoadConfigIPConfiguration(mockConfig, mockIPList)
		assert.Error(t, err)
	})
	t.Run("Test with empty env variable PATH", func(t *testing.T) {
		os.Setenv(constants.PATH, "")
		defer os.Unsetenv(constants.PATH)
		err := LoadConfigIPConfiguration(mockConfig, mockIPList)
		assert.NoError(t, err)
	})
}
