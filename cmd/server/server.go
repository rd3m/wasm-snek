package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()
var client *redis.Client

type scoreEntry struct {
	Name  string
	Score int
}

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
	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"showCurrentScore": true,
		})
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
		names, err := client.Keys(ctx, "*").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var scores []scoreEntry
		for _, name := range names {
			score, err := client.Get(ctx, name).Int()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			scores = append(scores, scoreEntry{Name: name, Score: score})
		}

		sort.Slice(scores, func(i, j int) bool {
			return scores[i].Score > scores[j].Score
		})

		c.HTML(http.StatusOK, "scores.html", gin.H{
			"scores": scores,
		})
	})

	r.Run()
}
