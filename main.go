package main

import (
	"github.com/gin-gonic/gin"
	"kingfisher/kf/common"
	"kingfisher/kf/common/log"
	"kingfisher/kf/common/rabbitmq"
	"kingfisher/kf/config"
	"kingfisher/kf/kit"
	_ "kingfisher/kf/middleware/Validator"
	_ "kingfisher/king-inspect/check/all"
	"kingfisher/king-inspect/router"
)

func main() {
	// Debug Mode
	gin.SetMode(config.Mode)
	g := gin.New()
	// 设置路由
	r := router.SetupRouter(kit.EnhanceGin(g))
	// 通过消息中间件获取更新kubeConfig文件消息
	consumer := rabbitmq.Consumer{
		Address:      config.RabbitMQURL,
		ExchangeName: common.UpdateKubeConfig,
		Handler:      &rabbitmq.UpdateKubeConfig{},
	}
	go consumer.Run()
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(config.Listen); err != nil {
		log.Fatalf("Listen error: %v", err)
	}
}
