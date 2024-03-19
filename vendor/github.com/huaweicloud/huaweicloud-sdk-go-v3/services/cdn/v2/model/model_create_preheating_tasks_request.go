package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePreheatingTasksRequest Request Object
type CreatePreheatingTasksRequest struct {

	// 当用户开启企业项目功能时，该参数生效，表示在当前企业项目下添加缓存预热任务，\"all\"代表所有项目。注意：当使用子帐号调用接口时，该参数必传。  您可以通过调用企业项目管理服务（EPS）的查询企业项目列表接口（ListEnterpriseProject）查询企业项目id。
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
