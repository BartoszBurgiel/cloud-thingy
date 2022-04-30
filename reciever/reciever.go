package reciever

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BartoszBurgiel/cloud/shared"
	"github.com/google/uuid"
)

// Reciever is a wrapper-struct for the request to the middleman
// to recieve the payload
type Reciever struct {
	macSecret             []byte
	packageConfigFilePath string
	apiKey                string
	destinationRootPath   string
	statsFile             string

	middlemanURL string
	log          *log.Logger
}

func NewRecieverFromConfig(confFilePath, destinationRootPath string) (Reciever, error) {
	conf, err := newConf(confFilePath)
	if err != nil {
		return Reciever{}, err
	}
	l := log.New(os.Stdout, fmt.Sprintf("RECIEVER(id: %s)> ", uuid.New().String()), log.Ldate|log.Ltime)
	logInit(l, destinationRootPath)
	macSecret, err := os.ReadFile(conf.MacSecretFile)
	logError(l, err)
	if err != nil {
		return Reciever{}, err
	}

	apiKey, err := os.ReadFile(conf.ApiKeyFile)
	logError(l, err)
	if err != nil {
		return Reciever{}, err
	}

	encryptionKey := make([]byte, 16)
	io.ReadFull(rand.Reader, encryptionKey)

	return Reciever{
		macSecret:             macSecret[:len(macSecret)-1],
		apiKey:                string(apiKey[:len(apiKey)-1]),
		packageConfigFilePath: conf.PackageConfigFile,
		destinationRootPath:   destinationRootPath,
		middlemanURL:          conf.MiddlemanDownloadURL,
		statsFile:             conf.StatsFile,
		log:                   l,
	}, nil
}

// NewReciever creates a new Reciever instance
func NewReciever(macSecretFilePath, packageConfigFilePath, apiKeyFilePath, destinationRootPath string) (Reciever, error) {
	l := log.New(os.Stdout, fmt.Sprintf("RECIEVER(id: %s)> ", uuid.New().String()), log.Ldate|log.Ltime)
	macSecret, err := os.ReadFile(macSecretFilePath)
	logError(l, err)
	if err != nil {
		return Reciever{}, err
	}

	apiKey, err := os.ReadFile(apiKeyFilePath)
	logError(l, err)
	if err != nil {
		return Reciever{}, err
	}

	encryptionKey := make([]byte, 16)
	io.ReadFull(rand.Reader, encryptionKey)

	return Reciever{
		macSecret:             macSecret[:len(macSecret)-1],
		apiKey:                string(apiKey[:len(apiKey)-1]),
		packageConfigFilePath: packageConfigFilePath,
		destinationRootPath:   destinationRootPath,
		log:                   l,
	}, nil
}

// AskForPackage sends a request to the middleman
// for a package. If there is a package aviable, the middleman
// sends a package back, and the reciever decrypts it and writes it to the destination folder
// elsewise, the method terminates.
func (r Reciever) AskForPackage() (bool, error) {

	req, err := http.NewRequest(
		http.MethodGet,
		r.middlemanURL,
		nil,
	)
	logError(r.log, err)
	if err != nil {
		return false, err
	}
	q := req.URL.Query()
	q.Add("api_key", r.apiKey)
	req.URL.RawQuery = q.Encode()

	c := &http.Client{}
	t := time.Now()
	res, err := c.Do(req)
	logError(r.log, err)
	if err != nil {
		return false, err
	}

	logMiddlemanResponse(r.log, res.Status, time.Since(t))
	if res.StatusCode != 200 {
		return false, fmt.Errorf("Invalid response from the middleman to the reciever when asking for package. \nStatusCode: %d\n", res.StatusCode)
	}

	finish := make(chan bool)
	go func() {
		fmt.Print("Downloading...")
		for {
			select {
			case _, ok := <-finish:
				if ok {
					fmt.Println("\nFinished!")
					return
				}
			default:
				fmt.Print(".")
				time.Sleep(time.Second)
			}
		}
	}()
	t = time.Now()
	body, err := ioutil.ReadAll(res.Body)
	finish <- true
	if err != nil {
		return false, err
	}

	sf, err := os.OpenFile(r.statsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer sf.Close()
	if err != nil {
		return true, err
	}
	_, err = sf.WriteString(
		stats{
			size: len(body),
			dur:  time.Since(t),
			time: time.Now(),
		}.toCSV(),
	)
	if err != nil {
		return true, err
	}

	sub := shared.Submission{}
	err = json.Unmarshal(body, &sub)
	logError(r.log, err)
	if err != nil {
		return false, err
	}

	if len(sub.Package.Payload) == 0 {
		return false, shared.DownloadReturnsEmptyPackage
	}

	authenticated := r.verifyMAC(sub)
	logAuthentication(r.log, authenticated, sub.ID)
	if !r.verifyMAC(sub) {
		return false, shared.MiddlemanPackageHasInvalidMAC
	}

	pack := sub.Package
	conf, err := shared.NewPackageConfig(r.packageConfigFilePath)
	logError(r.log, err)
	if err != nil {
		return false, err
	}

	err = pack.Decrypt(conf)
	logError(r.log, err)
	logDecryption(r.log, err != shared.FailedDecryptionOfThePackage, sub.ID)
	if err != nil {
		return false, err
	}

	pathToZip := r.destinationRootPath + string(os.PathSeparator) + time.Now().Format("2006-01-02-15-04-05") + ".zip"
	f, err := os.Create(pathToZip)
	defer f.Close()
	logError(r.log, err)
	if err != nil {
		return false, err
	}

	t = time.Now()
	_, err = io.Copy(f, bytes.NewReader(pack.Payload))
	logError(r.log, err)
	if err != nil {
		return false, err
	}

	logSuccessfullDownload(r.log, sub.ID)
	return true, nil
}

// verifyMAC tells if the MAC was sealed with the correct
// secret
func (r Reciever) verifyMAC(s shared.Submission) bool {

	message, _ := json.Marshal(s.Package)
	// generate the MAC
	macHasher := hmac.New(sha512.New, r.macSecret)
	macHasher.Write(message)
	MAC := macHasher.Sum(nil)
	return bytes.Equal(MAC, s.MAC)
}
