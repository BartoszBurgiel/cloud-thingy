package shared

import (
	"encoding/hex"
	"log"
	"time"
)

func logError(l *log.Logger, err error) {
	if err != nil {
		l.Println(err.Error())
	}
}

func logStartingPreparingArchive(l *log.Logger) {
	l.Printf("The preparation of the archive has begun at %s.\n", time.Now())
}

func logArchiveResult(l *log.Logger, dur time.Duration, size int) {
	l.Printf("The creation of the archive took %s, size of the archive in bytes: %d\n", dur, size)
}

func logInitPackage(l *log.Logger, c PackageConfig) {
	l.Printf("The creation of the package has started at %s, max payload size: %d\n", time.Now(), c.MaxPackageSize)
}

func logEncryptionResult(l *log.Logger, dur time.Duration) {
	l.Printf("The encryption of the archive took %s.\n", dur)
}

func logPackagePrepFinish(l *log.Logger, p Package, dur time.Duration) {
	l.Printf("The preparation of the package is finished and took %s. Package checksum: %s\n", dur, hex.EncodeToString(p.CheckSum))
}
