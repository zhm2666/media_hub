package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtHashKey string

func SetJwtHashKey(key string) {
	jwtHashKey = key
}

type UserClaims struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	jwt.RegisteredClaims
}

func GenerateTokenWithUserCenter(userID int, name, avatar string) (string, error) {
	hs := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"avatar":  avatar,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})
	return hs.SignedString([]byte(jwtHashKey))
}

func VerifyToken(tokenString string) (*UserClaims, error) {
	if jwtHashKey == "" {
		return nil, errors.New("jwt hash key not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtHashKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		name := claims["name"].(string)
		avatar := ""
		if a, ok := claims["avatar"].(string); ok {
			avatar = a
		}
		return &UserClaims{
			UserID: userID,
			Name:   name,
			Avatar: avatar,
		}, nil
	}

	return nil, errors.New("invalid token")
}

func HttpGet(client *http.Client, url string) ([]byte, error) {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024)
	var result []byte
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
		}
		if err != nil {
			break
		}
	}
	return result, nil
}
