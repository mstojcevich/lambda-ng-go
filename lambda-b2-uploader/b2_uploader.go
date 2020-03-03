// Uploads files to backblaze b2 at regular intervals

package main

import (
	"bufio"
	"database/sql"
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Needed for postgres DB
	"gopkg.in/kothar/go-backblaze.v0"
)

var getFilenamesStmt *sql.Stmt

var dbConnString string

var blazeAppID string
var blazeAppKey string
var blazeBucketID string

var uploadDir string

func main() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}

	dbc := connectToDB()
	prepareStatements(dbc)

	for {
		filesToUpload, err := getNonB2Filenames(dbc)
		if err != nil {
			log.Printf("Failed to get non-b2 filenames: %s\n", err)
			time.Sleep(2 * time.Minute)
			continue
		}

		// Hacky workaround to avoid uploading a file before thumbnail generation and clamav scanner are done
		time.Sleep(30 * time.Second)

		log.Printf("Uploading %d files to B2...\n", len(filesToUpload))
		startTime := time.Now()

		uploadDone := make(chan bool)
		fileChunks := chunkifyFiles(filesToUpload, 8)
		for _, chunk := range fileChunks {
			go func(chunk []string) {
				err := uploadFiles(dbc, chunk)
				if err != nil {
					log.Printf("Error uploading chunk of files: %s\n", err)
				}
				uploadDone <- true
			}(chunk)
		}
		for i := 0; i < len(fileChunks); i++ {
			<-uploadDone
		}

		endTime := time.Now()
		elapsed := endTime.Sub(startTime)

		log.Printf("File uploading took %.1f minutes. Going to sleep...\n", elapsed.Minutes())
		time.Sleep(6 * time.Hour)
	}
}

func loadConfig() error {
	var exists bool

	blazeAppID, exists = os.LookupEnv("LMDA_BLAZE_APP_ID")
	os.Unsetenv("LMDA_BLAZE_APP_ID")
	if !exists {
		return errors.New("Missing LMDA_BLAZE_APP_ID environment variable")
	}

	blazeAppKey, exists = os.LookupEnv("LMDA_BLAZE_KEY")
	os.Unsetenv("LMDA_BLAZE_KEY")
	if !exists {
		return errors.New("Missing LMDA_BLAZE_KEY environment variable")
	}

	blazeBucketID, exists = os.LookupEnv("LMDA_BLAZE_BUCKET")
	os.Unsetenv("LMDA_BLAZE_BUCKET")
	if !exists {
		return errors.New("Missing LMDA_BLAZE_BUCKET environment variable")
	}

	dbConnString, exists = os.LookupEnv("LMDA_DB_CONNSTR")
	os.Unsetenv("LMDA_DB_CONNSTR")
	if !exists {
		return errors.New("Missing LMDA_DB_CONNSTR environment variable")
	}

	uploadDir, exists = os.LookupEnv("LMDA_UPLOAD_DIR")
	os.Unsetenv("LMDA_UPLOAD_DIR")
	if !exists {
		return errors.New("Missing LMDA_UPLOAD_DIR environment variable")
	}

	return nil
}

func prepareStatements(dbc *sqlx.DB) {
	var err error
	getFilenamesStmt, err = dbc.Prepare("SELECT CONCAT(name, '.', extension) AS filename FROM files WHERE in_b2=false")
	if err != nil {
		panic(err)
	}
}

func getNonB2Filenames(dbc *sqlx.DB) ([]string, error) {
	rows, err := getFilenamesStmt.Query()
	if err != nil {
		return []string{}, err
	}
	var filenames = make([]string, 0)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		filenames = append(filenames, name)
	}
	return filenames, nil
}

func connectToDB() *sqlx.DB {
	// Open instead of connect to more gracefully handle the DB being initially unavailable.
	db, err := sqlx.Open("postgres", dbConnString)

	// Error establishing connection to DB
	if err != nil {
		panic(err)
	}

	for {
		err = db.Ping()
		if err == nil {
			break
		} else {
			log.Printf("DB unavailable. Trying again... %s\n", err)
			time.Sleep(5 * time.Second)
		}
	}

	return db
}

func chunkifyFiles(files []string, targetChunkCount int) [][]string {
	fileChunks := make([][]string, 0, targetChunkCount)
	chunkSize := (len(files) + targetChunkCount - 1) / targetChunkCount
	for chunkStart := 0; chunkStart < len(files); chunkStart += chunkSize {
		chunkEnd := chunkStart + chunkSize
		if chunkEnd > len(files) {
			chunkEnd = len(files)
		}
		fileChunks = append(fileChunks, files[chunkStart:chunkEnd])
	}
	return fileChunks
}

func createB2Client() (*backblaze.Bucket, error) {
	b2, err := backblaze.NewB2(backblaze.Credentials{
		KeyID:          blazeAppID,
		ApplicationKey: blazeAppKey,
	})
	if err != nil {
		return nil, err
	}

	bucket, err := b2.Bucket(blazeBucketID)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func uploadFiles(dbc *sqlx.DB, files []string) error {
	b2bucket, err := createB2Client()
	if err != nil {
		panic(err)
	}

	fileNames := make([]string, len(files))
	for i, filename := range files {
		file, err := os.Open(path.Join(uploadDir, filename))
		if err != nil {
			log.Printf("Skipping %s\n", filename)
			continue
		}
		defer file.Close()

		fileReader := bufio.NewReader(file)
		_, err = b2bucket.UploadFile(filename, make(map[string]string, 0), fileReader)
		if err != nil {
			return err
		}

		fileNames[i] = strings.Split(filename, ".")[0]
	}

	query, args, err := sqlx.In("UPDATE files SET in_b2=true WHERE name IN (?);", fileNames)
	if err != nil {
		return err
	}
	query = dbc.Rebind(query)
	_, err = dbc.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
