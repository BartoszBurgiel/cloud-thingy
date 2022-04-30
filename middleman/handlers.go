package middleman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BartoszBurgiel/cloud/shared"
	"github.com/google/uuid"
)

// handle a get request for a package
func (m *Middleman) handleGetPackage(w http.ResponseWriter, r *http.Request) {
	m.updateEnvVariables()
	id := uuid.New().String()
	logRequest(m.log, r.Method, id)
	if !r.URL.Query().Has("api_key") {
		logMissingAPIKey(m.log, id)
		serve401(w, r)
		return
	}
	if !m.verifyAPIKey([]byte(r.URL.Query().Get("api_key"))) {
		logFalseAPIKey(m.log, id)
		serve401(w, r)
		return
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()

	s := shared.NewSubmission(m.hmacSecret, m.currentPackage)
	fmt.Fprintf(w, string(s.JSON()))

	m.currentPackage = shared.Package{}
	return
}

// handle upload of a package
func (m *Middleman) handlePostPackage(w http.ResponseWriter, r *http.Request) {
	m.updateEnvVariables()
	id := uuid.New().String()
	logRequest(m.log, r.Method, id)

	// check if the API key was passed
	if !r.URL.Query().Has("api_key") {
		logMissingAPIKey(m.log, id)
		serve401(w, r)
		return
	}
	if !m.verifyAPIKey([]byte(r.URL.Query().Get("api_key"))) {
		logFalseAPIKey(m.log, id)
		serve401(w, r)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	submission := shared.NewEmptySubmission()
	err = json.Unmarshal(body, &submission)
	if err != nil {
		fmt.Fprintln(w, err, string(body))
		return
	}

	err = m.tryAcceptSubmission(submission)
	if err == nil {
		logPackageAccepted(m.log, submission.ID)
		fmt.Fprintln(w, SuccessfullUploadResponse)
		return
	}
	fmt.Fprintln(w, err)
}

func serve401(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	return
}
