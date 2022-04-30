package shared

import (
	"testing"
)

func TestEncryption(t *testing.T) {
	conf := PackageConfig{
		KeyFilePath:    "./testdata/encryption_key",
		MaxPackageSize: 2 << 16,
	}

	files := []string{
		"./testdata/file1.txt",
		"./testdata/file2.txt",
		"./testdata/file3.txt",
	}

	p, err := NewPackage(conf, files)
	if err != nil {
		t.Error(err)
	}

	// check if the package is indeed encrypted
	if p.isPayloadValid() {
		t.Errorf("Package should be encrypted, but isn't")
	}

	err = p.Decrypt(conf)
	if err != nil {
		t.Error(err)
	}

	if !p.isPayloadValid() {
		t.Errorf("The decrytion was not succesfull")
	}
}
