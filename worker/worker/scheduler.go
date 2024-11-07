package worker

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ChuhC/crontab/common"
	. "github.com/ChuhC/crontab/worker/logger"
)

type Scheduler struct {
	JobEventChan    chan *common.JobEvent
	jobPlanTable    map[string]*common.JobSchedulePlan
	jobExecutingJob map[string]*common.JobExecuteInfo
}

var G_Scheduler *Scheduler

func InitScheduler() {
	sch := &Scheduler{
		JobEventChan: make(chan *common.JobEvent, 65535),
		jobPlanTable: make(map[string]*common.JobSchedulePlan),
	}

	go sch.ScheduleLoop()

	G_Scheduler = sch

}

func (s *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {

}

func (s *Scheduler) RunCommand(job *common.Job) error {
	cmd := strings.Split(job.Command, " ")
	exc := exec.Command(cmd[0], cmd[1:]...)
	return exc.Run()
}

func (s *Scheduler) TrySchedule() time.Duration {

	fmt.Println(" in try schedule", len(s.jobPlanTable))
	//sleep random time
	if len(s.jobPlanTable) == 0 {
		return 1 * time.Second
	}

	var nearTime *time.Time
	// range all jobs
	for _, jobPlan := range s.jobPlanTable {
		fmt.Println(jobPlan.NextTime)
		fmt.Println(jobPlan.Job.CronExpr)
		if jobPlan.NextTime.Before(time.Now()) || jobPlan.NextTime.Equal(time.Now()) {
			LG.Debug("执行任务, %s", jobPlan.Job.Name)
			err := s.RunCommand(jobPlan.Job)
			if err != nil {
				LG.Error("command: %s exec failed. err: %s", jobPlan.Job.Command, err.Error())
			}
			jobPlan.NextTime = jobPlan.Expr.Next(time.Now())
		}

		//统计最近一下要过期的时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	scheduleAfter := (*nearTime).Sub(time.Now())
	return scheduleAfter
}

func (s *Scheduler) ScheduleLoop() {

	scheduleAfter := s.TrySchedule()

	timer := time.NewTimer(scheduleAfter)
	for {
		select {
		case event := <-s.JobEventChan:
			switch event.EventType {
			case common.JOB_EVENT_SAVE:
				LG.Debug("recv event save. name: %s", event.Job.Name)
				jobPlan, err := common.NewJobSchedulePlan(event.Job)
				if err != nil {
					LG.Error("NewJobSchedulePlan err, ", err.Error())
					break
				}
				fmt.Println(jobPlan.Job.Name)
				fmt.Println(jobPlan.Job.CronExpr)

				s.jobPlanTable[event.Job.Name] = jobPlan

			case common.JOB_EVENT_DELETE:

				LG.Debug("recv event delete. name: %s", event.Job.Name)
				_, existed := s.jobPlanTable[event.Job.Name]
				if existed {
					delete(s.jobPlanTable, event.Job.Name)
				}

			}
		case <-timer.C:
		}
		scheduleAfter = s.TrySchedule()
		timer.Reset(scheduleAfter)

	}
}

func (s *Scheduler) PushScheduler(event *common.JobEvent) {
	s.JobEventChan <- event
}

func (s *Scheduler) handleScheduler(event *common.JobEvent) {

}
