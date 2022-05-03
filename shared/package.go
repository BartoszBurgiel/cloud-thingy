package shared

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Package represents the actual payload which will be sent to the middleman
// and forwarded to the reciever
type Package struct {

	// Payload holds the compressed and encrypted
	// files
	Payload []byte

	// CheckSum of the decrypted archive to verify if the
	// decryption was succesful
	CheckSum []byte
}

// NewPackage creates a new package of the compressed and encrypted payload with its checksum
func NewPackage(config PackageConfig, files []string) (Package, error) {
	l := log.New(os.Stdout, "PACKAGE>", log.Ltime)

	// id for logging purposes
	start := time.Now()

	t := time.Now()
	payload, err := preparePayload(files)
	logError(l, err)
	if err != nil {
		return Package{}, err
	}
	logArchiveResult(l, time.Since(t), len(payload))

	hasher := md5.New()
	io.Copy(hasher, bytes.NewReader(payload))
	checkSum := hasher.Sum(nil)

	t = time.Now()
	encryptedPayload, err := encryptPayload(config.KeyFilePath, payload)
	logEncryptionResult(l, time.Since(t))
	if err != nil {
		return Package{}, err
	}

	if config.MaxPackageSize < len(encryptedPayload) {
		return Package{}, fmt.Errorf(
			"The payload exceeds the size limit. Payload size: %d, Max size: %d\n",
			len(encryptedPayload),
			config.MaxPackageSize,
		)
	}

	p := Package{
		Payload:  encryptedPayload,
		CheckSum: checkSum,
	}
	logPackagePrepFinish(l, p, time.Since(start))
	return p, nil
}

// read and compress all of the files
func preparePayload(files []string) ([]byte, error) {

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for _, v := range files {

		f, err := os.Open(v)
		defer f.Close()
		if err != nil {
			return []byte{}, err
		}

		wr, err := w.Create(v)
		if err != nil {
			return []byte{}, err
		}

		_, err = io.Copy(wr, f)
		if err != nil {
			return []byte{}, err
		}
	}
	err := w.Close()
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func encryptPayload(saltFilePath string, archive []byte) ([]byte, error) {
	key, err := ioutil.ReadFile(saltFilePath)
	if err != nil {
		return []byte{}, err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	p := make([]byte, aes.BlockSize+len(archive))
	iv := p[:aes.BlockSize]

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return []byte{}, err
	}

	stream := cipher.NewCTR(c, iv)
	stream.XORKeyStream(p[aes.BlockSize:], archive)
	return p, nil
}

// Decrypt and verify the contents of the package.
// The method returns an error if there was an error while reading the
// SaltFile or when the decryption failed.
func (p *Package) Decrypt(conf PackageConfig) error {

	b := make([]byte, len(p.Payload))
	copy(b, p.Payload)

	key, err := ioutil.ReadFile(conf.KeyFilePath)
	if err != nil {
		return err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := p.Payload[:aes.BlockSize]
	p.Payload = p.Payload[aes.BlockSize:]

	stream := cipher.NewCTR(c, iv)
	stream.XORKeyStream(p.Payload, p.Payload)

	if !p.isPayloadValid() {
		p.Payload = b
		return FailedDecryptionOfThePackage
	}

	return nil
}

// isPayloadValid returns if the checksum of the current payload
// matches the original unencrypted payload
func (p Package) isPayloadValid() bool {
	hasher := md5.New()
	io.Copy(hasher, bytes.NewReader(p.Payload))
	return bytes.Equal(hasher.Sum(nil), p.CheckSum)
}
