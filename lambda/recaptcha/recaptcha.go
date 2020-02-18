package recaptcha

import (
	"time"

	"github.com/valyala/fasthttp"
)

// Client for recaptcha using fasthttp and easyjson
// Based on github.com/dpapathanasiou/go-recaptcha

// verifyURL is the URL used to verify a captcha response
const verifyURL = "https://www.google.com/recaptcha/api/siteverify"

// recaptchaResponse is a response recieved from Recaptcha's validation server
type recaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// Instance is an instance of a recaptcha client
type Instance struct {
	secretKey string
}

// NewInstance creates a new instance of a Recaptcha client
// with the specified secret key
func NewInstance(secretKey string) Instance {
	return Instance{secretKey}
}

// check validates a client's challenge-response response with the recatcha server
func (instance Instance) check(remoteIP, response string) (r recaptchaResponse, e error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(verifyURL)

	req.Header.SetMethod("POST")

	// We're expecting a JSON response back from Google
	req.Header.Add("Accept", "application/json")

	// Set the POST arguments
	pa := req.PostArgs()
	pa.Add("secret", instance.secretKey)
	pa.Add("remoteip", remoteIP)
	pa.Add("response", response)

	// Set the body of the request to the POST args
	req.SetBodyString(pa.String())

	// Make the request
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{
		ReadTimeout:        1 * time.Minute,
		WriteTimeout:       1 * time.Minute,
	}
	err := client.Do(req, resp)

	if err != nil {
		return recaptchaResponse{}, err
	}

	// Parse the JSON resposne
	bodyBytes := resp.Body()
	responseJSON := recaptchaResponse{}
	err = responseJSON.UnmarshalJSON(bodyBytes)

	if err != nil {
		return recaptchaResponse{}, err
	}

	return responseJSON, nil
}

// Check validates a client's challenge-response response with the recatcha server
func (instance Instance) Check(remoteIP, response string) (bool, error) {
	resp, err := instance.check(remoteIP, response)

	if err != nil {
		return false, err
	}

	return resp.Success, nil
}
