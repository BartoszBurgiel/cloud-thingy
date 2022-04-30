package reciever

import (
	"log"
	"time"
)

func logError(l *log.Logger, err error) {
	if err != nil {
		l.Println(err.Error())
	}
}

func logInit(l *log.Logger, dest string) {
	l.Printf("Starting a new reciever instance. Destination path: %s\n.", dest)
}

func logAuthentication(l *log.Logger, success bool, id string) {

	if success {
		l.Printf("The recieved package (id: %s) is authenticated.\n", id)
		return
	}

	l.Printf("The recieved package (id: %s) is not authenticated. Terminating handling.\n", id)
}

func logDecryption(l *log.Logger, success bool, id string) {

	if success {
		l.Printf("The recieved package (id: %s) could be successfully decrypted.\n", id)
		return
	}

	l.Printf("The recieved package (id: %s) could not successfully decrypted. Terminating handling.\n", id)
}

func logMiddlemanResponse(l *log.Logger, responseCode string, dur time.Duration) {
	l.Printf("Middleman responded with the code: %s after: %s\n", responseCode, dur)
}

func logRequest(l *log.Logger, apikey string) {
	l.Printf("Sending request to the middleman for a package with the following API key: %s\n", apikey)
}

func logSuccessfullDownload(l *log.Logger, id string) {
	l.Printf("The download of the package (id: %s), was successful.\n", id)
}
