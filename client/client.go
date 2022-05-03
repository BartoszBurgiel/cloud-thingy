package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/BartoszBurgiel/cloud/shared"
	"github.com/google/uuid"
)

// Client is a wrapper-struct for request to the middleman
// to upload the payload
type Client struct {
	macSecret []byte

	// files for the package
	files []string

	apiKey []byte

	packageConfigFilePath string

	uploadURL string

	statsFile string

	log *log.Logger
}

func NewClientFromConfigFile(configFilePath, payloadRootPath string) (Client, error) {
	l := log.New(os.Stdout, "CLIENT>", log.Ltime)
	logInitClient(l, payloadRootPath)
	conf, err := newConf(configFilePath)
	if err != nil {
		logError(l, err)
		return Client{}, err
	}
	cl := Client{
		log:                   l,
		packageConfigFilePath: conf.PackageConfigFilePath,
		uploadURL:             conf.MiddlemanUploadURL,
		statsFile:             conf.StatsFile,
	}

	macSecret, err := os.ReadFile(conf.MACSecretFilePath)
	logError(cl.log, err)
	if err != nil {
		return Client{}, err
	}
	cl.macSecret = macSecret[:len(macSecret)-1]

	apiKey, err := os.ReadFile(conf.ApiKeyFilePath)
	logError(cl.log, err)
	if err != nil {
		return Client{}, err
	}
	cl.apiKey = apiKey[:len(apiKey)-1]
	filepath.Walk(payloadRootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			logAddingFileToArchive(cl.log, path)
			cl.files = append(cl.files, path)
		}
		return nil
	})
	return cl, nil
}

// NewClient reads all filenames recursively, i.e., if a directory was submitted
// and all constant information
func NewClient(payloadRootPath, apiKeyFilePath, packageConfigFilePath, macSecretFilePath string) (Client, error) {

	l := log.New(os.Stdout, fmt.Sprintf("CLIENT(id: %s)> ", uuid.New().String()), log.Ldate|log.Ltime)
	logInitClient(l, payloadRootPath)
	cl := Client{
		packageConfigFilePath: packageConfigFilePath,
		log:                   l,
	}
	macSecret, err := os.ReadFile(macSecretFilePath)
	logError(cl.log, err)
	if err != nil {
		return Client{}, err
	}
	cl.macSecret = macSecret

	apiKey, err := os.ReadFile(apiKeyFilePath)
	logError(cl.log, err)
	if err != nil {
		return Client{}, err
	}
	cl.apiKey = apiKey[:len(apiKey)-1]
	filepath.Walk(payloadRootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			logAddingFileToArchive(cl.log, path)
			cl.files = append(cl.files, path)
		}
		return nil
	})

	return cl, nil
}

func (c Client) Sumbit() error {

	// read the package config
	packageConfig, err := os.ReadFile(c.packageConfigFilePath)
	logError(c.log, err)
	if err != nil {
		return err
	}

	conf := &shared.PackageConfig{}
	err = json.Unmarshal(packageConfig, conf)
	logError(c.log, err)
	if err != nil {
		return err
	}

	t := time.Now()
	pack, err := shared.NewPackage(*conf, c.files)
	logError(c.log, err)
	if err != nil {
		return err
	}

	s := shared.NewSubmission(c.macSecret, pack)
	logSubmissionPreparationResult(c.log, s)
	requestData, err := json.Marshal(s)
	logError(c.log, err)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.uploadURL,
		bytes.NewReader(requestData),
	)
	logError(c.log, err)
	q := req.URL.Query()
	q.Add("api_key", string(c.apiKey))
	req.URL.RawQuery = q.Encode()

	httpC := &http.Client{}
	t = time.Now()
	finish := make(chan bool)
	go shared.DisplayProgressBar("Uploading...", finish)
	res, err := httpC.Do(req)
	finish <- true
	logError(c.log, err)
	if err != nil {
		return err
	}

	dur := time.Since(t)
	msg, err := ioutil.ReadAll(res.Body)
	logError(c.log, err)
	if err != nil {
		return err
	}
	logRequestToTheMiddleman(c.log, string(c.apiKey), dur, string(msg))
	f, err := os.OpenFile(c.statsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(
		stats{
			size: len(requestData),
			dur:  dur,
			time: time.Now(),
		}.toCSV(),
	)
	if err != nil {
		return err
	}
	return nil
}
