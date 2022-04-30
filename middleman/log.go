package middleman

import "log"

func logError(l *log.Logger, err error) {
	if err != nil {
		l.Println(err.Error())
	}
}

func logStart(l *log.Logger) {
	l.Println("Starting the middleman.")
}

func logRequest(l *log.Logger, method, id string) {
	l.Printf("New %s request (id: %s).\n", method, id)
}
func logMissingAPIKey(l *log.Logger, id string) {
	l.Printf("Request (id: %s) has no API key.\n", id)
}
func logFalseAPIKey(l *log.Logger, id string) {
	l.Printf("Request (id: %s) has an invalid API key.\n", id)
}

func logPackageAleadyInMemory(l *log.Logger, id string) {
	l.Printf("Rejecting the submission (id: %s), as there is a package in memory already.\n", id)
}

func logPackageAlreadyUploaded(l *log.Logger, id string) {
	l.Printf("Rejecting the submission (id: %s), as it was uploaded already before.\n", id)
}

func logInvalidMAC(l *log.Logger, submissionID string) {
	l.Printf("The submission (id: %s) is not authenticated, i.e. it has an invalid MAC\n", submissionID)
}

func logPackageAccepted(l *log.Logger, id string) {
	l.Printf("The submission (id: %s) has been accepted.\n", id)
}

func logVariableUpdate(l *log.Logger, m Middleman) {
	l.Printf("The environment variables have been updated: api_key_hash: %s, mac_secret: %s, port: %s\n",
		string(m.apiKeyHash), string(m.hmacSecret), m.port)
}
