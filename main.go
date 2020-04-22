package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/open-kingfisher/king-inspect/check/all"
	"github.com/open-kingfisher/king-inspect/router"
	"github.com/open-kingfisher/king-utils/common"
	"github.com/open-kingfisher/king-utils/common/log"
	"github.com/open-kingfisher/king-utils/common/rabbitmq"
	"github.com/open-kingfisher/king-utils/config"
	"github.com/open-kingfisher/king-utils/kit"
	_ "github.com/open-kingfisher/king-utils/middleware/Validator"
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
