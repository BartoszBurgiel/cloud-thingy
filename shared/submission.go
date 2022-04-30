package shared

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"

	"github.com/google/uuid"
)

// Submission represents a single submission of a package
// This struct is the payload of the actual POST request to the middleman
type Submission struct {
	MAC     []byte
	ID      string
	Package Package
}

// NewSubmission creates a new submission instance with the generated MAC based
// on the jsonified package
func NewSubmission(macSecret []byte, p Package) Submission {
	message, _ := json.Marshal(p)
	// generate the MAC
	macHasher := hmac.New(sha512.New, macSecret)
	macHasher.Write(message)
	MAC := macHasher.Sum(nil)
	return Submission{
		MAC:     MAC,
		ID:      uuid.New().String(),
		Package: p,
	}
}

// NewEmptySubmission returns a submission instance with uninitialized fields
func NewEmptySubmission() Submission {
	return Submission{}
}

// JSON returns the submission as a JSON object
func (s Submission) JSON() []byte {
	j, _ := json.Marshal(s)
	return j
}
