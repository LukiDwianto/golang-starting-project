package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"golang.org/x/crypto/bcrypt"
)



func GetPrivileges(data interface{}) ([]string, error) {
	privilegesSlice, ok := data.([]interface{})
	if !ok {
		return []string{}, errors.New("failed to assret")
	}

	privileges := make([]string, len(privilegesSlice))
	for i, v := range privilegesSlice {
		privilege, ok := v.(string)
		if !ok {
			return []string{}, errors.New("failed to assret")
		}
		privileges[i] = privilege
	}

	return privileges, nil

}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateToken(username string, role string, privileges []string, c *gin.Context) (string, string, error) {

	claims := jwt.MapClaims{
		"username":   username,
		"role":       role,
		"privileges": privileges,
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := os.Getenv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"username":   username,
		"role":       role,
		"privileges": privileges,
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
		"iat":        time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(key))
	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := extractTokenFromHeader(c.GetHeader("Authorization"))
		if accessToken == "" {
			unauthorizedResponse(c, "Unauthorized")
			return
		}

		token, err := parseToken(accessToken)
		if err != nil || !token.Valid {
			refreshToken := extractTokenFromHeader(c.GetHeader("refresh_token"))
			if refreshToken == "" {
				unauthorizedResponse(c, "Unauthorized")
				return
			}

			if err := handleRefreshToken(refreshToken, c); err != nil {
				unauthorizedResponse(c, "Invalid tokens")
				return
			}
		} else {
			if err := setTokenClaims(token, c); err != nil {
				unauthorizedResponse(c, "Invalid token claims")
				return
			}
		}
	}
}

func extractTokenFromHeader(header string) string {
	return strings.TrimPrefix(header, "Bearer ")
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
}

func handleRefreshToken(refreshToken string, c *gin.Context) error {
	refreshClaims, err := parseToken(refreshToken)
	if err != nil || !refreshClaims.Valid {
		return err
	}

	claims, ok := refreshClaims.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("Invalid token claims")
	}

	return generateAndSetNewAccessToken(claims, c)
}

func setTokenClaims(token *jwt.Token, c *gin.Context) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("Invalid token claims")
	}
	return generateAndSetNewAccessToken(claims, c)
}

func generateAndSetNewAccessToken(claims jwt.MapClaims, c *gin.Context) error {
	username, role, privileges, err := extractClaims(claims)
	if err != nil {
		return err
	}

	privilegesData, err := GetPrivileges(privileges)
	if err != nil {
		return fmt.Errorf("error get privileges: %v", err)
	}

	newAccessToken, _, err := GenerateToken(username, role, privilegesData, c)
	if err != nil {
		return fmt.Errorf("failed to generate new access token")
	}

	c.Header("access_token", newAccessToken)
	c.Set("username", username)
	c.Set("role", role)
	c.Set("privileges", privilegesData)
	c.Set("claims", claims)
	return nil
}

func extractClaims(claims jwt.MapClaims) (string, string, interface{}, error) {
	username, ok := claims["username"].(string)
	if !ok {
		return "", "", nil, fmt.Errorf("missing or invalid username")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", nil, fmt.Errorf("missing or invalid role")
	}

	privileges, ok := claims["privileges"]
	if !ok {
		return "", "", nil, fmt.Errorf("missing or invalid privileges")
	}

	return username, role, privileges, nil
}

func unauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
	c.Abort()
}


func CheckPrivilegesMiddleware(privilege []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userprivilege := c.GetStringSlice("privileges")

		hasRequiredprivilege := false

		if privilegesContains(userprivilege, privilege) {
			hasRequiredprivilege = true
		}
		if !hasRequiredprivilege {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
	}
}

func privilegesContains(privileges []string, privilegesUser []string) bool {
	for _, v := range privileges {

		for _, value := range privilegesUser {
			if v == value {
				return true
			}

		}

	}
	return false
}