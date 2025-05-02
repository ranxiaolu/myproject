package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"myproject/model"
	"myproject/service"
	"net/http"
	"strconv"
)

// CreateOrder 添加订单
func CreateOrder(ctx context.Context, c *app.RequestContext) {
	//从请求中获取用户 ID
	userIDStr := c.GetString("user_id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)
	//定义一个结构体 orderData 来绑定请求的 JSON 数据，该结构体包含产品列表
	var orderData struct {
		Products []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"products"`
	}
	//绑定数据到orderData
	if err := c.BindJSON(&orderData); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"info": "invalid request"})
		return
	}

	var products []model.Product
	//orderData 中的产品列表，调用 services.GetProductDetails 获取每个产品的详细信息。
	for _, item := range orderData.Products {
		product, err := service.GetProductDetails(item.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.H{"info": err.Error()})
			return
		}
		products = append(products, product)
	}
	//将获取到的产品添加到 products 列表中。
	orderID, err := service.CreateOrder(uint(userID), products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"info": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 10000, "info": "success", "data": utils.H{"order_id": orderID}})
}
