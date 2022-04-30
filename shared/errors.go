package shared

import "fmt"

var (
	MiddlemanPackageHasInvalidMAC   = fmt.Errorf("Middleman's response submission was sealed has invalid MAC.")
	DownloadReturnsEmptyPackage     = fmt.Errorf("The downloaded package was empty.")
	FailedDecryptionOfThePackage    = fmt.Errorf("The decryption of the payload was not successful.")
	MiddlemanHasAPackageInMemory    = fmt.Errorf("The middleman already has a package in the memory. The incoming package won't be accepted.")
	ChecksumHasAleadyBeenUploaded   = fmt.Errorf("A package with this checksum has already been uploaded before. The incoming package won't be accepted.")
	TheSubmissionIsNotAuthenticated = fmt.Errorf("The incoming package is not authenticated.")
)
