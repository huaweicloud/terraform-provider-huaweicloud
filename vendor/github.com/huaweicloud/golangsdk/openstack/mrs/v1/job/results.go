package job

import "github.com/huaweicloud/golangsdk"

type Job struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenant_id"`
	JobID          string `json:"job_id"`
	JobName        string `json:"job_name"`
	StartTime      int    `json:"start_time"`
	EndTime        int    `json:"end_time"`
	ClusterID      string `json:"cluster_id"`
	GroupID        string `json:"group_id"`
	JarPath        string `json:"jar_path"`
	Input          string `json:"input"`
	Output         string `json:"output"`
	JobLog         string `json:"job_log"`
	JobType        int    `json:"job_type"`
	FileAction     string `json:"file_action"`
	Arguments      string `json:"arguments"`
	Hql            string `json:"hql"`
	JobState       int    `json:"job_state"`
	JobFinalStatus int    `json:"job_final_status"`
	HiveScriptPath string `json:"hive_script_path"`
	CreateBy       string `json:"create_by"`
	FinishedStep   int    `json:"finished_step"`
	JobMainID      string `json:"job_main_id"`
	JobStepID      string `json:"job_step_id"`
	PostponeAt     int    `json:"postpone_at"`
	StepName       string `json:"step_name"`
	StepNum        int    `json:"step_num"`
	TaskNum        int    `json:"task_num"`
	UpdateBy       string `json:"update_by"`
	SpendTime      int    `json:"spend_time"`
	StepSeq        int    `json:"step_seq"`
	Progress       string `json:"progress"`
}

type JobResult struct {
	Templated      bool   `json:"templated"`
	ID             string `json:"id"`
	TenantID       string `json:"tenant_id"`
	JobID          string `json:"job_id"`
	JobName        string `json:"job_name"`
	InputID        string `json:"input_id"`
	OutputID       string `json:"output_id"`
	StartTime      int    `json:"start_time"`
	EndTime        int    `json:"end_time"`
	ClusterID      string `json:"cluster_id"`
	EngineJobID    string `json:"engine_job_id"`
	ReturnCode     int    `json:"return_code"`
	IsPublic       bool   `json:"is_public"`
	IsProtected    bool   `json:"is_protected"`
	GroupID        string `json:"group_id"`
	JarPath        string `json:"jar_path"`
	Input          string `json:"input"`
	Output         string `json:"output"`
	JobLog         string `json:"job_log"`
	JobType        int    `json:"job_type"`
	FileAction     string `json:"file_action"`
	Arguments      string `json:"arguments"`
	JobState       int    `json:"job_state"`
	JobFinalStatus int    `json:"job_final_status"`
	HiveScriptPath string `json:"hive_script_path"`
	CreateBy       string `json:"create_by"`
	FinishedStep   int    `json:"finished_step"`
	JobMainID      string `json:"job_main_id"`
	JobStepID      string `json:"job_step_id"`
	PostponeAt     int    `json:"postpone_at"`
	StepName       string `json:"step_name"`
	StepNum        int    `json:"step_num"`
	TaskNum        int    `json:"task_num"`
	UpdateBy       string `json:"update_by"`
	Credentials    string `json:"credentials"`
	UserID         string `json:"user_id"`
	JobConfigs     string `json:"job_configs"`
	Extra          string `json:"extra"`
	DataSourceUrls string `json:"data_source_urls"`
	Info           string `json:"info"`
}

type CreateResult struct {
	golangsdk.Result
}

func (r CreateResult) Extract() (*JobResult, error) {
	var s JobResult
	err := r.ExtractInto(&s)
	return &s, err
}

func (r CreateResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "job_execution")
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*Job, error) {
	var s Job
	err := r.ExtractInto(&s)
	return &s, err
}

func (r GetResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "job_execution")
}

type DeleteResult struct {
	golangsdk.ErrResult
}
