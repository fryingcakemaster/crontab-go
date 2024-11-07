package worker

import (
	"fmt"
	"time"

	"github.com/ChuhC/crontab/common"
	. "github.com/ChuhC/crontab/worker/logger"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
)

var G_JobMgr *JobMgr

// job manager
type JobMgr struct {
	Client *clientv3.Client // etcd client
}

// init etcd connect
func InitJobMgr(endPoints []string, dialTimeout int) {
	timeout := time.Duration(dialTimeout) * time.Millisecond
	config := clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: timeout,
	}
	client, err := clientv3.New(config)
	if err != nil {
		panic(fmt.Errorf("ectd client init failed! err: %s ", err.Error()))
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = client.Status(timeoutCtx, config.Endpoints[0])
	if err != nil {
		panic(fmt.Errorf("error checking etcd status: %v", err))
	}

	LG.Debug("%s", "etcd init succeed!")
	G_JobMgr = &JobMgr{Client: client}
	G_JobMgr.watchJobs()
}

func (j *JobMgr) watchJobs() error {

	getResp, err := j.Client.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// read all jobs
	for _, kv := range getResp.Kvs {

		job, err := common.UnpackJob(kv.Value)
		if err != nil {
			err = nil
			continue
		}
		//build a save event
		evt := common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
		G_Scheduler.PushScheduler(evt)
	}

	// start watching from this reversion
	go func() {
		startRevision := getResp.Header.Revision + 1
		watchChan := j.Client.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(startRevision), clientv3.WithPrefix())
		for watchResp := range watchChan {
			for _, watchEvent := range watchResp.Events {
				var evt *common.JobEvent
				switch watchEvent.Type {
				case mvccpb.PUT:
					job, err := common.UnpackJob(watchEvent.Kv.Value)
					if err != nil {
						continue
					}
					//build a save event
					evt = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)

				case mvccpb.DELETE:
					// Delete /cron/jobs/job10
					jobName := common.ExtractJobName(string(watchEvent.Kv.Key))
					job := &common.Job{Name: jobName}
					//build a delete event
					evt = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				G_Scheduler.PushScheduler(evt)
			}
		}
	}()
	return nil
}
