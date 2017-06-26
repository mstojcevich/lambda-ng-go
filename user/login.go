package user

import (
	"bytes"
	"html/template"
	"log"

	"github.com/valyala/fasthttp"

	"io/ioutil"

	"github.com/mstojcevich/lambda-ng-go/database"
	"github.com/mstojcevich/lambda-ng-go/user/session"
)

var userLoginStmt, _ = database.DB.Preparex(`SELECT password, api_key FROM users WHERE username=$1`)
var userByKeyStmt, _ = database.DB.Prepare(`SELECT id, username FROM users WHERE api_key=$1`)

func init() {
	createLoginTemplate()
}

func createLoginTemplate() {
	// Create the template
	t, err := template.ParseFiles("html/login.html", "html/partials/shared_head.html")
	if err != nil {
		panic(err)
	}

	// Render the template into a byte buffer
	var tpl bytes.Buffer
	err = t.Execute(&tpl, nil)
	if err != nil {
		panic(err)
	}

	// Output the template to a file
	ioutil.WriteFile("html/compiled/login.html", tpl.Bytes(), 0644)
}

// LoginPage renders the login page HTML
func LoginPage(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("html/compiled/login.html")
}

// LoginAPI handles a user login
func LoginAPI(ctx *fasthttp.RequestCtx) {
	username := string(ctx.FormValue("username"))
	password := string(ctx.FormValue("password"))

	var user User
	err := userLoginStmt.Get(&user, username)
	if err != nil {
		loginError(ctx, "No such user exists", fasthttp.StatusPreconditionFailed)
		return
	}

	correctPW, err := CheckPassword(password, user.Password)
	if err != nil {
		log.Println(err)

		loginError(ctx, "Error checking password", fasthttp.StatusInternalServerError)
		return
	}
	if !correctPW {
		loginError(ctx, "Incorrect password", fasthttp.StatusUnauthorized)
		return
	}

	// Set the API key on the user's session
	sess := session.Sessions.StartFasthttp(ctx)
	sess.Set("api_key", user.APIKey)

	result := LoginResult{Errors: nil, APIKey: user.APIKey, Success: true}
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

// GetSessionAPI provides information about the current logged in session
func GetSessionAPI(ctx *fasthttp.RequestCtx) {
	user, err := GetLoggedInUser(ctx)
	if err != nil {
		sessionGetError(ctx, "User not logged in or no API key sent", fasthttp.StatusUnauthorized)
		return
	}

	response := SessionResult{UserID: user.ID, Username: user.Username, APIKey: user.APIKey}
	responseJSON, err := response.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.Error("{errors:[\"Failed to create JSON response. Contact an admin\"]}", fasthttp.StatusInternalServerError)
		ctx.SetContentType("text/json")
		return
	}

	ctx.Write(responseJSON)
	ctx.SetContentType("text/json")
}

// GetLoggedInUser gets the current logged in user by checking the API key both in POST and in the session
func GetLoggedInUser(ctx *fasthttp.RequestCtx) (User, error) {
	// Try to get API key first via session then via POST
	sess := session.Sessions.StartFasthttp(ctx)
	apiKey := sess.GetString("api_key")
	if apiKey == "" {
		apiKey = string(ctx.FormValue("api_key"))
	}

	var user User
	user.APIKey = apiKey
	row := userByKeyStmt.QueryRow(apiKey)
	err := row.Scan(&user.ID, &user.Username)

	return user, err
}

func LogoutAPI(ctx *fasthttp.RequestCtx) {
	// Logout the user by destroying their session
	session.Sessions.DestroyFasthttp(ctx)
}

// loginError writes out JSON for a failed login
func loginError(ctx *fasthttp.RequestCtx, errStr string, statusCode int) {
	result := LoginResult{Errors: []string{errStr}, Success: false}
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

// sessionGetError writes out JSON for when a session info request fails
func sessionGetError(ctx *fasthttp.RequestCtx, errStr string, statusCode int) {
	result := SessionResult{Errors: []string{errStr}}
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
