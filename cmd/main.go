package main

import (
	//"context"
	//"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/zap"
	"log"
	"myproject/api"
	"myproject/middleware"
	//"myproject/middleware"
)

func main() {
	//初始日志，记录错误信息
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("无法同步日志记录器: %v", err)
		}
	}()
	h := server.Default() // 创建engine
	//h.Use(func(c *hertz.Context) {
	//	logger.Info("Request received",
	//		zap.String("method", c.Request.Method),
	//		zap.String("path", c.Request.URL.Path),
	//	)
	//	c.Next()
	//})
	// 自定义日志中间件
	//h.Use(func(c context.Context, ctx *app.RequestContext) {
	//	logger.Info("Request received",
	//		zap.String("method", ctx.Request.Method),
	//		zap.String("path", ctx.Request.URI().Path()),
	//	)
	//	ctx.Next(c)
	//})

	// 允许跨域请求
	//h.Use(middleware.CORS())
	//配置用户相关路由
	//注册
	h.POST("user/register", api.RegisterUser)
	//登录
	h.GET("user/token", api.LoginUser)
	//刷新token
	h.GET("user/token/refresh", middleware.AuthMiddleware(), api.RefreshUserToken)
	//修改密码
	h.PUT("user/password", middleware.AuthMiddleware(), api.ChangePassword)
	//获取用户信息
	h.GET("user/info/:user_id", middleware.AuthMiddleware(), api.GetUserInfo)
	//修改用户信息
	h.PUT("user/info", middleware.AuthMiddleware(), api.UpdateUserInfo)

	//配置商品相关路由
	//列出商品信息
	h.GET("product/list", api.GetProductList)
	//搜索商品
	h.GET("book/search/:product_name", middleware.AuthMiddleware(), api.GetProductDetails)
	//商品添加到购物车
	h.PUT("product/addCart", middleware.AuthMiddleware(), api.AddProductTOCart)
	//获取购物车列表
	h.GET("product/cart", middleware.AuthMiddleware(), api.GetCartList)
	//获取商品详情
	h.GET("product/info/:product_name", middleware.AuthMiddleware(), api.GetProductDetails)
	//通过商品标签获得商品详情
	h.GET("product/:product_type", middleware.AuthMiddleware(), api.GetProductsByType)

	//配置评论相关路由
	//获得商品评论
	h.GET("comment/:product_id", api.GetComment)
	//给商品评论
	h.POST("comment/:product_id", middleware.AuthMiddleware(), api.AddComment)
	//删除商品评论
	h.DELETE("comment/Basketball", middleware.AuthMiddleware(), api.DeleteComment)
	//更新评论内容
	h.PUT("comment/:comment_id", middleware.AuthMiddleware(), api.UpdateComment)

	//配置操作相关路由
	//点赞点踩
	h.PUT("comment/praise", middleware.AuthMiddleware(), api.PraiseComment)
	//下单
	h.POST("operate/order", middleware.AuthMiddleware(), api.CreateOrder)
	// 启动服务器
	h.Spin() //// 启动监听，Hertz默认是8888端口
}
