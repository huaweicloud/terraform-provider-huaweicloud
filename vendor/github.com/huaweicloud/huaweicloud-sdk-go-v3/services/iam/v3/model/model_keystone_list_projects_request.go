package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListProjectsRequest struct {

	// 项目所属账号ID，获取方式请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId *string `json:"domain_id,omitempty"`

	// 项目名称，获取方式请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	Name *string `json:"name,omitempty"`

	// 如果查询自己创建的项目，则此处应填为所属区域的项目ID。  如果查询的是系统内置项目，如cn-north-4，则此处应填为账号ID。  获取项目ID和账号ID，请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	ParentId *string `json:"parent_id,omitempty"`

	// 项目是否启用。
	Enabled *bool `json:"enabled,omitempty"`

	// 该字段无需填写。
	IsDomain *bool `json:"is_domain,omitempty"`

	// 分页查询时数据的页数，查询值最小为1。需要与per_page同时存在。
	Page *int32 `json:"page,omitempty"`

	// 分页查询时每页的数据个数，取值范围为[1,5000]。需要与page同时存在。
	PerPage *int32 `json:"per_page,omitempty"`
}

func (o KeystoneListProjectsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListProjectsRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListProjectsRequest", string(data)}, " ")
}
