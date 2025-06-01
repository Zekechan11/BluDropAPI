package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"waterfalls/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Database connected successfully")

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go api.HandleMessages(db)

	api.AuthRoutes(r, db)
	api.RegisterRoutes(r, db)
	api.AdminRoutes(r, db)
	api.AgentRoutes(r, db)
	api.ChatRoutes(r, db)
	api.InventoryRoutes(r, db)
	api.Customer_OrderRoutes(r, db)
	api.StaffRoutes(r, db)
	api.ClientRoutes(r, db)
	api.PaymentRoutes(r, db)
	api.CustomerRoutes(r, db)
	api.TransactionRoutes(r, db)
	api.ScheduleRoutes(r, db)
	api.FGSRoutes(r, db)
	api.PricingRoutes(r, db)
	api.SalesReportRoutes(r, db)
	api.ManualOrderRoutes(r, db)
	api.RegisterRemittanceRoutes(r, db)
	api.ManageProfileRoutes(r, db)
	api.NotificationRoutes(r, db)

	r.Run(":9090")
}
