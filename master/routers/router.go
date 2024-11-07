package routers

import (
	"github.com/ChuhC/crontab/master/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:Get")
	beego.Router("/job/list", &controllers.JobController{}, "get:GetAll")
	beego.Router("/job/save", &controllers.JobController{}, "post:Save")
	beego.Router("/job/delete", &controllers.JobController{}, "post:Del")
	beego.Router("/job/kill", &controllers.JobController{}, "post:Kill")
}
