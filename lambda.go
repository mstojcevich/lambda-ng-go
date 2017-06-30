package main

import (
	"log"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/mstojcevich/lambda-ng-go/config"
	"github.com/mstojcevich/lambda-ng-go/fileserve"
	"github.com/mstojcevich/lambda-ng-go/index"
	"github.com/mstojcevich/lambda-ng-go/upload"
	"github.com/mstojcevich/lambda-ng-go/user"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/", index.Page)
	router.GET("/nojs/", index.PageNoJS)

	// User
	router.GET("/login", user.LoginPage)
	router.POST("/api/user/login", user.LoginAPI)
	router.GET("/register", user.RegisterPage)
	router.POST("/api/user/new", user.RegisterAPI)
	router.PUT("/api/user/new", user.RegisterAPI)
	router.GET("/api/session", user.GetSessionAPI)
	router.GET("/user/manage", user.ManagePage)
	router.GET("/nojs/user/manage", user.ManagePageNoJS)
	router.DELETE("/api/session", user.LogoutAPI)
	router.GET("/nojs/user/logout", user.LogoutNoJS)

	// Past uploads
	router.GET("/user/uploads", upload.PastUploadsPage)
	router.GET("/nojs/user/uploads", upload.PastUploadsPageNoJS)
	router.GET("/api/user/uploads", upload.PastUploadsAPI)
	router.GET("/generic/by-ext/:extension", upload.GenericImageByExtension)

	// Upload
	router.GET("/upload", upload.Page)
	router.GET("/nojs/upload", upload.PageNoJS)
	router.POST("/api/upload", upload.API)
	router.PUT("/api/upload", upload.API)
	router.DELETE("/file/:name", upload.DeleteAPI)

	// Paste
	router.POST("/api/paste", upload.PutPasteAPI)
	router.PUT("/api/paste", upload.PutPasteAPI)
	router.GET("/api/paste", upload.GetPasteAPI)
	router.GET("/paste", upload.PastePage)

	// Favicon handler
	router.GET("/favicon.ico", func(ctx *fasthttp.RequestCtx) {
		ctx.SendFile("static/img/favicon.ico")
	})

	router.ServeFiles("/static/*filepath", "static")

	router.NotFound = fileserve.Serve

	router.PanicHandler = panicHandler

	s := &fasthttp.Server{
		Name:               "Lambda",
		Handler:            router.Handler,
		MaxRequestBodySize: 1024 * 1024 * config.MaxUploadSize,
		MaxConnsPerIP:      1024,
		ReadTimeout:        2 * time.Minute,
		WriteTimeout:       2 * time.Minute,
		MaxRequestsPerConn: 512,
	}

	log.Fatal(s.ListenAndServe(":8080"))
}

func panicHandler(ctx *fasthttp.RequestCtx, err interface{}) {
	ctx.Error("500", fasthttp.StatusInternalServerError)
	log.Println("Panic while handling request")
	log.Println(err)
}
