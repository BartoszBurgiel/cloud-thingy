package client

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/BartoszBurgiel/cloud/shared"
)

func TestMessageSend(t *testing.T) {

	encryptionKey := make([]byte, 16)
	io.ReadFull(rand.Reader, encryptionKey)

	// read the package config
	conf, err := shared.NewPackageConfig(`/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/shared/testdata/packageConfig.json`)
	if err != nil {
		t.Error(err)
	}
	// start a dummy server
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			// get the message
			msg, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			submission := &shared.Submission{}
			err = json.Unmarshal(msg, submission)
			if err != nil {
				t.Error(err)
				return
			}
			pack := submission.Package
			err = pack.Decrypt(conf)
			if err != nil {
				t.Error(err)
				return
			}
		})
		fmt.Println(http.ListenAndServe(":"+shared.MiddlemanPORT, nil))
	}()
	path := "/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/"
	client, err := NewClient(
		path+"shared/testdata/",
		path+"shared/testdata/api_key",
		path+"shared/testdata/packageConfig.json",
		path+"shared/testdata/hmac_secret",
	)
	if err != nil {
		t.Error(err)
	}

	err = client.Sumbit()
	if err != nil {
		t.Error(err)
	}
}
