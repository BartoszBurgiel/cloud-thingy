package middleman

import (
	"context"
	"net/http"
	"testing"

	"github.com/BartoszBurgiel/cloud/client"
	"github.com/BartoszBurgiel/cloud/shared"
)

func TestUploadValidSubmission(t *testing.T) {
	path := "/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/"
	serv := &http.Server{Addr: ":" + shared.MiddlemanPORT}
	http.DefaultServeMux = new(http.ServeMux)

	m, err := NewMiddlemanFromFile(
		path+"shared/testdata/hmac_secret",
		path+"/middleman/testdata/apikey_hash",
	)

	if err != nil {
		t.Error(err)
	}
	http.HandleFunc(shared.MiddlemanUploadPath, m.handlePostPackage)
	go serv.ListenAndServe()

	client, err := client.NewClient(
		path+"shared/testdata/",
		path+"shared/testdata/api_key",
		path+"shared/testdata/packageConfig.json",
		path+"shared/testdata/hmac_secret",
	)

	if err := client.Sumbit(); err != nil {
		t.Error(err)
	}

	if len(m.currentPackage.Payload) == 0 {
		t.Errorf("No package has been uploaded, but should have.")
	}

	serv.Shutdown(context.TODO())
}

func TestUploadInvalidSubmission(t *testing.T) {

	path := "/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/"
	serv := &http.Server{Addr: ":" + shared.MiddlemanPORT}
	http.DefaultServeMux = new(http.ServeMux)

	m, err := NewMiddlemanFromFile(
		path+"shared/testdata/hmac_secret",
		path+"/middleman/testdata/apikey_hash",
	)

	if err != nil {
		t.Error(err)
	}
	http.HandleFunc(shared.MiddlemanUploadPath, m.handlePostPackage)
	go serv.ListenAndServe()

	client, err := client.NewClient(
		path+"shared/testdata/",
		path+"shared/testdata/api_key",
		path+"shared/testdata/packageConfig.json",
		path+"shared/testdata/invalid_hmac_secret",
	)

	if err := client.Sumbit(); err != nil {
		t.Error(err)
	}

	if len(m.currentPackage.Payload) != 0 {
		t.Errorf("Package has been uploaded, but should not have.")
	}
	serv.Shutdown(context.TODO())
}

func TestUploadRepeatedSubmission(t *testing.T) {

	path := "/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/"
	serv := &http.Server{Addr: ":" + shared.MiddlemanPORT}
	http.DefaultServeMux = new(http.ServeMux)

	m, err := NewMiddlemanFromFile(
		path+"shared/testdata/hmac_secret",
		path+"/middleman/testdata/apikey_hash",
	)

	if err != nil {
		t.Error(err)
	}
	http.HandleFunc(shared.MiddlemanUploadPath, m.handlePostPackage)
	go serv.ListenAndServe()

	client, err := client.NewClient(
		path+"shared/testdata/",
		path+"shared/testdata/api_key",
		path+"shared/testdata/packageConfig.json",
		path+"shared/testdata/hmac_secret",
	)

	for i := 0; i < 5; i++ {
		if err := client.Sumbit(); err != nil {
			t.Error(err)
		}
	}

	if len(m.packageChecksums) != 1 {
		t.Errorf("Package has been uploaded multiple times, but should only have been uploaded once.")
	}
	serv.Shutdown(context.TODO())
}
