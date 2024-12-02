package main

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	// Configure MySQL connection
	dsn := "root@tcp(localhost:3306)/waterfalls"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Database connected successfully")

	// Set up Gin router
	r := gin.Default()

	r.Use(cors.Default())

	// Example route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	RegisterRoutes(r, db)
	RegisterAgentRoutes(r, db)

	// Route to test MySQL query
	r.GET("/users", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, firstname FROM accounts")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var users []map[string]interface{}
		for rows.Next() {
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, gin.H{"id": id, "name": name})
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	r.POST("/accounts", func(c *gin.Context) {
		var account struct {
			Firstname string `json:"firstname" binding:"required"`
			Lastname  string `json:"lastname" binding:"required"`
			Email     string `json:"email" binding:"required"`
			Area      string `json:"area" binding:"required"`
			Password  string `json:"password" binding:"required"`
			Username  string `json:"username" binding:"required"`
		}

		// Parse and validate JSON input
		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert new account into the database
		result, err := db.Exec("INSERT INTO accounts (firstname, lastname, email, area, password, qrcode ) VALUES (?, ?, ?, ?, ?, ?)",
			account.Firstname, account.Lastname, account.Email, account.Area, account.Password, account.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get the ID of the newly inserted account
		id, _ := result.LastInsertId()
		c.JSON(http.StatusCreated, gin.H{
			"message": "Account created successfully",
			"account": gin.H{"id": id, "firstname": account.Firstname, "lastname": account.Lastname, "email": account.Email},
		})
	})

	r.POST("/login", func(c *gin.Context) {
		var account Account
		// Bind JSON to account struct
		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Check if the email and password are correct, and also select the role
		var user User
		err := db.QueryRow("SELECT id, email, password, role FROM accounts WHERE email = ?", account.Email).
			Scan(&user.ID, &account.Email, &user.Password, &user.Role)
		if err != nil || account.Password != user.Password {
			// If email doesn't exist or password is incorrect
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT token
		token, err := generateJWT(user.ID, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Return the JWT token
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
			"role":    user.Role, // Include the role in the response
		})
	})

	// Start the server
	r.Run(":9090") // Default port is 8080
}
