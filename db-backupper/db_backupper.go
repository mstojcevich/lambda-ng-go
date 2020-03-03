package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"golang.org/x/crypto/blake2b"
	"gopkg.in/kothar/go-backblaze.v0"
)

// Dumps the database, encrypts the dump, and uploads it to a B2 bucket.
// Expects a .pgpass file containing the DB password.
// Uses PGDATABASE, PGHOST, PGOPTIONS, PGPORT, PGUSER,
// LMDA_MAC_KEY, and LMDA_DB_ENCRYPTION_KEY
// environment vars for configuration.

// Encrypted database file structure:
// [6 bytes - file version (ASCII) - e.g. LMDA00]
// [64 bytes - BLAKE2b-512 MAC of IV + ciphertext]
// [16 bytes - AES CTR IV]
// [16 * n bytes - AES-256 CTR ciphertext, terminated with a null-byte]
// [n bytes - random garbage (padding due to buffer size), ignore anything after a null byte]

func main() {
	secretKey, exists := os.LookupEnv("LMDA_ENCRYPTION_KEY")
	os.Unsetenv("LMDA_ENCRYPTION_KEY")
	if !exists {
		log.Fatalln("LMDA_ENCRYPTION_KEY not defined")
	}
	if len(secretKey) != 64 {
		log.Fatalln("Encryption key must be 64 characters long")
	}
	secretKeyBytes, err := hex.DecodeString(secretKey)
	if err != nil {
		log.Fatalln("Encryption key must be hex")
	}

	macKey, exists := os.LookupEnv("LMDA_MAC_KEY")
	os.Unsetenv("LMDA_MAC_KEY")
	if !exists {
		log.Fatalln("LMDA_MAC_KEY not defined")
	}
	if len(macKey) != 64 {
		// It doesn't NEED to be this long, but it always will be.
		log.Fatalln("MAC key must be 64 characters long")
	}
	macKeyBytes, err := hex.DecodeString(macKey)
	if err != nil {
		log.Fatalln("MAC key must be hex")
	}

	if subtle.ConstantTimeCompare(secretKeyBytes, macKeyBytes) == 1 {
		log.Fatalln("MAC key and encryption key must be different")
	}

	blazeAppID, exists := os.LookupEnv("LMDA_BLAZE_APP_ID")
	os.Unsetenv("LMDA_BLAZE_APP_ID")
	if !exists {
		log.Fatalln("Missing LMDA_BLAZE_APP_ID environment variable")
	}

	blazeAppKey, exists := os.LookupEnv("LMDA_BLAZE_KEY")
	os.Unsetenv("LMDA_BLAZE_KEY")
	if !exists {
		log.Fatalln("Missing LMDA_BLAZE_KEY environment variable")
	}

	blazeBucketID, exists := os.LookupEnv("LMDA_BLAZE_BUCKET")
	os.Unsetenv("LMDA_BLAZE_BUCKET")
	if !exists {
		log.Fatalln("Missing LMDA_BLAZE_BUCKET environment variable")
	}

	log.Println("Waiting a bit to give the DB a chance to come up...")
	time.Sleep(30 * time.Second)
	for {
		log.Println("Creating a new DB dump")
		err = createDump(
			secretKeyBytes, macKeyBytes, blazeAppID, blazeAppKey, blazeBucketID)
		if err != nil {
			log.Printf("Failed to create DB dump: %s\n", err)
			log.Println("Going to sleep for 5 minutes")
			time.Sleep(5 * time.Minute)
		} else {
			log.Println("Successfully created and uploaded DB dump")
			log.Println("Going to sleep for 6 hours")
			time.Sleep(6 * time.Hour)
		}
	}
}

func createB2Client(appID, appKey, bucketID string) (*backblaze.Bucket, error) {
	b2, err := backblaze.NewB2(backblaze.Credentials{
		KeyID:          appID,
		ApplicationKey: appKey,
	})
	if err != nil {
		return nil, err
	}

	bucket, err := b2.Bucket(bucketID)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func createDump(secretKeyBytes, macKeyBytes []byte, blazeAppID, blazeAppKey, blazeBucketID string) error {
	encryptedDumpFile, err := ioutil.TempFile("", "lmda_dbdump_")
	if err != nil {
		return err
	}
	defer encryptedDumpFile.Close()

	// Plaintext header including the dump file version number.
	_, err = encryptedDumpFile.Write([]byte("LMDA00"))
	if err != nil {
		return err
	}

	// Create some room for the MAC. We'll seek back later to write the real MAC.
	_, err = encryptedDumpFile.Write(make([]byte, blake2b.Size))
	if err != nil {
		return err
	}

	mac, err := blake2b.New512(macKeyBytes)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil { // (len == BlockSize) iff (err == nil)
		return err
	}

	// Include the IV in the MAC
	_, err = mac.Write(iv)
	if err != nil {
		return err
	}

	// Put the IV right before the ciphertext
	_, err = encryptedDumpFile.Write(iv)
	if err != nil {
		return err
	}

	blockCipher, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		return err
	}
	aesctr := cipher.NewCTR(blockCipher, iv)

	// Not using the custom dump format since it only saves a bit of space
	// and makes it harder to do padding (it's a binary format that can contain
	// zeroes, so we'd have to indicate the number of padding bytes somewhere).
	pgdumpCmd := exec.Command("pg_dump", "--no-password")
	pgdumpOutput, err := pgdumpCmd.StdoutPipe()
	if err != nil {
		return err
	}
	pgdumpErr, err := pgdumpCmd.StderrPipe()
	if err != nil {
		return err
	}
	err = pgdumpCmd.Start()
	if err != nil {
		return err
	}

	go func() {
		stderr, err := ioutil.ReadAll(pgdumpErr)
		if err != nil {
			log.Println(string(stderr))
			return
		}
		log.Println(string(stderr))
	}()

	plaintextBuf := make([]byte, 4096)
	ciphertextBuf := make([]byte, 4096)
	for {
		n, err := io.ReadFull(pgdumpOutput, plaintextBuf)
		isEOF := err == io.ErrUnexpectedEOF || err == io.EOF
		if err != nil && !isEOF {
			return err
		}
		if n == 0 && isEOF {
			// Don't update the MAC with nothing (it actually changes it?)
			break
		}

		if n < len(plaintextBuf) {
			// Pad with a null byte followed by random bytes
			plaintextBuf[n] = 0
			if n < len(plaintextBuf)-1 {
				// There's more stuff to pad!
				_, err = rand.Read(plaintextBuf[n+1:])
				if err != nil {
					return err
				}
			}
		}
		aesctr.XORKeyStream(ciphertextBuf, plaintextBuf)
		_, err = mac.Write(ciphertextBuf)
		if err != nil {
			return err
		}
		_, err = encryptedDumpFile.Write(ciphertextBuf)
		if err != nil {
			return err
		}
		if isEOF {
			break
		}
	}

	err = pgdumpCmd.Wait()
	if err != nil {
		return err
	}

	// Go back and write the MAC output.
	_, err = encryptedDumpFile.Seek(6, io.SeekStart)
	if err != nil {
		return err
	}
	macOutput := mac.Sum(nil)
	_, err = encryptedDumpFile.Write(macOutput)
	if err != nil {
		return err
	}

	// Upload the encrypted db dump file to B2
	bucket, err := createB2Client(blazeAppID, blazeAppKey, blazeBucketID)
	if err != nil {
		return err
	}
	encryptedDumpFile.Seek(0, io.SeekStart)
	_, err = bucket.UploadFile(
		"lmda_dbdump",
		make(map[string]string, 0),
		bufio.NewReader(encryptedDumpFile),
	)
	if err != nil {
		return err
	}

	return nil
}
