package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type MainController struct {
	beego.Controller
}

func (m *MainController) Get() {
	logs.Debug(" in index ......................")
	m.TplName = "index.html"
}
