package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListStreamForbiddenRequest struct {

	// op账号需要携带的特定project_id，当使用op账号时该值为所操作租户的project_id
	SpecifyProject *string `json:"specify_project,omitempty"`

	// 推流域名
	Domain string `json:"domain"`

	// 应用名称，不指定则查询domain下所有应用的禁止直播推流信息
	AppName *string `json:"app_name,omitempty"`

	// 流名称
	StreamName *string `json:"stream_name,omitempty"`

	// 分页编号。 默认为0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。  取值范围：1-100。  默认为10。
	Size *int32 `json:"size,omitempty"`
}

func (o ListStreamForbiddenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListStreamForbiddenRequest struct{}"
	}

	return strings.Join([]string{"ListStreamForbiddenRequest", string(data)}, " ")
}
