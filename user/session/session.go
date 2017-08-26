package session

import (
	"github.com/kataras/go-sessions"
	"github.com/kataras/go-sessions/sessiondb/redis"
	"github.com/kataras/go-sessions/sessiondb/redis/service"
)

var lambdaSessionsConfig = sessions.Config{
	Cookie: "lambdasession",
}

// Sessions is the go-sessions instance of Lambda
var Sessions = sessions.New(lambdaSessionsConfig)

func init() {
	sessionDb := redis.New(service.Config{
		Network:       service.DefaultRedisNetwork,
		Addr:          service.DefaultRedisAddr,
		Password:      "",
		Database:      "",
		MaxIdle:       0,
		MaxActive:     0,
		IdleTimeout:   service.DefaultRedisIdleTimeout,
		Prefix:        "lmdasessions",
	})

	Sessions.UseDatabase(sessionDb)
}
