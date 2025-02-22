package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"myproject/service"
	"net/http"
)

// GetProductList 获取商品列表
func GetProductList(ctx context.Context, c *app.RequestContext) {

	products, err := service.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.H{"status": 10000, "info": "success", "data": utils.H{"products": products}})
}
func GetProductDetails(ctx context.Context, c *app.RequestContext) {
	productName := c.Param("product_name")

	product, err := service.GetProductDetails(productName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 10000, "info": "success", "data": utils.H{"product": product}})
}
func AddProductTOCart(ctx context.Context, c *app.RequestContext) {
	var cartItem struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.BindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 10001, "info": "invalid request"})
		return
	}
	//userID := c.GetString("user_id")
	//将商品添加到购物车中
	if err := service.AddProductTOCart(cartItem.ProductID, cartItem.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 10002, "info": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 10000, "info": "success"})

}
func GetCartList(ctx context.Context, c *app.RequestContext) {
	userID := c.GetString("user_id")
	cartItems, err := service.GetCartList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.H{"status": 10000, "info": "success", "data": utils.H{"products": cartItems}})
}
func GetProductsByType(ctx context.Context, c *app.RequestContext) {
	productType := c.Param("product_type")
	products, err := service.GetProductsByType(productType)
	if err != nil {
		c.JSON(400, utils.H{"info": "wrong"})
	}
	c.JSON(200, utils.H{"status": 10000, "info": "success", "data": utils.H{"products": products}})
}
