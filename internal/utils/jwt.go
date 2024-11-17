package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/Dostonlv/hackathon-nt/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTUtil struct {
	secretKey []byte
}

type Claims struct {
	UserID uuid.UUID       `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.StandardClaims
}

func NewJWTUtil(secretKey string) *JWTUtil {
	return &JWTUtil{
		secretKey: []byte(secretKey),
	}
}

func (j *JWTUtil) GenerateToken(userID uuid.UUID, role models.UserRole) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
func (j *JWTUtil) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ParseToken(tokenString string, secretKey []byte) (*Claims, error) {
	var jwtString string
	if strings.HasPrefix(tokenString, "Bearer ") {
		jwtString = strings.Split(tokenString, "Bearer ")[1]
	} else {
		jwtString = tokenString
	}

	token, err := jwt.ParseWithClaims(jwtString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
