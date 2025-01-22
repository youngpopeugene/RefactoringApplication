package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type server struct {
	router *gin.Engine
	group  *gin.RouterGroup
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func New(router *gin.Engine, db *gorm.DB, logger *zap.SugaredLogger) *server {
	return &server{
		router: router,
		db:     db,
		logger: logger,
	}
}

func Run(r *gin.Engine, db *gorm.DB, logger *zap.SugaredLogger) {
	server := New(r, db, logger)
	server.group = server.router.Group("/api/v1")
	server.group.GET("/ping", ping)
	server.Router()
	server.AuthRouter()

	if err := server.router.Run(":8080"); err != nil {
		logger.Fatal("Failed to run router", "err", err)
	}
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
