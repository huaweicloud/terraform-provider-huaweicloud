package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListQueryExportTaskRequest Request Object
type ListQueryExportTaskRequest struct {

	// 任务id
	TaskId string `json:"task_id"`

	// Region Id
	Region string `json:"region"`

	// 企业项目ID
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ListQueryExportTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListQueryExportTaskRequest struct{}"
	}

	return strings.Join([]string{"ListQueryExportTaskRequest", string(data)}, " ")
}
