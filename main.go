package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()

	// Create a gin route with a parameter called `id`
	r.GET("/hello/:id", func(c *gin.Context) {
		// Respond with a JSON, serialized from a map[string]interface{}
		// gin.H is a shortcut for map[string]interface{}
		c.JSON(200, gin.H{
			"Test": "Hello " + c.Param("id"), // Obtain value for parameter `:id`
		})
	})
	// Listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// Or the port you specified in the environment variable `PORT`
	r.Run()
}
