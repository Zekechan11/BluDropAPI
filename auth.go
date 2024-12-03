package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(r *gin.Engine, db *sql.DB) {
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
		var account struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		// Bind JSON to account struct
		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Check if the email and password are correct, and also select the role
		var user struct {
			ID       int    `json:"id"`
			Email    string `json:"email"`
			Role     string `json:"role"`
			Password     string `json:"password"`
		}
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

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
			"role":    user.Role,
		})
	})
}