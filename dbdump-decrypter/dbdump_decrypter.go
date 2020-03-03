package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/blake2b"
)

func main() {
	secretKey, exists := os.LookupEnv("LMDA_ENCRYPTION_KEY")
	os.Unsetenv("LMDA_ENCRYPTION_KEY")
	if !exists {
		panic("LMDA_ENCRYPTION_KEY not defined")
	}
	if len(secretKey) != 64 {
		panic("Encryption key must be 64 characters long")
	}
	secretKeyBytes, err := hex.DecodeString(secretKey)
	if err != nil {
		panic("Encryption key must be hex")
	}

	macKey, exists := os.LookupEnv("LMDA_MAC_KEY")
	os.Unsetenv("LMDA_MAC_KEY")
	if !exists {
		panic("LMDA_MAC_KEY not defined")
	}
	if len(macKey) != 64 {
		// It doesn't NEED to be this long, but it always will be.
		panic("MAC key must be 64 characters long")
	}
	macKeyBytes, err := hex.DecodeString(macKey)
	if err != nil {
		panic("MAC key must be hex")
	}

	if subtle.ConstantTimeCompare(secretKeyBytes, macKeyBytes) == 1 {
		panic("MAC key and encryption key must be different")
	}

	header := make([]byte, 6)
	_, err = io.ReadFull(os.Stdin, header)
	if err != nil {
		panic(err)
	}
	if string(header) != "LMDA00" {
		panic("Unsupported file type")
	}

	fileMAC := make([]byte, 64)
	_, err = io.ReadFull(os.Stdin, fileMAC)
	if err != nil {
		panic(err)
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(os.Stdin, iv)
	if err != nil {
		panic(err)
	}

	// Save the ciphertext into a temporary file to ensure that it's seekable
	ciphertextFile, err := ioutil.TempFile("", "lmda_dbdump_ciphertext")
	if err != nil {
		panic(err)
	}
	defer ciphertextFile.Close()

	// The decryption is done in two passes so that stdout isn't written
	// to until we are done verifying the integrity of the ciphertext.

	// First pass: verify the MAC to ensure integrity and copy to the tempfile as we go
	mac, err := blake2b.New512(macKeyBytes)
	if err != nil {
		panic(err)
	}
	_, err = mac.Write(iv)
	if err != nil {
		panic(err)
	}

	ciphertextBuf := make([]byte, 4096)
	for {
		n, err := io.ReadFull(os.Stdin, ciphertextBuf)
		isEOF := err == io.EOF
		if n == 0 && isEOF {
			break
		}
		if err != nil && !isEOF {
			panic(err)
		}
		_, err = mac.Write(ciphertextBuf)
		if err != nil {
			panic(err)
		}
		_, err = ciphertextFile.Write(ciphertextBuf)
		if err != nil {
			panic(err)
		}
	}

	// Verify the MAC
	actualMAC := mac.Sum(nil)
	if subtle.ConstantTimeCompare(actualMAC, fileMAC) == 0 {
		panic("Integrity check failed")
	}

	_, err = ciphertextFile.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}

	// Second pass: decrypt the ciphertext to stdout
	blockCipher, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		panic(err)
	}
	aesctr := cipher.NewCTR(blockCipher, iv)

	stdoutWriter := bufio.NewWriter(os.Stdout)
	plaintextBuf := make([]byte, 4096)
	for {
		_, err = io.ReadFull(ciphertextFile, ciphertextBuf)
		isEOF := err == io.EOF
		if err != nil && !isEOF {
			panic(err)
		}
		aesctr.XORKeyStream(plaintextBuf, ciphertextBuf)
		hadZero := false
		for _, b := range plaintextBuf {
			if b == 0 {
				hadZero = true
				break
			}
			err := stdoutWriter.WriteByte(b)
			if err != nil {
				panic(err)
			}
		}
		if hadZero || isEOF {
			break
		}
	}
	err = stdoutWriter.Flush()
	if err != nil {
		panic(err)
	}
}
