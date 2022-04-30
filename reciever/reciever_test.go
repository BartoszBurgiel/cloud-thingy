package reciever

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/BartoszBurgiel/cloud/shared"
)

func TestRecievePackage(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	srv := &http.Server{
		Addr: ":" + shared.MiddlemanPORT,
	}
	path := "/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/shared/testdata"
	rec, err := NewReciever(
		path+"/hmac_secret",
		path+"/packageConfig.json",
		"/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/shared/testdata/api_key",
		"/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/shared/testdata/dir/",
	)
	if err != nil {
		t.Error(err)
	}
	// start dummy server
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			files := []string{
				path + "/file1.txt",
				path + "/file2.txt",
				path + "/file3.txt",
			}
			conf, err := shared.NewPackageConfig(path + "/packageConfig.json")
			if err != nil {
				t.Error(err)
			}
			p, err := shared.NewPackage(conf, files)
			if err != nil {
				t.Error(err)
			}
			sub := shared.NewSubmission(rec.macSecret, p)
			if err != nil {
				t.Error(err)
			}
			fmt.Fprint(w, string(sub.JSON()))
		})
		srv.ListenAndServe()
	}()

	_, err = rec.AskForPackage()
	if err != nil {
		t.Error(err)
	}
	srv.Shutdown(context.TODO())
}

func TestRecievePackageServerHasInvalidHMACSecret(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	srv := &http.Server{
		Addr: ":" + shared.MiddlemanPORT,
	}
	path := "/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/shared/testdata"
	rec, err := NewReciever(
		path+"/hmac_secret",
		path+"/packageConfig.json",
		"/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/shared/testdata/api_key",
		"/home/bartosz/dev/go/src/github.com/BartoszBurgiel/cloud/shared/testdata/dir/",
	)
	if err != nil {
		t.Error(err)
	}
	// start dummy server
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			files := []string{
				path + "/file1.txt",
				path + "/file2.txt",
				path + "/file3.txt",
			}
			conf, err := shared.NewPackageConfig(path + "/packageConfig.json")
			if err != nil {
				t.Error(err)
			}
			p, err := shared.NewPackage(conf, files)
			if err != nil {
				t.Error(err)
			}
			sub := shared.NewSubmission([]byte("InvalidHMACSecret"), p)
			if err != nil {
				t.Error(err)
			}
			fmt.Fprint(w, string(sub.JSON()))
		})
		srv.ListenAndServe()
	}()

	_, err = rec.AskForPackage()
	if err != shared.MiddlemanPackageHasInvalidMAC {
		t.Errorf("Asking should fail but didn't")
	}
	srv.Shutdown(context.TODO())
}
