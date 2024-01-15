package api

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) RenderRoutes() {
	router := gin.Default()
	authRoutes := router.
		Group("/").
		Use(authMiddleware(s.tokenMaker))

	router.POST("/users/login", s.handlerLoginUser)
	router.POST("/users", s.handlerCreateUser)

	authRoutes.POST("/accounts", s.handlerCreateAccount)
	authRoutes.GET("/accounts/:id", s.handlerGetAccount)
	authRoutes.GET("/accounts", s.handlerListAccount)

	authRoutes.POST("/transfers", s.handlerCreateTransfer)

	s.router = router
}
