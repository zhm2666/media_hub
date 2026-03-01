package main

import (
	"flag"
	"shorturl-crontab/cron"
	"shorturl-crontab/pkg/config"
	"shorturl-crontab/pkg/db/mysql"
	"shorturl-crontab/pkg/db/redis"
	"shorturl-crontab/pkg/log"
)

var (
	configFile = flag.String("config", "dev.config.yaml", "")
)

func main() {
	flag.Parse()
	//初始化配置文件
	config.InitConfig(*configFile)
	cnf := config.GetConfig()

	log.SetLevel(cnf.Log.Level)
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetPrintCaller(true)

	//初始化MySQL
	mysql.InitMysql(cnf)
	//初始redis
	redis.InitRedisPool(cnf)

	cron.Run()
}
