package common

import (
	"encoding/json"
	"strings"
)

type Job struct {
	Name     string `json:"name" form:"name"`
	CronExpr string `json:"cronExpr" form:"cronExpr"`
	Command  string `json:"command" form:"command"`
}

// HTTP response
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func BuildResponse(errno int, msg string, data interface{}) []byte {
	var response Response
	response.Errno = errno
	response.Msg = msg
	response.Data = data

	resp, _ := json.Marshal(response)
	return resp
}

func UnpackJob(value []byte) (ret *Job, err error) {
	job := &Job{}
	err = json.Unmarshal(value, job)
	if err != nil {
		return
	}

	ret = job
	return
}

func ExtractJobName(key string) string {
	return strings.TrimPrefix(key, JOB_SAVE_DIR)
}
