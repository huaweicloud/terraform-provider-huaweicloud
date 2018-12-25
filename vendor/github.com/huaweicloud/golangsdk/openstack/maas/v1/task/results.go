package task

import "github.com/huaweicloud/golangsdk"

type Task struct {
	ID       int64  `json:"id"`
	TaskName string `json:"task_name"`
}

type TaskResult struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	SrcNode      SrcNode `json:"src_node"`
	DstNode      DstNode `json:"dst_node"`
	ThreadNum    int     `json:"thread_num"`
	Status       int     `json:"status"`
	Progress     float64 `json:"progress"`
	MigrateSpeed int64   `json:"migrate_speed"`
	EnableKMS    bool    `json:"enableKMS"`
	Description  string  `json:"description"`
	TotalSize    int64   `json:"total_size"`
	CompleteSize int64   `json:"complete_size"`
	StartTime    int64   `json:"start_time"`
	LeftTime     int64   `json:"left_time"`
	TotalTime    int64   `json:"total_time"`
	SuccessNum   int64   `json:"success_num"`
	FailNum      int64   `json:"fail_num"`
	TotalNum     int64   `json:"total_num"`
	SmnInfo      SmnInfo `json:"smnInfo"`
}

type SrcNode struct {
	Region    string   `json:"region"`
	ObjectKey []string `json:"object_key"`
	Bucket    string   `json:"bucket"`
}

type DstNode struct {
	Region    string `json:"region"`
	ObjectKey string `json:"object_key"`
	Bucket    string `json:"bucket"`
}

type SmnInfo struct {
	NotifyResult       string `json:"notifyResult"`
	NotifyErrorMessage string `json:"notifyErrorMessage"`
	TopicName          string `json:"topicName"`
}

type CreateResult struct {
	golangsdk.Result
}

func (r CreateResult) Extract() (*Task, error) {
	var s Task
	err := r.ExtractInto(&s)
	return &s, err
}

func (r CreateResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*TaskResult, error) {
	var s TaskResult
	err := r.ExtractInto(&s)
	return &s, err
}

func (r GetResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type DeleteResult struct {
	golangsdk.ErrResult
}
