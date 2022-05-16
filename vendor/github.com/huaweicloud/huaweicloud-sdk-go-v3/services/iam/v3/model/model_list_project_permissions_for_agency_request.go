package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListProjectPermissionsForAgencyRequest struct {

	// 委托方的项目ID，获取方式请参见：[获取项目名称、项目ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	ProjectId string `json:"project_id"`

	// 委托ID，获取方式请参见：[获取委托名、委托ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	AgencyId string `json:"agency_id"`
}

func (o ListProjectPermissionsForAgencyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProjectPermissionsForAgencyRequest struct{}"
	}

	return strings.Join([]string{"ListProjectPermissionsForAgencyRequest", string(data)}, " ")
}
