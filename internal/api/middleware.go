package api

//import (
//	"github.com/dgrijalva/jwt-go"
//	"github.com/gin-gonic/gin"
//	"go-canteen-devpart/auth"
//	"net/http"
//	"strconv"
//	"strings"
//)
//
//var jwtKey = []byte("your_secret_key")
//
//func TokenAuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		tokenString := c.GetHeader("Authorization")
//		tokenString = strings.TrimPrefix(tokenString, "Bearer ") // Удаление префикса "Bearer "
//		if tokenString == "" {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
//			return
//		}
//		claims := &jwt.StandardClaims{}
//		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//			return jwtKey, nil
//		})
//
//		if err != nil || !token.Valid || claims.Audience != "admin" {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid tokensss"})
//			return
//		}
//
//		userID, err := strconv.Atoi(claims.Subject)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
//			return
//		}
//		c.Set("userID", userID)
//
//		c.Next()
//	}
//}
//
////func GenerateToken(userID uint, isAdmin bool) (string, error) {
////	expirationTime := time.Now().Add(24 * time.Hour) // Токен истекает через 24 часа
////	claims := &jwt.StandardClaims{
////		Subject:   fmt.Sprintf("%d", userID),
////		ExpiresAt: expirationTime.Unix(),
////		Issuer:    "your-app-name",
////		// Вы можете добавить другие поля в claims, например, роль пользователя
////		Audience: "admin", // Указание, что токен предназначен для администраторов
////	}
////
////	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
////	tokenString, err := token.SignedString(jwtKey)
////	log.Printf("Generated Token: %s", tokenString)
////	return tokenString, err
////}
//
//// В вашем Go файле сервера
//func AdminDashboardHandler(c *gin.Context) {
//	c.HTML(http.StatusOK, "adminPage.html", nil)
//}
//
//func AuthAdminMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		authHeader := r.Header.Get("Authorization")
//		if authHeader == "" {
//			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
//			return
//		}
//
//		token := strings.TrimPrefix(authHeader, "Bearer ")
//		claims, ok := auth.VerifyToken(token)
//		if !ok {
//			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
//			return
//		}
//
//		if !claims.Admin {
//			http.Error(w, "Access restricted to administrators", http.StatusForbidden)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}
