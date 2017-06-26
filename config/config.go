package config

import (
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

// DBUser is the username to connect to postgres with
var DBUser = "lambda_dev"

// DBPass is the password to connect to postgres with
var DBPass = "testing"

// DBHost is the hostname of the postgres server
var DBHost = "localhost"

// DBPort is the port postgres is running on
var DBPort uint16 = 5432

// DBName is the name of the postgres database to use
var DBName = "lambda_dev"

// MinifiedAssets is whether minified versions of js and css should be used
var MinifiedAssets bool // TODO implement

func init() {
	s, exists := os.LookupEnv("LMDA_RECAPTCHA_SECRET")
	if exists {
		RecaptchaSecret = s
	}

	s, exists = os.LookupEnv("LMDA_RECAPTCHA_SITE_KEY")
	if exists {
		RecaptchaSiteKey = s
	}

	s, exists = os.LookupEnv("LMDA_ALLOWED_FILETYPES")
	if exists {
		AllowedFiletypes = strings.Split(s, ",")
		AllowedFiletypesStr = s
	}

	s, exists = os.LookupEnv("LMDA_MAX_UPLOAD_SIZE")
	if exists {
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
		UploadDir = s
	}
}
