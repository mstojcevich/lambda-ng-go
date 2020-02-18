package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// RecaptchaEnabled is whether or not Recaptcha should be required for registration
var RecaptchaEnabled bool // TODO implement

// RecaptchaSecret is the secret key obtained from Recaptcha
var RecaptchaSecret string

// RecaptchaSiteKey is the site key obtained from Recaptcha
var RecaptchaSiteKey string

// AllowedFiletypesStr is the list of comma-separated allowed file extensions
var AllowedFiletypesStr = ".png,.jpg,.jpeg,.pdf,.zip,.7z,.mp3,.opus,.mp4,.webm,.webp,.gif,.ogg"

// AllowedFiletypes is the list of allowed file extensions
var AllowedFiletypes = []string{".png", ".jpg", ".jpeg", ".pdf", ".zip", ".7z", ".mp3", ".opus", ".mp4", ".webm", ".webp", ".gif", ".ogg"}

// ThumbnailExtensions is the list of extensions to create thumbnails for
var ThumbnailExtensions = []string{"png", "jpg", "jpeg", "gif"}

// MaxUploadSize is the maximum upload size (in MB)
var MaxUploadSize = 15

// UploadDir is the directory that uploads are stored in
var UploadDir = "files/"

// DBString is the string to use to connect to the database. It is the second argument to db.Connect
var DBString = "host=db port=5432 user=lambda_dev password=super-secret-dev-password dbname=lambda_dev sslmode=disable"

// MinifiedAssets is whether minified versions of js and css should be used
var MinifiedAssets bool // TODO implement

// ClamAVScanning is whether or not uploads should be scanned with ClamAV
var ClamAVScanning bool

// ClamSock is the path to the ClamAV socket file
var ClamSock = "/var/lib/clamav/clamd.sock"

// BackblazeBucket is the name of the B2 bucket to use for the B2 integration
var BackblazeBucket = "lambda"

// BackblazeAccountID is the Backblaze account ID for the B2 integration
var BackblazeAccountID = ""

// BackblazeAppKey is the Backblaze application key for the B2 integration
var BackblazeAppKey = ""

// Address used to connected to redis
var RedisAddr = "localhost:6379"

// Password for redis authentication. Empty string = no authentication.
var RedisPassword = ""

// String to specify the interface and port to listen on
var ListenStr = ":8080"

func init() {
	// Using environment vars for config isn't ideal, but we unset them ASAP
	// and using a file would require some way to drop privs after to mitigate
	// a path traversal vulnerability. This should be relatively safe.

	s, exists := os.LookupEnv("LMDA_RECAPTCHA_SECRET")
	if exists {
		os.Unsetenv("LMDA_RECAPTCHA_SECRET")
		RecaptchaSecret = s
	}

	s, exists = os.LookupEnv("LMDA_RECAPTCHA_SITE_KEY")
	if exists {
		os.Unsetenv("LMDA_RECAPTCHA_SITE_KEY")
		RecaptchaSiteKey = s
	}

	s, exists = os.LookupEnv("LMDA_ALLOWED_FILETYPES")
	if exists {
		os.Unsetenv("LMDA_ALLOWED_FILETYPES")
		AllowedFiletypes = strings.Split(s, ",")
		AllowedFiletypesStr = s
	}

	s, exists = os.LookupEnv("LMDA_MAX_UPLOAD_SIZE")
	if exists {
		os.Unsetenv("LMDA_MAX_UPLOAD_SIZE")
		maxUploadSize, err := strconv.Atoi(s)
		if err != nil {
			log.Println("Error when parsing LMDA_MAX_UPLOAD_SIZE")
			log.Println(err)
		} else {
			MaxUploadSize = maxUploadSize
		}
	}

	s, exists = os.LookupEnv("LMDA_UPLOAD_DIR")
	if exists {
		os.Unsetenv("LMDA_UPLOAD_DIR")
		UploadDir = s
	}

	s, exists = os.LookupEnv("LMDA_DB_CONNSTR")
	if exists {
		os.Unsetenv("LMDA_DB_CONNSTR")
		DBString = s
	}

	s, exists = os.LookupEnv("LMDA_CLAMAV")
	if exists {
		os.Unsetenv("LMDA_CLAMAV")
		cav, err := strconv.ParseBool(s)
		if err != nil {
			fmt.Println("Error parsing LMDA_CLAMAV")
			fmt.Println(err)
		} else {
			ClamAVScanning = cav
		}
	}

	s, exists = os.LookupEnv("LMDA_CLAM_SOCK")
	if exists {
		os.Unsetenv("LMDA_CLAM_SOCK")
		ClamSock = s
	}

	s, exists = os.LookupEnv("LMDA_BLAZE_ID")
	if exists {
		os.Unsetenv("LMDA_BLAZE_ID")
		BackblazeAccountID = s
	}

	s, exists = os.LookupEnv("LMDA_BLAZE_KEY")
	if exists {
		os.Unsetenv("LMDA_BLAZE_KEY")
		BackblazeAppKey = s
	}

	s, exists = os.LookupEnv("LMDA_BLAZE_BUCKET")
	if exists {
		os.Unsetenv("LMDA_BLAZE_BUCKET")
		BackblazeBucket = s
	}

	s, exists = os.LookupEnv("LMDA_REDIS_ADDR")
	if exists {
		os.Unsetenv("LMDA_REDIS_ADDR")
		RedisAddr = s
	}

	s, exists = os.LookupEnv("LMDA_REDIS_PASS")
	if exists {
		os.Unsetenv("LMDA_REDIS_PASS")
		RedisPassword = s
	}

	s, exists = os.LookupEnv("LMDA_LISTEN_STR")
	if exists {
		os.Unsetenv("LMDA_LISTEN_STR")
		ListenStr = s
	}
}
