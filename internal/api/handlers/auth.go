package handlers

import (
	"net/http"
	"strings"

	"github.com/Dostonlv/hackathon-nt/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param input body service.RegisterInput true "Register Input"
// @Success 201 {object} service.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
		return
	}

	if input.Role != "client" && input.Role != "contractor" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid role"})
		return
	}

	// Custom validation logic
	if input.Email == "" || input.Username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "username or email cannot be empty"})
		return
	}

	// Basic email format validation
	if !strings.Contains(input.Email, "@") || !strings.Contains(input.Email, ".") {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid email format"})
		return
	}

	// Call AuthService to register
	resp, err := h.authService.Register(c.Request.Context(), input)
	if err != nil {
		// Handle specific errors
		switch err.Error() {
		case "email already exists":
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Email already exists"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to register user"})
		}
		return
	}

	// Success response
	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary Login a user
// @Description Login a user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param input body service.LoginInput true "Login Input"
// @Success 200 {object} service.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
		return
	}

	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Username and password are required"})
		return
	}

	// Call AuthService to login
	resp, err := h.authService.Login(c.Request.Context(), input)
	if err != nil {
		// Handle specific errors
		switch err.Error() {
		case "invalid credentials":
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid username or password"})
		case "user not found":
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to login user"})
		}
		return
	}

	// Success response
	c.JSON(http.StatusOK, resp)
}
