package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/luonannet/playground-backend/models"
	"github.com/luonannet/playground-backend/routers"
	"github.com/luonannet/playground-backend/util"
	"github.com/sirupsen/logrus"
)

func init() {

	util.InitConfig("")
}
func main() {
	if util.GetConfig().AppModel == "dev" || util.GetConfig().AppModel == "debug" {
		util.Log.SetLevel(logrus.DebugLevel)
		util.GetConfig().AppModel = gin.DebugMode
	} else {
		util.Log.SetLevel(logrus.InfoLevel)
		util.GetConfig().AppModel = gin.ReleaseMode
	}

	//init database
	err := util.InitDB()
	if err != nil {
		util.Log.Errorf("database connect failed , err:%v\n", err)
		return
	}

	//init Authing.cn config
	models.InitAuthing("", "")
	//init kubernetes client-go
	// models.InitK8sClient()
	util.InitStatisticsLog()
	//startup a webscoket server to wait client ws
	// go models.StartWebSocket()
	gin.SetMode(util.GetConfig().AppModel)
	r := routers.InitRouter()
	address := fmt.Sprintf(":%d", util.GetConfig().AppPort)
	util.Log.Infof(" startup meta http service at port %s .and %s mode \n", address, util.GetConfig().AppModel)
	if err := r.Run(address); err != nil {
		util.Log.Infof("startup meta  http service failed, err:%v\n", err)
	}
}
