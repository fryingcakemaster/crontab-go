package models

import (
	"context"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
)

var (
	G_JobMgr *JobMgr
)

func init() {
	dialTimeout := 3 * time.Second
	hosts := []string{beego.AppConfig.String("etcdHost")}
	config := clientv3.Config{
		Endpoints:   hosts,
		DialTimeout: dialTimeout,
	}
	client, err := clientv3.New(config)
	if err != nil {
		panic(fmt.Errorf("ectd client init failed! err: %s ", err.Error()))
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	_, err = client.Status(timeoutCtx, config.Endpoints[0])
	if err != nil {
		panic(fmt.Errorf("error checking etcd status: %v", err))
	}
	logs.Info("etcd init succeed!")
	G_JobMgr = &JobMgr{Client: client}
}
