package middleman

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/BartoszBurgiel/cloud/shared"
)

// Middleman handles the connection between the client and the server.
// Again, this struct will be used as a wrapper for a simple tls server
// that will have additional authenticating options
type Middleman struct {
	hmacSecret []byte
	apiKeyHash []byte

	// buffer for the lastly submitted package
	currentPackage shared.Package

	mutex *sync.Mutex

	// packageChecksums contains checksums of
	// all previously uploaded packages to prevent
	// multiple upload of the same package
	packageChecksums [][]byte

	port string

	log *log.Logger
}

// NewMiddlemanFromEnv returns a new Middleman instance
// where the sensitive values have been pulled from the environment variables
// of the server
func NewMiddlemanFromEnv() Middleman {
	return Middleman{
		hmacSecret: []byte(os.Getenv("HMAC_SECRET")),
		apiKeyHash: []byte(os.Getenv("API_KEY_HASH")),
		mutex:      &sync.Mutex{},
		port:       ":" + os.Getenv("PORT"),
		log:        log.New(os.Stdout, "MIDDLEMAN>", log.Ldate|log.Ltime),
	}
}

// NewMiddlemanFromFile returns a new Middleman instance
// where the sensitive values have been read from files
func NewMiddlemanFromFile(hmacSecretFilePath, apiKeyHashFilePath string) (Middleman, error) {

	hmacSecret, err := os.ReadFile(hmacSecretFilePath)
	if err != nil {
		return Middleman{}, err
	}
	apiKeyHash, err := os.ReadFile(apiKeyHashFilePath)
	if err != nil {
		return Middleman{}, err
	}
	return Middleman{
		hmacSecret: hmacSecret,
		apiKeyHash: apiKeyHash[:len(apiKeyHash)-1],
		mutex:      &sync.Mutex{},
		port:       ":7777",
		log:        log.New(os.Stdout, "MIDDLEMAN>", log.Ldate|log.Ltime),
	}, nil
}

func (m *Middleman) Start() error {

	logStart(m.log)
	http.HandleFunc(shared.MiddlemanUploadPath, m.handlePostPackage)
	http.HandleFunc(shared.MiddlemanDownloadPath, m.handleGetPackage)

	return http.ListenAndServe(m.port, nil)
}

// authenticateSubmission tells if the submission indeed was sent by the
// client or not
func (m Middleman) authenticateSubmission(s shared.Submission) bool {
	message, _ := json.Marshal(s.Package)
	macHasher := hmac.New(sha512.New, m.hmacSecret)
	macHasher.Write(message)
	MAC := macHasher.Sum(nil)
	return bytes.Equal(s.MAC, MAC)
}

// verifyAPIKey tells if the provided APIKey's hash matches the middleman's stored hash
func (m Middleman) verifyAPIKey(apiKey []byte) bool {
	hasher := sha512.New()
	hasher.Write(apiKey)
	h := make([]byte, len(m.apiKeyHash))
	hex.Encode(h, hasher.Sum(nil))
	return bytes.Equal(m.apiKeyHash, h)
}

// tryAcceptSubmission verifies and authenticates the submission.
// If all security checks are passed, the package is added to the buffer
// and the checksum is added to the set
func (m *Middleman) tryAcceptSubmission(s shared.Submission) error {
	if !m.authenticateSubmission(s) {
		logInvalidMAC(m.log, s.ID)
		return shared.TheSubmissionIsNotAuthenticated
	}

	// check if the checksum has already been uploaded
	// this can be optimised to use the BinarySearchFunc
	for _, v := range m.packageChecksums {
		if bytes.Equal(v, s.Package.CheckSum) {
			logPackageAlreadyUploaded(m.log, s.ID)
			return shared.ChecksumHasAleadyBeenUploaded
		}
	}

	if len(m.currentPackage.Payload) != 0 {
		logPackageAleadyInMemory(m.log, s.ID)
		return shared.MiddlemanHasAPackageInMemory
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.currentPackage = s.Package
	m.packageChecksums = append(m.packageChecksums, s.Package.CheckSum)
	return nil
}

// reload all environment variables
func (m *Middleman) updateEnvVariables() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.hmacSecret = []byte(os.Getenv("HMAC_SECRET"))
	m.apiKeyHash = []byte(os.Getenv("API_KEY_HASH"))
	m.port = ":" + os.Getenv("PORT")
	logVariableUpdate(m.log, *m)
}
