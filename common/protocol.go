package common

import (
	"time"

	"github.com/gorhill/cronexpr"
)

type JobEvent struct {
	EventType int
	Job       *Job
}

func BuildJobEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

// 任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 // 要调度的任务信息
	Expr     *cronexpr.Expression // 解析好的cronexpr表达式
	NextTime time.Time            // 下次调度时间
}

func NewJobSchedulePlan(job *Job) (*JobSchedulePlan, error) {
	expr, err := cronexpr.Parse(job.CronExpr)
	if err != nil {
		return nil, err
	}

	return &JobSchedulePlan{Job: job, Expr: expr, NextTime: expr.Next(time.Now())}, nil

}

//任务执行情况
type JobExecuteInfo struct {
	Job      *Job
	PlanTime time.Time //理论调度时间
	RealTime time.Time //实际调度时间
}

func NewJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) *JobExecuteInfo {
	jobExecuteInfo := &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime,
		RealTime: time.Now(),
	}

	return jobExecuteInfo
}
