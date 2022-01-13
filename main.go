package main

import (
	"log"
	"net/http"

	handlers "./handlers"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
)

const CONN_PORT    = ":9090"

func main() {
	
	router := gin.Default()
		var httpServer = http.Server{
		Addr: CONN_PORT,
		Handler: router,
	}
	var http2Server = http2.Server{}//http2 server
	_ = http2.ConfigureServer(&httpServer, &http2Server)

	router.LoadHTMLGlob("static/template/*")
	router.Static("/static/js", "./static/js")
	router.Static("/static/css", "./static/css")

	router.GET("/", handlers.GetHome)
	router.POST("/login", handlers.PostHome)
	router.GET("/register", handlers.GetRegister)
	router.POST("/register", handlers.PostRegister)
	router.POST("/quizz", handlers.SwitchQuestion)

	err := router.RunTLS(CONN_PORT, "./server.crt", "./server.key")
	if err != nil {
		log.Fatal("Error run TLS")
	}
}
