package api

import (
	"net/http"

	_ "github.com/Dostonlv/hackathon-nt/docs"
	"github.com/Dostonlv/hackathon-nt/internal/api/handlers"
	"github.com/Dostonlv/hackathon-nt/internal/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/k0kubun/pp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Claims represents the JWT claims structure
type Claims struct {
	RolePayload string `json:"role"`
	jwt.StandardClaims
}

// AuthorizationMiddleware checks permissions using Casbin with JWT role
func AuthorizationMiddleware(enforcer *casbin.Enforcer, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		pp.Println(authHeader)

		// Extract the token from the "Bearer <token>" format
		// tokenParts := strings.Split(authHeader, " ")
		// if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
		// 	c.Abort()
		// 	return
		// }

		// tokenString := tokenParts[1]

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Extract role from claims
		userRole := claims.RolePayload
		if userRole == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found in token"})
			c.Abort()
			return
		}

		// Get request path and method
		path := c.Request.URL.Path
		method := c.Request.Method

		// Check permission using Casbin
		allowed, err := enforcer.Enforce(userRole, path, method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization error"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		// Store role in context for later use if needed
		c.Set("userRole", userRole)
		c.Next()
	}
}

// NewRouter -.
// Swagger spec:
// @title       hackathon
// @description Backend
// @version     1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func SetupRouter(authService *service.AuthService, tenderService *service.TenderService, enforcer *casbin.Enforcer, jwtSecret string) *gin.Engine {
	router := gin.Default()

	authHandler := handlers.NewAuthHandler(authService)
	tenderHandler := handlers.NewTenderHandler(tenderService)
	// Public routes (no authorization required)
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Protected routes (require authorization)
	api := router.Group("/api")
	api.Use(AuthorizationMiddleware(enforcer, jwtSecret))
	{
		api.POST("/tenders", tenderHandler.CreateTender)
		// Add your protected routes here
		// Example:
		// api.GET("/users", userHandler.GetUsers)
		// api.POST("/users", userHandler.CreateUser)
	}

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
