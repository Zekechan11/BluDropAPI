package api

import (
	"fmt"
	"net/http"
	"bludrop-api/dto"
	"bludrop-api/util"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AuthRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/users", func(c *gin.Context) {
		var users []struct {
			ID   int    `db:"id"`
			Name string `db:"firstname"`
			Area string `db:"area"`
		}
		err := db.Select(&users, "SELECT id, firstname, area FROM accounts")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	r.GET("/users/count", func(c *gin.Context) {
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM accounts WHERE role = 'customer'")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"total_users": count})
	})

	r.POST("/accounts", func(c *gin.Context) {
		var account struct {
			Firstname string `json:"firstname" binding:"required"`
			Lastname  string `json:"lastname" binding:"required"`
			Email     string `json:"email" binding:"required"`
			Area      string `json:"area" binding:"required"`
			Password  string `json:"password" binding:"required"`
			Username  string `json:"username" binding:"required"`
			Role      string `json:"role" binding:"required,oneof=customer staff"`
		}

		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO accounts (firstname, lastname, email, area, password, qrcode, role) 
				  VALUES (:firstname, :lastname, :email, :area, :password, :username, :role)`
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
				"username":  account.Username,
				"role":      account.Role,
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
			ID        int    `db:"id"`
			FirstName string `db:"firstName"`
			LastName  string `db:"lastName"`
			Area      string `db:"area"`
			Email     string `db:"email"`
			Role      string `db:"role"`
			Password  string `db:"password"`
		}

		err := db.Get(&user, "SELECT id, firstName, lastName, area, email, password, role FROM accounts WHERE email = ?", account.Email)
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
			"message":   "Login successful",
			"token":     token,
			"id":        user.ID, // Corrected to use user.ID
			"role":      user.Role,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"area":      user.Area,
			"email":     user.Email,
		})
	})



	r.POST("/v2/api/login", func(c *gin.Context) {
		var account struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user dto.LoginData

		query := `
			SELECT staff_id AS uid, area_id, firstname, lastname, NULL AS username, email, password, role, area, NULL AS type
			FROM account_staffs
			LEFT JOIN areas ON id = area_id
			WHERE email = ?
			UNION ALL
			SELECT client_id AS uid, area_id, firstname, lastname, username, email, password, role, area, type
			FROM account_clients
			LEFT JOIN areas ON id = area_id
			WHERE email = ?`

		err := db.Get(&user, query, account.Email, account.Email)
		if err != nil {
			fmt.Println("Error fetching user:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if account.Password != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := util.GenerateJWT(user.Uid, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Login successful",
			"token":     token,
			"user_info": gin.H{
				"uid":   	 user.Uid,
				"role":      user.Role,
				"firstname": user.FirstName,
				"lastname":  user.LastName,
				"email":     user.Email,
				"username":	 user.UserName,
				"area": 	 user.Area,
				"area_id":	 user.AreaId,
				"type":		 user.Type,
			},
		})
	})

	r.POST("/api/auth/forgot-password", func(c *gin.Context) {
		var account struct {
			Email    string `json:"email" binding:"required"`
		}

		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request",
			})
			return
		}

		query := `
			SELECT staff_id AS uid, email
			FROM account_staffs
			LEFT JOIN areas ON id = area_id
			WHERE email = ?
			UNION ALL
			SELECT client_id AS uid, email
			FROM account_clients
			LEFT JOIN areas ON id = area_id
			WHERE email = ?`

		var user struct {
			UID   int    `db:"uid"`
			Email string `db:"email"`
		}
		err := db.Get(&user, query, account.Email, account.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid request",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message":   "Email found",
		})
	})
}
