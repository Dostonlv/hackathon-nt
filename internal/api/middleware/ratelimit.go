package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)


type BidRateLimiter struct {
	mu       sync.RWMutex
	requests map[string]*bidRequests
	limit    int                     
	window   time.Duration          
}

type Claims struct {
	RolePayload string `json:"role"`
	UserID      string `json:"user_id"`
	jwt.StandardClaims
}

type bidRequests struct {
	count     int       
	startTime time.Time 
}


func NewBidRateLimiter() *BidRateLimiter {
	limiter := &BidRateLimiter{
		requests: make(map[string]*bidRequests),
		limit:    5,           
		window:   time.Minute,
	}


	go limiter.cleanupLoop()

	return limiter
}


func (rl *BidRateLimiter) BidRateLimitMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

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


		if userRole != "contractor" {

			c.Next()
			return
		}

		userId := claims.UserID
		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
			c.Abort()
			return
		}

		if !rl.allowBidSubmission(userId) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded. Please try again later.",
				"retry_after": time.Now().Add(rl.window).Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}


func (rl *BidRateLimiter) allowBidSubmission(userID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	req, exists := rl.requests[userID]
	if !exists {
		rl.requests[userID] = &bidRequests{
			count:     1,
			startTime: now,
		}
		return true
	}


	if now.Sub(req.startTime) > rl.window {
		req.count = 1
		req.startTime = now
		return true
	}


	if req.count < rl.limit {
		req.count++
		return true
	}
	return false
}


func (rl *BidRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for userID, req := range rl.requests {
			if now.Sub(req.startTime) > rl.window {
				delete(rl.requests, userID)
			}
		}
		rl.mu.Unlock()
	}
}
