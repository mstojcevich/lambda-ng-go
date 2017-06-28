package user

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"regexp"
	"time"

	"github.com/dchest/uniuri"
	"github.com/mstojcevich/lambda-ng-go/config"
	"github.com/mstojcevich/lambda-ng-go/database"
	"github.com/mstojcevich/lambda-ng-go/recaptcha"
	tplt "github.com/mstojcevich/lambda-ng-go/template"
	"github.com/mstojcevich/lambda-ng-go/user/session"
	"github.com/valyala/fasthttp"
)

var captcha = recaptcha.NewInstance(config.RecaptchaSecret)

var checkUserStmt, checkUserErr = database.DB.Prepare(`SELECT exists(SELECT 1 FROM users WHERE username=$1)`)
var createUserStmt, createUserErr = database.DB.Prepare(`INSERT INTO users (username, password, creation_date, api_key, encryption_enabled) VALUES ($1, $2, $3, $4, false)`)

var isAlnum = regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString

type registerTplContext struct {
	tplt.CommonTemplateCtx

	RecaptchaSiteKey string
}

func init() {
	createRegisterTemplate()
}

func createRegisterTemplate() {
	// Create the template
	t, err := template.ParseFiles("html/register.html", "html/partials/shared_head.html")
	if err != nil {
		panic(err)
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = t.Execute(&tpl, registerTplContext{RecaptchaSiteKey: config.RecaptchaSiteKey})
	if err != nil {
		panic(err)
	}

	// Output the template to a file
	ioutil.WriteFile("html/compiled/register.html", tpl.Bytes(), 0644)
}

// LoginPage renders the login page HTML
func RegisterPage(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/register.html")
}

// RegisterAPI handles a user registration, creating the user if successful
func RegisterAPI(ctx *fasthttp.RequestCtx) {
	username := string(ctx.FormValue("username"))
	password := string(ctx.FormValue("password"))
	captchaResponse := string(ctx.FormValue("g-recaptcha-response"))
	remoteIP := ctx.RemoteIP()

	captchaValid, err := captcha.Check(remoteIP.String(), captchaResponse)

	if err != nil {
		// TODO proper logging
		log.Println(err)

		registerError(ctx, "Failed to check captcha", fasthttp.StatusInternalServerError)
		return
	}

	if !isAlnum(username) {
		registerError(ctx, "Username can only contain English letters and numbers", fasthttp.StatusUnprocessableEntity)
		return
	}

	if len(username) > MaxUsernameLength {
		registerError(ctx, fmt.Sprintf("Username length > %d", MaxUsernameLength), fasthttp.StatusUnprocessableEntity)
		return
	}

	if len(username) < MinUsernameLength {
		registerError(ctx, fmt.Sprintf("Username length < %d", MinUsernameLength), fasthttp.StatusUnprocessableEntity)
		return
	}

	if len(password) > MaxPassLength {
		registerError(ctx, fmt.Sprintf("Password length > %d", MaxPassLength), fasthttp.StatusUnprocessableEntity)
		return
	}

	if len(password) < MinPassLength {
		registerError(ctx, fmt.Sprintf("Password length < %d", MinPassLength), fasthttp.StatusUnprocessableEntity)
		return
	}

	if !captchaValid {
		registerError(ctx, "Invalid captcha response", fasthttp.StatusUnprocessableEntity)
		return
	}

	// Check if the username is already in use
	var exists bool
	err = checkUserStmt.QueryRow(username).Scan(&exists)

	if err != nil {
		// TODO log error

		registerError(ctx, "Failed to check if username was in use", fasthttp.StatusInternalServerError)
		return
	}

	if exists {
		registerError(ctx, "Username already in use", fasthttp.StatusUnprocessableEntity)
		return
	}

	hashedPw := HashPassword(password)

	apiKey := genAPIKey()

	_, err = createUserStmt.Exec(username, hashedPw, time.Now(), apiKey)
	if err != nil {
		log.Println(err)

		registerError(ctx, "Error creating user. Please try again.", fasthttp.StatusInternalServerError)
		return
	}

	// Sign the user in
	sess := session.Sessions.StartFasthttp(ctx)
	sess.Set("api_key", apiKey)

	// Return the result as JSON
	result := RegisterResult{Errors: nil, APIKey: apiKey, Success: true}
	resultJSON, err := result.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.Error("{errors:[\"Failed to create JSON response. Contact an admin\"]}", fasthttp.StatusInternalServerError)
		ctx.SetContentType("text/json")
		return
	}
	ctx.SetContentType("text/json")
	ctx.Write(resultJSON)
}

// registerError writes out JSON for a failed registration
func registerError(ctx *fasthttp.RequestCtx, errStr string, statusCode int) {
	result := RegisterResult{Errors: []string{errStr}, Success: false}
	resultJSON, err := result.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.Error("{errors:[\"Failed to create JSON response. Contact an admin\"]}", fasthttp.StatusInternalServerError)
		ctx.SetContentType("text/json")
		return
	}
	ctx.Error(string(resultJSON), statusCode)
	ctx.SetContentType("text/json")
}

func genAPIKey() string {
	return uniuri.NewLen(32)
}
