package api

import "github.com/gin-gonic/gin"

func (s *Server) RenderRoutes() {
	router := gin.Default()

	router.POST("/users", s.handlerCreateUser)
	router.POST("/users/login", s.handlerLoginUser)

	router.POST("/accounts", s.handlerCreateAccount)
	router.GET("/accounts/:id", s.handlerGetAccount)
	router.GET("/accounts", s.handlerListAccount)

	router.POST("/transfers", s.handlerCreateTransfer)

	s.router = router
}
