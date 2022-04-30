package shared

import (
	"encoding/json"
	"os"
)

// PackageConfig holds all configuration information for all packages
type PackageConfig struct {

	// MaxPackageSize defines the maximal size in bytes
	// of the compressed and encrypted package of files
	// that will be accepted by the middleman
	MaxPackageSize int

	// KeyFilePath holds the path to the file
	// which contains the key for the AES encryption
	KeyFilePath string
}

// NewPackageConfig reads the package config file and creates a new config instance
func NewPackageConfig(path string) (PackageConfig, error) {
	packageConfig, err := os.ReadFile(path)
	if err != nil {
		return PackageConfig{}, err
	}

	conf := &PackageConfig{}
	err = json.Unmarshal(packageConfig, conf)
	if err != nil {
		return PackageConfig{}, err
	}
	return *conf, nil
}
