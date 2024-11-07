package models

import (
	"context"
	"encoding/json"

	common "github.com/ChuhC/crontab/common"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
)

// job manager
type JobMgr struct {
	Client *clientv3.Client // etcd client
}

func (j *JobMgr) SaveJob(ctx context.Context, job *common.Job) (oldJob *common.Job, err error) {
	key := common.JOB_SAVE_DIR + "/" + job.Name
	var jobValue []byte
	jobValue, err = json.Marshal(job)
	if err != nil {
		return
	}

	var putResp *clientv3.PutResponse

	// save job to etcd
	putResp, err = j.Client.Put(ctx, key, string(jobValue), clientv3.WithPrevKV())
	if err != nil {
		return
	}
	// return old job if already exist
	if putResp.PrevKv != nil {
		var j = &common.Job{}
		err = json.Unmarshal(putResp.PrevKv.Value, j)
		oldJob = j
		logs.Debug("old value: %s", oldJob)
		return oldJob, err
	}
	return
}

// del job from etcd
func (j *JobMgr) DelJob(ctx context.Context, name string) (oldJob *common.Job, err error) {

	key := common.JOB_SAVE_DIR + "/" + name
	var dResp *clientv3.DeleteResponse
	dResp, err = j.Client.Delete(ctx, key, clientv3.WithPrevKV())
	if err != nil {
		logs.Info("1111111111111111111111111111", err.Error())
		return
	}
	logs.Info("22222222222222222222222222222")
	if len(dResp.PrevKvs) != 0 {
		var jb = &common.Job{}
		logs.Debug("oldJob: %s", string(dResp.PrevKvs[0].Value))
		err = json.Unmarshal(dResp.PrevKvs[0].Value, jb)
		if err == nil {
			oldJob = jb
			logs.Debug("oldJob2: %+v", oldJob)
			return
		}
	}

	return
}

// get all jobs
func (j *JobMgr) GetAll(ctx context.Context) (jobList []common.Job, err error) {

	var (
		resp *clientv3.GetResponse
	)
	resp, err = j.Client.Get(ctx, common.JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	jobList = make([]common.Job, 0)
	for _, kv := range resp.Kvs {

		job := &common.Job{}
		err = json.Unmarshal(kv.Value, job)
		if err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, *job)
	}
	return
}

// kill a job
func (j *JobMgr) KillJob(jobName string) error {
	// write to key=/cron/killer
	killkey := common.JOB_KILLER_DIR + jobName
	lResp, err := j.Client.Grant(context.Background(), 1)
	if err != nil {
		return nil
	}

	_, err = j.Client.Put(context.Background(), killkey, "", clientv3.WithLease(lResp.ID))
	if err != nil {
		return err
	}

	return nil
}
