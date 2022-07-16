package main

import (
	"fmt"
	"golibrary/Dao"
	"golibrary/config"
	"golibrary/routers"
	// "library/controller"
)

func main() {

	// 载入配置文件	./config.json
	cfgPath := "./cfg.json"
	c, err := config.LoadConfig(cfgPath)
	if err != nil {
		fmt.Printf("载入配置文件错误:%v\n", err)
		return
	}
	// 初始化数据库
	err = Dao.InitDB(c)
	if err != nil {
		fmt.Printf("初始化数据库错误:%v\n", err)
		return
	}
	// 注册路由
	router := routers.SetupRouter()

	router.Run(fmt.Sprintf("%s:%s", c.Host, c.Port))
}
