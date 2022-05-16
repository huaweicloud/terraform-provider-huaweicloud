package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateProjectStatusRequest struct {

	// 待设置状态的项目ID，获取方式请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	ProjectId string `json:"project_id"`

	Body *UpdateProjectStatusRequestBody `json:"body,omitempty"`
}

func (o UpdateProjectStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProjectStatusRequest struct{}"
	}

	return strings.Join([]string{"UpdateProjectStatusRequest", string(data)}, " ")
}
