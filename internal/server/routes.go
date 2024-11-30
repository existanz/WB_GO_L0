package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/orders", s.GetOrdersHandler)
	r.GET("/orders/:id", s.GetOrderHandler)

	return r
}

func (s *Server) GetOrdersHandler(c *gin.Context) {
	orders, err := s.db.GetOrdersPlain()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, orders)
}

func (s *Server) GetOrderHandler(c *gin.Context) {
	id := c.Param("id")
	order, err := s.db.GetOrder(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.Data(http.StatusOK, gin.MIMEJSON, []byte(order))
}
