package logger

import (
	"fmt"
	"os"
	"path"

	"chufutao/bz/uniqued/utils"
	"github.com/Jeffail/gabs"
	"github.com/astaxie/beego/logs"
)

var LG *logs.BeeLogger

// init logger
func InitLogger(fileName string, daily bool, maxDays int, console bool, level int) {
	lg := logs.NewLogger(10000)
	lg.EnableFuncCallDepth(true)

	jsonObj := gabs.New()
	jsonObj.Set(fileName, "filename")
	jsonObj.Set(daily, "daily")
	jsonObj.Set(1<<31, "maxSize")     // 2GB
	jsonObj.Set(10000000, "maxLines") // 10,000, 000
	jsonObj.Set(maxDays, "maxDays")
	jsonObj.Set(level, "level")

	dir := path.Dir(fileName)
	ok, _ := utils.PathExists(dir)
	if !ok {
		err := os.Mkdir(dir, 0700)
		if err != nil {
			panic(fmt.Sprintf("cant't mkdir. err: %s", err.Error()))
		}
	}
	jsonObj.Set("0660", "perm")

	lg.SetLogger("file", jsonObj.String())

	if console == true {
		lg.SetLogger("console")
	}

	LG = lg
}
