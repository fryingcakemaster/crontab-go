package models

import (
	common "github.com/ChuhC/crontab/common"
)

type JobGroup struct {
	GroupName   string
	Jobs        []common.Job
	Description string
	CreateTime  string
	Authority
}

func (g *JobGroup) Drop() {

}

func (g *JobGroup) Rename() {

}
