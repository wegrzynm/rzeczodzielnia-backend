package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type TokenClaims struct {
	Id    uint
	Email string
	Role  uint
	Exp   time.Time
}

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func CreateToken(id uint, username string, role uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    id,
			"email": username,
			"role":  role,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (error, *TokenClaims) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err, nil
	}

	if !token.Valid {
		return fmt.Errorf("invalid token"), nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims"), nil
	}
	tokenClaims := TokenClaims{
		Id:    uint(claims["id"].(float64)),
		Email: claims["email"].(string),
		Role:  uint(claims["role"].(float64)),
		Exp:   time.Unix(int64(claims["exp"].(float64)), 0),
	}
	return nil, &tokenClaims
}

func EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
