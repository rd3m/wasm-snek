package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()
var client *redis.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	client = redis.NewClient(opt)
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/wasm", "../wasm")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/saveScore", func(c *gin.Context) {
		name := c.PostForm("name")
		score := c.PostForm("score")

		err := client.Set(ctx, name, score, 0).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Score saved"})
	})

	r.GET("/scores", func(c *gin.Context) {
		keys, err := client.Keys(ctx, "*").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		scores := make(map[string]int)
		for _, name := range keys {
			score, err := client.Get(ctx, name).Int()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			scores[name] = score
		}

		c.HTML(http.StatusOK, "scores.html", gin.H{
			"scores": scores,
		})
	})

	r.Run()
}
