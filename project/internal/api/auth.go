package api

import (
	"app/internal/domain"
	"github.com/gin-gonic/gin"
	jwt "github.com/kyfk/gin-jwt"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) NewAuth() (jwt.Auth, error) {
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

			var user domain.User
			s.db.Where("username = ?", req.Username).First(&user)

			if checkPasswordHash(req.Password, user.Password) {
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
			var user = domain.User{
				Username: username,
			}
			s.db.First(&user)
			return user, nil
		},
	})
}

func (s *server) Register(c *gin.Context) {
	var user = domain.User{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, err)
		return
	}

	hash, _ := hashPassword(user.Password)
	user.Password = hash

	s.db.Create(&user)
	c.JSON(200, user)
}

func Worker(m jwt.Auth) gin.HandlerFunc {
	return m.VerifyPerm(func(claims jwt.MapClaims) bool {
		return role(claims) == domain.RoleWorker
	})
}

func Dispatcher(m jwt.Auth) gin.HandlerFunc {
	return m.VerifyPerm(func(claims jwt.MapClaims) bool {
		return role(claims) == domain.RoleDispatcher
	})
}

func role(claims jwt.MapClaims) domain.Role {
	return domain.Role(claims["role"].(string))
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil
}
