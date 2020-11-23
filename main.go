package main

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()

	// Use built-in Static handler to serve static files efficiently
	r.Static("/assets", "./assets")

	// Load HTML templates inside templates folder
	r.LoadHTMLGlob("templates/*")

	// Render index.html.tmpl as index page
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html.tmpl", struct{}{})
	})

	// Return current time (formatted) for AJAX requests
	r.GET("/now", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"time": time.Now().Local().Format("2006/01/02 15:04:05"),
		})
	})

	// Response to <form> submission
	r.POST("/", func(c *gin.Context) {
		// Define a struct to receive form data
		params := struct {
			Username string `form:"username"`
			Password string `form:"password"`
		}{}
		// Parse POST body (JSON/XML/x-www-form-urlencode) and save to `params`
		c.Bind(&params)
		// Hash user-provided password
		hashedPwd := md5.Sum([]byte(params.Password))
		params.Password = hex.EncodeToString(hashedPwd[:])
		// Render result.html.tmpl with params
		c.HTML(200, "result.html.tmpl", params)
	})

	// Listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// Or the port you specified in the environment variable `PORT`
	r.Run()
}
