package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"mediahub/controller"
	"mediahub/middleware"
	"mediahub/pkg/config"
	"mediahub/pkg/log"
	"mediahub/pkg/storage/cos"
	"mediahub/routers"
	"net/http"
)

var configFile = flag.String("config", "dev.config.yaml", "")

func main() {
	flag.Parse()
	//初始化配置文件
	config.InitConfig(*configFile)
	cnf := config.GetConfig()

	log.SetLevel(cnf.Log.Level)
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetPrintCaller(true)

	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)

	sf := cos.NewCosStorageFactory(cnf.Cos.BucketUrl, cnf.Cos.SecretId, cnf.Cos.SecretKey, cnf.Cos.CDNDomain)
	controller := controller.NewController(sf, logger, cnf)

	gin.SetMode(cnf.Http.Mode)
	r := gin.Default()
	r.Use(middleware.Cors(), middleware.Auth())
	r.GET("/health", func(*gin.Context) {})
	api := r.Group("/api")
	routers.InitRouters(api, controller)

	fs := http.FileServer(http.Dir("www"))
	r.NoRoute(func(ctx *gin.Context) {
		fs.ServeHTTP(ctx.Writer, ctx.Request)
	})
	r.GET("/", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "www/index.html")
	})

	err := r.Run(fmt.Sprintf("%s:%d", cnf.Http.IP, cnf.Http.Port))
	if err != nil {
		log.Error(err)
	}
}
