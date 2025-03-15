package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"myproject/service"
	"net/http"
	"strconv"
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

// SearchProduct 搜索商品
func SearchProduct(ctx context.Context, c *app.RequestContext) {
	productName := c.Param("product_name")
	products, err := service.SearchProduct(productName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"info": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 10000, "info": "success", "data": gin.H{"products": products}})
}

// GetProductDetails 获取商品详情
func GetProductDetails(ctx context.Context, c *app.RequestContext) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"status": err.Error()})
	}
	product, err := service.GetProductDetails(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 10000, "info": "success", "data": utils.H{"product": product}})
}
func AddProductTOCart(ctx context.Context, c *app.RequestContext) {
	var cartItem struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
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
