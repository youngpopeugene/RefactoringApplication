package rest

import (
	"app/internal/models"
	"github.com/gin-gonic/gin"
	jwt "github.com/kyfk/gin-jwt"
)

func (h *handler) NewAuth() (jwt.Auth, error) {
	return jwt.New(jwt.Auth{
		SecretKey: []byte("secret"),
		Authenticator: func(c *gin.Context) (jwt.MapClaims, error) {
			var req struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			if err := c.ShouldBind(&req); err != nil {
				return nil, jwt.ErrorAuthenticationFailed
			}

			var user models.User
			h.db.Where("username = ?", req.Username).First(&user)

			if CheckPasswordHash(req.Password, user.Password) {
				return nil, jwt.ErrorAuthenticationFailed
			}

			return jwt.MapClaims{
				"username": user.Username,
				"role":     user.Role,
			}, nil
		},
		UserFetcher: func(c *gin.Context, claims jwt.MapClaims) (interface{}, error) {
			username, ok := claims["username"].(string)
			if !ok {
				return nil, nil
			}
			var user = models.User{
				Username: username,
			}
			h.db.First(&user)
			return user, nil
		},
	})
}

func (h *handler) Register(c *gin.Context) {
	var user = models.User{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, err)
		return
	}

	hash, _ := HashPassword(user.Password)
	user.Password = hash

	h.db.Create(&user)
	c.JSON(200, user)
}

func Worker(m jwt.Auth) gin.HandlerFunc {
	return m.VerifyPerm(func(claims jwt.MapClaims) bool {
		return role(claims) == models.RoleWorker
	})
}

func Dispatcher(m jwt.Auth) gin.HandlerFunc {
	return m.VerifyPerm(func(claims jwt.MapClaims) bool {
		return role(claims) == models.RoleDispatcher
	})
}

func role(claims jwt.MapClaims) models.Role {
	return models.Role(claims["role"].(string))
}
