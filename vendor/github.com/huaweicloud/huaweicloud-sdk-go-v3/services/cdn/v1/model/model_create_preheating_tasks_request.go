package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreatePreheatingTasksRequest struct {

	// 当用户开启企业项目功能时，该参数生效，表示在当前企业项目下添加缓存预热任务，\"all\"代表所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	Body *PreheatingTaskRequest `json:"body,omitempty"`
}

func (o CreatePreheatingTasksRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePreheatingTasksRequest struct{}"
	}

	return strings.Join([]string{"CreatePreheatingTasksRequest", string(data)}, " ")
}
