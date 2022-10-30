package main

import (
	"net/http"
	id_generator "utils/id_generator"
	rate_limiter "utils/rate_limiter_leakybucket"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func main() {
	gin.SetMode("debug")

	router.Use(rate_limiter.RateLimit(10))

	router.GET("/id", func(c *gin.Context) {
		id, _ := id_generator.NextID()
		c.JSON(http.StatusOK, gin.H{
			"message": id,
		})
	})
	router.Run()
}
