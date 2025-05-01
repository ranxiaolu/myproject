package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"myproject/middleware"
)

func Router(h *server.Hertz) {
	//配置用户相关路由
	//注册
	h.POST("user/register", RegisterUser)
	//登录
	h.GET("user/token", LoginUser)
	//刷新token
	h.GET("user/token/refresh", middleware.AuthMiddleware(), RefreshUserToken)
	//修改密码
	h.PUT("user/password", middleware.AuthMiddleware(), ChangePassword)
	//获取用户信息
	h.GET("user/info/:user_id", middleware.AuthMiddleware(), GetUserInfo)
	//修改用户信息
	h.PUT("user/info", middleware.AuthMiddleware(), UpdateUserInfo)

	//配置商品相关路由
	//列出商品信息
	h.GET("product/list", GetProductList)
	//搜索商品
	h.GET("book/search/:product_name", middleware.AuthMiddleware(), SearchProduct)
	//商品添加到购物车
	h.PUT("product/addCart", middleware.AuthMiddleware(), AddProductTOCart)
	//获取购物车列表
	h.GET("product/cart", middleware.AuthMiddleware(), GetCartList)
	//获取商品详情
	h.GET("product/info/:product_id", middleware.AuthMiddleware(), GetProductDetails)
	//通过商品标签获得商品详情
	h.GET("product/:product_type", middleware.AuthMiddleware(), GetProductsByType)

	//配置评论相关路由
	//获得商品评论
	h.GET("comment/:product_id", GetComment)
	//给商品评论
	h.POST("comment/:product_id", middleware.AuthMiddleware(), AddComment)
	//删除商品评论
	h.DELETE("comment/:comment_id", middleware.AuthMiddleware(), DeleteComment)
	//更新评论内容
	h.PUT("comment/:comment_id", middleware.AuthMiddleware(), UpdateComment)

	//配置操作相关路由
	//点赞点踩
	h.PUT("comment/praise/:comment_id", middleware.AuthMiddleware(), PraiseComment)
	//下单
	h.POST("operate/order", middleware.AuthMiddleware(), CreateOrder)
}
