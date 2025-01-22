package api

import (
	"app/internal/domain"
	"app/internal/gateway"
	"app/internal/infrastructure/httputil"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) AuthRouter() {
	s.group.POST("/login", s.Login)
	s.group.POST("/register", s.Register)
	s.group.POST("/refresh", s.Refresh)
	s.group.POST("/logout", s.Logout)
}

func (s *server) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			httputil.NewResponse(c, http.StatusUnauthorized, "error", "missing or invalid token", nil)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validateJWT(token)
		if err != nil {
			s.logger.Errorw("Invalid access token", "err", err)
			httputil.NewResponse(c, http.StatusUnauthorized, "error", "invalid access token", nil)
			return
		}

		c.Set("claims", claims)

		c.Next()
	}
}

func (s *server) RoleMiddleware(requiredRole domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			httputil.NewResponse(c, http.StatusUnauthorized, "error", "missing or invalid token", nil)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validateJWT(token)
		if err != nil {
			s.logger.Errorw("Invalid token", "err", err)
			httputil.NewResponse(c, http.StatusUnauthorized, "error", "invalid token", nil)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			s.logger.Errorw("Role not found in token claims")
			httputil.NewResponse(c, http.StatusForbidden, "error", "role not found in token", nil)
			return
		}

		if role != string(requiredRole) {
			s.logger.Errorw("Unauthorized role access", "role", role, "requiredRole", requiredRole)
			httputil.NewResponse(c, http.StatusForbidden, "error", fmt.Sprintf("access restricted to %s only", requiredRole), nil)
			return
		}

		c.Next()
	}
}

func (s *server) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBind(&req); err != nil {
		s.logger.Errorw("Failed to bind user", "err", err)
		httputil.NewResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	tx := s.db.Begin()
	gw := gateway.NewGateway(tx)

	user, err := gw.GetUserByName(req.Username)
	if err != nil {
		tx.Rollback()
		s.logger.Errorw("User not found", "err", err)
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}

	if checkPassword(user.Password, req.Password) {
		tx.Rollback()
		s.logger.Error("Authentication failed")
		httputil.NewResponse(c, http.StatusUnauthorized, "error", "authentication failled", nil)
		return
	}

	tokens, err := createBothTokens(user.Username, string(user.Role))
	if err != nil {
		tx.Rollback()
		s.logger.Errorw("Failed to create tokens", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
	}

	claims, _ := validateJWT(tokens.RefreshToken)

	var refreshToken = domain.Token{
		UUID:     claims["jit"].(string),
		Username: claims["login"].(string),
		Active:   true,
	}

	if err = gw.DeactivateRefreshTokens(refreshToken.Username); err != nil {
		tx.Rollback()
		s.logger.Errorw("Failed to deactivate previous refresh tokens", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	if err = gw.SaveRefreshToken(refreshToken); err != nil {
		tx.Rollback()
		s.logger.Errorw("Failed to save refresh token", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	tx.Commit()

	httputil.NewResponse(c, http.StatusOK, "success", "", tokens)
}

func (s *server) Register(c *gin.Context) {
	var user = domain.User{}
	if err := c.ShouldBind(&user); err != nil {
		s.logger.Errorw("Failed to bind user", "err", err)
		httputil.NewResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	if user.Role != domain.RoleWorker && user.Role != domain.RoleDispatcher {
		s.logger.Errorw("Invalid role")
		httputil.NewResponse(c, http.StatusBadRequest, "error", "invalid role", nil)
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		s.logger.Errorw("Failed to hash password", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	user.Password = hash

	tx := s.db.Begin()
	gw := gateway.NewGateway(tx)

	if err = gw.SaveUser(user); err != nil {
		tx.Rollback()
		s.logger.Errorw("Failed to save user", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	tokens, err := createBothTokens(user.Username, string(user.Role))
	if err != nil {
		tx.Rollback()
		s.logger.Errorw("Failed to create tokens", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	claims, _ := validateJWT(tokens.RefreshToken)

	var refreshToken = domain.Token{
		UUID:     claims["jit"].(string),
		Username: claims["login"].(string),
		Active:   true,
	}

	if err = gw.DeactivateRefreshTokens(refreshToken.Username); err != nil {
		tx.Rollback()
		s.logger.Errorw("Failed to deactivate previous refresh tokens", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	if err = gw.SaveRefreshToken(refreshToken); err != nil {
		tx.Rollback()
		s.logger.Errorw("Failed to save refresh token", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	tx.Commit()

	httputil.NewResponse(c, http.StatusOK, "success", "", tokens)
}

func (s *server) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBind(&req); err != nil {
		s.logger.Errorw("Failed to bind request", "err", err)
		httputil.NewResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	claims, err := validateJWT(req.RefreshToken)
	if err != nil {
		s.logger.Errorw("Invalid refresh token", "err", err)
		httputil.NewResponse(c, http.StatusUnauthorized, "error", "invalid refresh token", nil)
		return
	}

	if claims["type"] != "refresh" {
		s.logger.Errorw("Invalid token type for refresh")
		httputil.NewResponse(c, http.StatusUnauthorized, "error", "invalid token type", nil)
		return
	}

	username, ok1 := claims["login"].(string)
	role, ok2 := claims["role"].(string)
	if !ok1 || !ok2 {
		s.logger.Errorw("Failed to extract claims from token")
		httputil.NewResponse(c, http.StatusInternalServerError, "error", "invalid token claims", nil)
		return
	}

	tokens, err := createBothTokens(username, role)
	if err != nil {
		s.logger.Errorw("Failed to create tokens", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	claims2, _ := validateJWT(tokens.RefreshToken)

	var refreshToken = domain.Token{
		UUID:     claims2["jit"].(string),
		Username: claims2["login"].(string),
		Active:   true,
	}

	tx := s.db.Begin()
	gw := gateway.NewGateway(tx)

	if err = gw.DeactivateRefreshTokens(refreshToken.Username); err != nil {
		tx.Rollback()
		s.logger.Errorw("Failed to deactivate previous refresh tokens", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	if err = gw.SaveRefreshToken(refreshToken); err != nil {
		tx.Rollback()
		s.logger.Errorw("Failed to save refresh token", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	tx.Commit()

	httputil.NewResponse(c, http.StatusOK, "success", "", tokens)
}

func (s *server) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBind(&req); err != nil {
		s.logger.Errorw("Failed to bind request", "err", err)
		httputil.NewResponse(c, http.StatusBadRequest, "error", "invalid request payload", nil)
		return
	}

	claims, _ := validateJWT(req.RefreshToken)
	username := claims["login"].(string)

	gw := gateway.NewGateway(s.db)
	if err := gw.DeactivateRefreshTokens(username); err != nil {
		s.logger.Errorw("Failed to deactivate refresh token", "err", err)
		httputil.NewResponse(c, http.StatusInternalServerError, "error", "failed to logout", nil)
		return
	}

	httputil.NewResponse(c, http.StatusOK, "success", "logged out successfully", nil)
}

func createBothTokens(login, role string) (domain.JWT, error) {
	accessToken, err := createAccessToken(login, role)
	if err != nil {
		return domain.JWT{}, err
	}
	refreshToken, err := createRefreshToken(login, role)
	if err != nil {
		return domain.JWT{}, err
	}

	var data domain.JWT
	data.AccessToken = accessToken
	data.RefreshToken = refreshToken

	return data, nil
}

func createAccessToken(login, role string) (string, error) {
	payload := map[string]any{
		"login": login,
		"role":  role,
		"type":  "access",
		"jit":   uuid.New(),
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}

	return generateJWT(payload)
}

func createRefreshToken(login, role string) (string, error) {
	payload := map[string]any{
		"login": login,
		"role":  role,
		"type":  "refresh",
		"jit":   uuid.New(),
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	return generateJWT(payload)
}

func generateJWT(payload map[string]any) (string, error) {
	header := map[string]any{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", errors.New("error marshalling header token")
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", errors.New("error marshalling payload token")
	}

	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	signatureInput := fmt.Sprintf("%s.%s", headerBase64, payloadBase64)

	secretKey := viper.GetString("jwt.secret")
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signatureInput))
	signature := h.Sum(nil)

	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	token := fmt.Sprintf("%s.%s.%s", headerBase64, payloadBase64, signatureBase64)
	return token, nil
}

func validateJWT(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("failed to decode payload")
	}

	signatureInput := fmt.Sprintf("%s.%s", parts[0], parts[1])

	secretKey := viper.GetString("jwt.secret")
	expectedSignature := hmac.New(sha256.New, []byte(secretKey))
	expectedSignature.Write([]byte(signatureInput))
	expectedSignatureBase64 := base64.RawURLEncoding.EncodeToString(expectedSignature.Sum(nil))

	if parts[2] != expectedSignatureBase64 {
		return nil, errors.New("invalid token signature")
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("failed to parse claims")
	}

	claims["exp"] = int64(claims["exp"].(float64))

	if exp, ok := claims["exp"].(int64); ok {
		if exp < time.Now().Unix() {
			return nil, errors.New("token expired")
		}
	}

	return claims, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil
}
