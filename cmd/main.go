package main

import (
	//"context"
	//"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"go.uber.org/zap"
	"log"
	"myproject/api"
	"myproject/dao"
	//"myproject/middleware"
)

func main() {
	dao.Init()
	h := server.Default() // 创建engine

	//初始日志，记录错误信息
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("无法同步日志记录器: %v", err)
		}
	}()

	// 允许跨域请求
	//h.Use(middleware.CORS())
	// 配置 CORS 中间件（允许所有来源）
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	api.Router(h)
	// 启动服务器
	h.Spin() //// 启动监听，Hertz默认是8888端口
}
