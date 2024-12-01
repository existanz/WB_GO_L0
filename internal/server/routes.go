package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/orders", s.GetOrdersHandler)
	r.GET("/orders/:id", s.GetOrderHandler)

	r.LoadHTMLGlob("web/template/*")

	r.GET("/web/orders/:id", s.getOrderHTML)
	r.GET("/search", s.getSearchHTML)

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

	c.JSON(http.StatusOK, order)
}

func (s *Server) getOrderHTML(c *gin.Context) {
	orderID := c.Param("id")
	order, err := s.db.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	var prevOrderID, nextOrderID string
	if orderIDInt, err := strconv.Atoi(orderID); err == nil {
		if orderIDInt > 1 {
			prevOrderID = strconv.Itoa(orderIDInt - 1)
		}
		nextOrderID = strconv.Itoa(orderIDInt + 1)
	}

	c.HTML(http.StatusOK, "order.html", gin.H{
		"Order":       order,
		"PrevOrderID": prevOrderID,
		"NextOrderID": nextOrderID,
	})
}

func (s *Server) getSearchHTML(c *gin.Context) {
	uid := c.Query("order_uid")

	if uid == "" {
		c.HTML(http.StatusOK, "search.html", nil)
		return
	} else if uid[0] != '"' {
		uid = fmt.Sprintf("\"%s\"", uid)
	}

	order, err := s.db.GetOrderByUID(uid)
	isFind := true
	if err != nil {
		isFind = false
	}

	c.HTML(http.StatusOK, "search.html", gin.H{
		"Order":  order,
		"isFind": isFind,
		"UID":    uid,
	})
}
