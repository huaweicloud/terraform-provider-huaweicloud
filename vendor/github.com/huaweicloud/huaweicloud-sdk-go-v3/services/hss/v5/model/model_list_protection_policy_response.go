package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListProtectionPolicyResponse Response Object
type ListProtectionPolicyResponse struct {

	// 策略总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 查询防护策略列表
	DataList       *[]ProtectionPolicyInfo `json:"data_list,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ListProtectionPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProtectionPolicyResponse struct{}"
	}

	return strings.Join([]string{"ListProtectionPolicyResponse", string(data)}, " ")
}
