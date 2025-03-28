package middleware

import (
	"context"
	"net/http"
	"strings"

	"chat-llama/internal/api/handlers"

	"github.com/dgrijalva/jwt-go"
)

// JWTMiddleware JWT认证中间件
type JWTMiddleware struct {
	jwtSecret []byte
}

// NewJWTMiddleware 创建JWT中间件
func NewJWTMiddleware(jwtSecret string) *JWTMiddleware {
	return &JWTMiddleware{
		jwtSecret: []byte(jwtSecret),
	}
}

// Middleware 中间件处理函数
func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头或URL参数获取令牌
		tokenString := extractToken(r)
		if tokenString == "" {
			handlers.ErrorResponse(w, http.StatusUnauthorized, "需要认证")
			return
		}

		// 解析JWT令牌
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return m.jwtSecret, nil
		})

		// 处理令牌错误
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				handlers.ErrorResponse(w, http.StatusUnauthorized, "无效的令牌签名")
				return
			}
			handlers.ErrorResponse(w, http.StatusUnauthorized, "无效或过期的令牌")
			return
		}

		// 验证令牌
		if !token.Valid {
			handlers.ErrorResponse(w, http.StatusUnauthorized, "无效的令牌")
			return
		}

		// 将用户ID添加到请求上下文
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken 从请求中提取token
func extractToken(r *http.Request) string {
	// 从Authorization头部获取
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	// 从URL参数获取
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	// 从Cookie获取
	cookie, err := r.Cookie("token")
	if err == nil {
		return cookie.Value
	}

	return ""
}
