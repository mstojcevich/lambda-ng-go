package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/mstojcevich/lambda-ng-go/config"
	"github.com/valyala/fasthttp"
)

// Cookie expiration
const sessionDuration = 86400 * 30 // 30 days

var redisClient = redis.NewClient(&redis.Options{
	Addr:     config.RedisAddr,
	Password: config.RedisPassword,
})

func generateSessionID() (string, error) {
	buf := make([]byte, 36, 36)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buf), nil
}

// Convert a session ID to a redis key used to fetch the session
func sessionIDToKey(sessionID string) string {
	return "lmda_session_v2|" + sessionID
}

// Check if a session ID is correctly formatted.
func validateSessionID(sessionID []byte) bool {
	// Session IDs are valid base64 and are at least 36 bytes (48 base64 chars)
	rawBytes := make([]byte, base64.URLEncoding.DecodedLen(len(sessionID)))
	_, err := base64.URLEncoding.Strict().Decode(rawBytes, sessionID)
	return err == nil && len(sessionID) >= 48
}

// Gets the session ID from the request. Returns an empty string if no session was set.
func getSessionIDFromRequest(ctx *fasthttp.RequestCtx) (string, error) {
	cookieBytes := ctx.Request.Header.Cookie("__Secure-lambdasession2")
	if cookieBytes == nil {
		return "", nil
	}

	cookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(cookie)
	err := cookie.ParseBytes(cookieBytes)
	if err != nil {
		return "", err
	}

	sessionIDBytes := cookie.Value()
	if !validateSessionID(sessionIDBytes) {
		return "", errors.New("Invalid session cookie")
	}

	return string(sessionIDBytes), nil
}

// LoadFromContext loads the session for a request, or returns nil if there is no session.
func LoadFromContext(ctx *fasthttp.RequestCtx) (*LambdaSession, error) {
	sessionID, err := getSessionIDFromRequest(ctx)
	if sessionID == "" {
		return nil, err
	}

	sessionRedisKey := sessionIDToKey(sessionID)
	sessionJSON, err := redisClient.Get(sessionRedisKey).Bytes()
	if err != nil {
		return nil, err
	}

	sess := LambdaSession{}
	sess.UnmarshalJSON(sessionJSON)
	return &sess, nil
}

// SaveToContext saves the current session onto the request (in a cookie).
func (session *LambdaSession) SaveToContext(context *fasthttp.RequestCtx) error {
	var err error
	if len(session.sessionID) == 0 {
		// The session didn't already exist, so we need to create a new ID and cookie.
		session.sessionID, err = generateSessionID()
		if err != nil {
			return err
		}

		// It's astronomically unlikely, but check for a conflict.
		exists, err := redisClient.Exists(sessionIDToKey(session.sessionID)).Result()
		if err != nil {
			return err
		}
		if exists > 0 {
			return errors.New("Conflict when generated session ID")
		}

		cookie := fasthttp.AcquireCookie()
		defer fasthttp.ReleaseCookie(cookie)
		cookie.SetSecure(true)
		cookie.SetHTTPOnly(true)
		cookie.SetMaxAge(sessionDuration)
		cookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
		// cookie.SetPath("/api")  TODO add this restriction after figuring out nojs
		cookie.SetPath("/")
		cookie.SetKey("__Secure-lambdasession2")
		cookie.SetValue(session.sessionID)
		context.Response.Header.SetCookie(cookie)
	}

	// Persist to redis.
	redisKey := sessionIDToKey(session.sessionID)
	sessionJSON, err := session.MarshalJSON()
	if err != nil {
		return err
	}
	err = redisClient.Set(redisKey, sessionJSON, time.Duration(sessionDuration)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

// RemoveSession unsets the session cookie and deletes it from the session store.
func RemoveSession(ctx *fasthttp.RequestCtx) error {
	sessionID, err := getSessionIDFromRequest(ctx)
	if sessionID == "" {
		// note: err will be nil if the session cookie wasn't set
		return err
	}

	err = redisClient.Del(sessionID).Err()
	if err != nil {
		return err
	}

	// Only try to delete the cookie _after_ redis delete was successful.
	// That way a client can continue to attempt to log out in the event
	// of a redis error.
	cookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(cookie)
	cookie.SetSecure(true)
	cookie.SetHTTPOnly(true)
	cookie.SetMaxAge(sessionDuration)
	cookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	// cookie.SetPath("/api")  TODO add this restriction after figuring out nojs
	cookie.SetPath("/")
	cookie.SetKey("__Secure-lambdasession2")
	cookie.SetValue("")
	cookie.SetExpire(fasthttp.CookieExpireDelete)
	ctx.Response.Header.SetCookie(cookie)

	return nil
}
