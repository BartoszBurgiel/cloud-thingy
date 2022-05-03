package client

import (
	"log"
	"time"

	"github.com/BartoszBurgiel/cloud/shared"
)

func logError(l *log.Logger, err error) {
	if err != nil {
		l.Println(err.Error())
	}
}

func logInitClient(l *log.Logger, path string) {
	l.Printf("Starting a new client instance at %s for: %s\n.", time.Now(), path)
}

func logAddingFileToArchive(l *log.Logger, name string) {
	l.Printf("Adding %s to the list of files to compress.\n", name)
}

func logSubmissionPreparationResult(l *log.Logger, s shared.Submission) {
	l.Printf("Submission (id: %s) is ready to be sent to the middleman.\n", s.ID)
}

func logRequestToTheMiddleman(l *log.Logger, apikey string, dur time.Duration, msg string) {
	l.Printf("The submission has been uploaded to the middleman. Used apikey: %s, upload took: %s, middleman's response: %s", apikey, dur, msg)
}
