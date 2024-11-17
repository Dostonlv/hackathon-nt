package api

import (
	"net/http"
	"strings"

	_ "github.com/Dostonlv/hackathon-nt/docs"
	"github.com/Dostonlv/hackathon-nt/internal/api/handlers"
	"github.com/Dostonlv/hackathon-nt/internal/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Claims represents the JWT claims structure
type Claims struct {
	RolePayload string    `json:"role"`
	UserID      uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

// AuthorizationMiddleware checks permissions using Casbin with JWT role
func AuthorizationMiddleware(enforcer *casbin.Enforcer, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
			c.Abort()
			return
		}
		var jwtString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			jwtString = strings.Split(authHeader, "Bearer ")[1]
		} else {
			jwtString = authHeader
		}

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
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

		userId := claims.UserID.String()
		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
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
		c.Set("userId", userId)
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
func SetupRouter(authService *service.AuthService, tenderService *service.TenderService, bidService *service.BidService, enforcer *casbin.Enforcer, jwtSecret string) *gin.Engine {
	router := gin.Default()

	authHandler := handlers.NewAuthHandler(authService)
	tenderHandler := handlers.NewTenderHandler(tenderService)

	bidHandler := handlers.NewBidHandler(bidService)
	// Public routes (no authorization required)
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Protected routes (require authorization)
	api := router.Group("/api")
	api.Use(AuthorizationMiddleware(enforcer, jwtSecret))
	{
		api.POST("/client/tenders", tenderHandler.CreateTender)
		api.GET("/client/tenders", tenderHandler.ListTenders)
		api.PUT("/client/tenders/:id", tenderHandler.UpdateTenderStatus)
		api.DELETE("/client/tenders/:id", tenderHandler.DeleteTender)
		api.GET("/client/tenders/:tender_id/bids", bidHandler.GetBidsByClientID)
		api.POST("/client/tenders/:tender_id/award/:bid_id", bidHandler.AwardBid)
		api.GET("/client/tenders/filter", tenderHandler.ListTendersFiltering)

		api.POST("/contractor/tenders/:tender_id/bid", bidHandler.CreateBid)
		api.GET("/contractor/bids", bidHandler.GetBidsByContractorID)
		api.DELETE("/contractor/bids/:bid_id", bidHandler.DeleteBidByContractorID)
	}

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
