package session

import (
	"github.com/kataras/go-sessions"
	"github.com/kataras/go-sessions/sessiondb/redis"
	"github.com/kataras/go-sessions/sessiondb/redis/service"

	"github.com/mstojcevich/lambda-ng-go/config"
)

// TODO sign the cookie. It's random, so it should be fine, but signing sure wouldn't hurt.

var lambdaSessionsConfig = sessions.Config{
	Cookie: "lambdasession",
	CookieSecureTLS: true,
}

// Sessions is the go-sessions instance of Lambda
var Sessions = sessions.New(lambdaSessionsConfig)

// TODO FIXME: This session library doesn't support Secure cookies on proxied connections!

func init() {
	sessionDb := redis.New(service.Config{
		Network:       "tcp",
		Addr:          config.RedisAddr,
		Password:      config.RedisPassword,
		Database:      "",
		MaxIdle:       0,
		MaxActive:     0,
		IdleTimeout:   service.DefaultRedisIdleTimeout,
		Prefix:        "lmdasessions",
	})

	Sessions.UseDatabase(sessionDb)
}
