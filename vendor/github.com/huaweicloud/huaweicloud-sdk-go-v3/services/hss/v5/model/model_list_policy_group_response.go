package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPolicyGroupResponse Response Object
type ListPolicyGroupResponse struct {

	// 总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 策略组列表
	DataList       *[]PolicyGroupResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ListPolicyGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPolicyGroupResponse struct{}"
	}

	return strings.Join([]string{"ListPolicyGroupResponse", string(data)}, " ")
}
