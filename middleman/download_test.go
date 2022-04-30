package middleman

import (
	"context"
	"net/http"
	"testing"

	"github.com/BartoszBurgiel/cloud/client"
	"github.com/BartoszBurgiel/cloud/reciever"
	"github.com/BartoszBurgiel/cloud/shared"
)

func TestDownloadNoPackageSubmission(t *testing.T) {

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
	http.HandleFunc(shared.MiddlemanDownloadPath, m.handleGetPackage)

	go serv.ListenAndServe()
	rec, err := reciever.NewReciever(
		path+"shared/testdata/hmac_secret",
		path+"shared/testdata/packageConfig.json",
		path+"shared/testdata/api_key",
		path+"shared/testdata/dir/",
	)
	if err != nil {
		t.Error(err)
	}

	_, err = rec.AskForPackage()
	if err != shared.DownloadReturnsEmptyPackage {
		t.Errorf("Reciever didnt get any package, but it should have.")
	}
	serv.Shutdown(context.TODO())
}

func TestDownloadSubmission(t *testing.T) {

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
	http.HandleFunc(shared.MiddlemanDownloadPath, m.handleGetPackage)
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
	rec, err := reciever.NewReciever(
		path+"shared/testdata/hmac_secret",
		path+"shared/testdata/packageConfig.json",
		path+"shared/testdata/api_key",
		path+"shared/testdata/dir/",
	)
	if err != nil {
		t.Error(err)
	}

	ok, _ := rec.AskForPackage()
	if !ok {
		t.Errorf("Reciever didnt get any package, but it should have.")
	}
	serv.Shutdown(context.TODO())
}
