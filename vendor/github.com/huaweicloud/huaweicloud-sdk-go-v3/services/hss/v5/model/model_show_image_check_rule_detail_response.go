package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowImageCheckRuleDetailResponse Response Object
type ShowImageCheckRuleDetailResponse struct {

	// 检查项描述
	Description *string `json:"description,omitempty"`

	// 参考
	Reference *string `json:"reference,omitempty"`

	// 审计描述
	Audit *string `json:"audit,omitempty"`

	// 修改建议
	Remediation *string `json:"remediation,omitempty"`

	// 检测用例信息
	CheckInfoList  *[]ImageCheckRuleCheckCaseResponseInfo `json:"check_info_list,omitempty"`
	HttpStatusCode int                                    `json:"-"`
}

func (o ShowImageCheckRuleDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowImageCheckRuleDetailResponse struct{}"
	}

	return strings.Join([]string{"ShowImageCheckRuleDetailResponse", string(data)}, " ")
}
