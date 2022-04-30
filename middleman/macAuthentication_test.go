package middleman

import (
	"os"
	"testing"

	"github.com/BartoszBurgiel/cloud/shared"
)

func TestMacAuthentication(t *testing.T) {
	path := "/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud"

	conf, err := shared.NewPackageConfig(path + "/shared/testdata/packageConfig.json")
	if err != nil {
		t.Error(err)
	}

	p, err := shared.NewPackage(conf,
		[]string{
			path + "/shared/testdata/file1.txt",
			path + "/shared/testdata/file2.txt",
			path + "/shared/testdata/file3.txt",
		},
	)
	if err != nil {
		t.Error(err)
	}

	hmacSecret, err := os.ReadFile(path + "/shared/testdata/hmac_secret")
	if err != nil {
		t.Error(err)
	}
	sCorrect := shared.NewSubmission(hmacSecret, p)

	m := Middleman{
		hmacSecret: hmacSecret,
	}

	if !m.authenticateSubmission(sCorrect) {
		t.Errorf("Message was not correctly authenticated, but should.")
	}

	sNotCorrect := shared.NewSubmission([]byte("obviouslyOtherKey"), p)
	if m.authenticateSubmission(sNotCorrect) {
		t.Errorf("Message was correctly authenticated, but should not.")
	}
}
