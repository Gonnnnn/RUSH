package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"rush/server"
)

type SignInRequest struct {
	Token string `json:"token"`
}

func handleSignIn(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := server.SignIn(req.Token)
		if err != nil {
			code := getHttpStatus(err)
			if code == http.StatusBadRequest {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
				return
			}
			if code == http.StatusNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			if code == http.StatusInternalServerError {
				log.Printf("Error signing in: %+v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func handleAuth(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("userId")
		if userId == "" {
			log.Printf("Error getting user ID from context, it is supposed to be set by the middleware")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userId})
	}
}

// TODO(#86): Fix it so that only certain users can see the list of users.
func handleListUsers(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
			return
		}
		pageSize, err := strconv.Atoi(c.Query("pageSize"))
		if err != nil || pageSize < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize"})
			return
		}

		result, err := server.ListUsers(offset, pageSize)
		if err != nil {
			log.Printf("Error getting users: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"users":       result.Users,
			"is_end":      result.IsEnd,
			"total_count": result.TotalCount,
		})
	}
}

// TODO(#86): Fix it so that only the admin or the user itself can see the user's information.
func handleGetUser(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		user, err := server.GetUser(id)
		if err != nil {
			if isNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}

			log.Printf("Error getting user: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

type addUserRequest struct {
	Name       string  `json:"name"`
	Generation float64 `json:"generation"`
	IsActive   bool    `json:"is_active"`
	Email      string  `json:"email"`
}

func handleAddUser(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req addUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := server.AddUser(
			req.Name,
			req.Generation,
			req.IsActive,
			req.Email,
		); err != nil {
			log.Printf("Error adding user: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User added successfully"})
	}
}

func handleListSessions(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
			return
		}
		pageSize, err := strconv.Atoi(c.Query("pageSize"))
		if err != nil || pageSize < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize"})
			return
		}

		result, err := server.ListSessions(offset, pageSize)
		if err != nil {
			log.Printf("Error getting sessions: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"sessions":    result.Sessions,
			"is_end":      result.IsEnd,
			"total_count": result.TotalCount,
		})
	}
}

func handleGetSession(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		session, err := server.GetSession(id)
		if err != nil {
			if isNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
				return
			}

			log.Printf("Error getting session: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, session)
	}
}

type addSessionRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	Score       int       `json:"score"`
}

func handleAddSession(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req addSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId := c.GetString("userId")
		if userId == "" {
			log.Printf("Error getting user ID from context, it is supposed to be set by the middleware")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		id, err := server.AddSession(req.Name, req.Description, userId, req.StartsAt, req.Score)
		if err != nil {
			log.Printf("Error adding session: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

func handleCreateAttendanceForm(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.Param("id")
		formUrl, err := server.CreateAttendanceForm(sessionId)
		if err != nil {
			code := getHttpStatus(err)
			if code == http.StatusNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
				return
			}

			if code == http.StatusBadRequest {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Form already exists"})
				return
			}

			log.Printf("Error creating session form: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"form_url": formUrl})
	}
}

func handleApplyAttendance(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.Param("id")
		if err := server.CloseSession(sessionId); err != nil {
			if isBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Session already closed"})
				return
			}

			if isNotFound(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
				return
			}

			log.Printf("Error closing session: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Session closed successfully"})
	}
}

func handleGetAttendanceForUser(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
			return
		}

		attendances, err := server.GetAttendanceByUserId(userId)
		if err != nil {
			log.Printf("Error getting attendance for user: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"attendances": attendances})
	}
}

func handleHalfYearAttendance(server *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := server.GetHalfYearAttendance()
		if err != nil {
			log.Printf("Error applying half year attendance: %+v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sessions": result.Sessions, "users": result.Users, "attendances": result.Attendances})
	}
}
