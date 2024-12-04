package main

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"waterfalls/api"
)

func main() {
	dsn := "root@tcp(localhost:3306)/waterfalls"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Database connected successfully")

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go api.HandleMessages(db)

	api.AuthRoutes(r, db)
	api.RegisterRoutes(r, db)
	api.RegisterAgentRoutes(r, db)
	api.ChatRoutes(r, db)
	
	r.Run(":9090")
}
