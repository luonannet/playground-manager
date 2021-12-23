package task

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"playground_backend/handler"
)

// start task
func StartTask() {
	toolbox.StartTask()
}

func StopTask() {
	toolbox.StopTask()
}

//InitYamlTask Get yaml data source
func ClearInstanceTask(clInvalidInstance string) {
	invalidTask := toolbox.NewTask("ClearInvaildResource", clInvalidInstance, handler.ClearInvaildResource)
	toolbox.AddTask("ClearInvaildResource", invalidTask)
}

//InitTask Timing task initialization
func InitTask() bool {
	// Get the original yaml data
	clInvalidInstanFlag, err := beego.AppConfig.Int("crontab::cl_invalid_instances_flag")
	if clInvalidInstanFlag == 1 && err == nil {
		clInvalidInstance := beego.AppConfig.String("crontab::cl_invalid_instances")
		ClearInstanceTask(clInvalidInstance)
	}
	return true
}