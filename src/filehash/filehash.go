package filehash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func getCheckSum(fileName string) string {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hasher := sha256.New()

	_, err = io.Copy(hasher, file)
	if err != nil {
		panic(err)
	}

	sum := hasher.Sum(nil)

	checksum := hex.EncodeToString(sum)

	return checksum
}
