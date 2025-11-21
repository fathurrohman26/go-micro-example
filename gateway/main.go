package main

import (
	"elproxy.cloud/elproxy-server/common/model"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v5"
	"go-micro.dev/v5/client"
)

type Gateway struct {
	client client.Client
}

func (g *Gateway) GetUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "missing user id"})
		return
	}

	req := g.client.NewRequest("users", "UsersService.Get", &model.GetUserRequest{ID: id})
	rsp := &model.GetUserResponse{}

	if err := g.client.Call(c, req, rsp); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, rsp)
}

func (g *Gateway) CreateOrder(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	rpcReq := g.client.NewRequest("orders", "OrdersService.Create", &req)
	rpcRsp := &model.CreateOrderResponse{}

	if err := g.client.Call(c, rpcReq, rpcRsp); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, rpcRsp)
}

func main() {
	svc := micro.NewService(
		micro.Name("api.gateway"),
	)
	svc.Init()

	gtw := &Gateway{
		client: svc.Client(),
	}

	app := gin.Default()

	app.GET("/users", gtw.GetUser)
	app.POST("/orders", gtw.CreateOrder)

	app.Run(":8080")
}
