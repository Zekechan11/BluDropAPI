package api

import (
	"net/http"
	"waterfalls/util"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AuthRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/users", func(c *gin.Context) {
		var users []struct {
			ID   int    `db:"id"`
			Name string `db:"firstname"`
		}
		err := db.Select(&users, "SELECT id, firstname FROM accounts")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
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

		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO accounts (firstname, lastname, email, area, password, qrcode) 
				  VALUES (:firstname, :lastname, :email, :area, :password, :username)`
		result, err := db.NamedExec(query, account)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Account created successfully",
			"account": gin.H{
				"id":        id,
				"firstname": account.Firstname,
				"lastname":  account.Lastname,
				"email":     account.Email,
			},
		})
	})

	r.POST("/login", func(c *gin.Context) {
		var account struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user struct {
			ID       int    `db:"id"`
			Email    string `db:"email"`
			Role     string `db:"role"`
			Password string `db:"password"`
		}

		err := db.Get(&user, "SELECT id, email, password, role FROM accounts WHERE email = ?", account.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if account.Password != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := util.GenerateJWT(user.ID, user.Email, user.Role)
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
