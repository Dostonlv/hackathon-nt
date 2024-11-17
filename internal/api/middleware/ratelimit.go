package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// BidRateLimiter contractor user larning bid submission so'rovlarini cheklash uchun
type BidRateLimiter struct {
	mu       sync.RWMutex
	requests map[string]*bidRequests // user ID ga bog'langan so'rovlar
	limit    int                     // Maksimal so'rovlar soni (5)
	window   time.Duration           // Vaqt oralig'i (1 minut)
}

type Claims struct {
	RolePayload string `json:"role"`
	UserID      string `json:"user_id"`
	jwt.StandardClaims
}

type bidRequests struct {
	count     int       // So'rovlar soni
	startTime time.Time // Birinchi so'rov vaqti
}

// NewBidRateLimiter creates a new rate limiter for bid submissions
func NewBidRateLimiter() *BidRateLimiter {
	limiter := &BidRateLimiter{
		requests: make(map[string]*bidRequests),
		limit:    5,           // minutiga 5 ta so'rov
		window:   time.Minute, // 1 minutlik oyna
	}

	// Eski yozuvlarni tozalash uchun goroutine ishga tushirish
	go limiter.cleanupLoop()

	return limiter
}

// BidRateLimitMiddleware bid submission endpointi uchun Gin middleware
func (rl *BidRateLimiter) BidRateLimitMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Faqat POST so'rovlarni cheklash
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

		// Check if user role is "contractor"
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

// allowBidSubmission checks if the contractor can submit another bid
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

	// Vaqt oynasi tugagan bo'lsa, qayta boshlash
	if now.Sub(req.startTime) > rl.window {
		req.count = 1
		req.startTime = now
		return true
	}

	// Limitni tekshirish
	if req.count < rl.limit {
		req.count++
		return true
	}
	return false
}

// cleanupLoop removes expired entries periodically
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
