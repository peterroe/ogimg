package server

import (
	"ogimg/internal/handler"
	"ogimg/internal/middleware"
	"ogimg/pkg/log"

	"github.com/gin-gonic/gin"
)

func NewServerHTTP(
	logger *log.Logger,
	userHandler *handler.UserHandler,
	imageHandler *handler.ImageHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(
		middleware.CORSMiddleware(),
	)
	r.GET("/", imageHandler.GetOgImageByUrl)
	r.GET("/desc", imageHandler.GetOgDescByUrl)
	r.GET("/user", userHandler.GetUserById)

	return r
}
