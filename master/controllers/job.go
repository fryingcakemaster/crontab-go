package controllers

import (
	"fmt"

	"github.com/ChuhC/crontab/common"
	"github.com/ChuhC/crontab/master/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang.org/x/net/context"
)

// job controller
type JobController struct {
	beego.Controller
}

// @Description save job to etcd
// @Param	name     query 	    string	true        "crontab task name"
// @Param	command  query 	    string	true        "crontab command"
// @Param	cronExpr query 	    string	true        "crontab expression"
// @Success 200 {string} save job success
// @Failure 403
// @router /save [post]
func (j *JobController) Save() {
	var (
		err    error
		oldJob *common.Job
	)
	var job = &common.Job{}
	err = j.ParseForm(job)
	if err != nil {
		logs.Error("parse form error! err: %s", err.Error())
		return
	}

	if oldJob, err = models.G_JobMgr.SaveJob(context.Background(), job); err != nil {
		logs.Error("save job error! err: %s", err.Error())
		return
	}

	resp := common.BuildResponse(0, "success", oldJob)
	logs.Debug("save job success. oldJob: %+v", oldJob)
	fmt.Println(resp)
	j.Ctx.ResponseWriter.Write(resp)
}

// @Description del job from etcd
// @Param	name     query 	    string	true        "crontab task name"
// @Success 200 {string} delete job success
// @Failure 403
// @router /del [post]
func (j *JobController) Del() {
	var (
		err    error
		oldJob *common.Job
		resp   []byte
	)
	jobName := j.GetString("name")

	oldJob, err = models.G_JobMgr.DelJob(context.Background(), jobName)
	if err != nil {
		goto ERR
	}

	if oldJob == nil {
		err = fmt.Errorf("job not exist")
		goto ERR
	}

	logs.Debug("del job success! oldJob: %+v", oldJob)
	resp = common.BuildResponse(0, "success", oldJob)
	j.Ctx.ResponseWriter.Write(resp)
	return

ERR:
	logs.Error("del job failed! err: %s", err.Error())
	resp = common.BuildResponse(-1, err.Error(), oldJob)
	j.Ctx.ResponseWriter.Write(resp)
	return
}

// @Title GetAll
// @Description get all jobs
// @Success 200 {string} get all jobs success
// @Failure 403
// @router /list [get]
func (j *JobController) GetAll() {
	var resp []byte
	jobList, err := models.G_JobMgr.GetAll(context.Background())
	if err != nil {
		goto ERR
	}
	logs.Debug("get all job success! jobList: %+v", jobList)
	resp = common.BuildResponse(0, "success", jobList)
	j.Ctx.ResponseWriter.Write(resp)
	return

ERR:
	logs.Debug("get all job failed! err: %s", err.Error())
	resp = common.BuildResponse(-1, err.Error(), jobList)
	j.Ctx.ResponseWriter.Write(resp)

	return
}

// @Title killjob
// @Description kill a job
// @Param	name     query 	    string	true        "crontab task name"
// @Success 200 {string} kill a job success
// @Failure 403
// @router /kill [post]
func (j *JobController) Kill() {
	var (
		err  error
		resp []byte
	)
	jobName := j.GetString("name")
	err = models.G_JobMgr.KillJob(jobName)
	if err != nil {

	}
	err = models.G_JobMgr.KillJob(jobName)
	if err != nil {
		goto ERR
	}

	logs.Error("kill job succeed! name: %s", jobName)
	resp = common.BuildResponse(0, "success", nil)
	j.Ctx.ResponseWriter.Write(resp)
	return

ERR:
	logs.Error("kill job failed! err: %s", err.Error())
	resp = common.BuildResponse(-1, err.Error(), nil)
	j.Ctx.ResponseWriter.Write(resp)
	return

}
