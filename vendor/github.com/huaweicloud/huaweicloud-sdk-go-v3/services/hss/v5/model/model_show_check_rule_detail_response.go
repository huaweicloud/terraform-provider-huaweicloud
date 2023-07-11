package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCheckRuleDetailResponse Response Object
type ShowCheckRuleDetailResponse struct {

	// 描述
	Description *string `json:"description,omitempty"`

	// 根据
	Reference *string `json:"reference,omitempty"`

	// 审计描述
	Audit *string `json:"audit,omitempty"`

	// 修改建议
	Remediation *string `json:"remediation,omitempty"`

	// 检测用例信息
	CheckInfoList  *[]CheckRuleCheckCaseResponseInfo `json:"check_info_list,omitempty"`
	HttpStatusCode int                               `json:"-"`
}

func (o ShowCheckRuleDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCheckRuleDetailResponse struct{}"
	}

	return strings.Join([]string{"ShowCheckRuleDetailResponse", string(data)}, " ")
}
