package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOrganizationTreeResponse Response Object
type ListOrganizationTreeResponse struct {

	// 总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 事件列表详情
	DataList *[]OrganizationNodeResponseInfo `json:"data_list,omitempty"`

	XRequestId     *string `json:"X-request-id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListOrganizationTreeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOrganizationTreeResponse struct{}"
	}

	return strings.Join([]string{"ListOrganizationTreeResponse", string(data)}, " ")
}
