package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeprecatedApiRisks 废弃API风险来源
type DeprecatedApiRisks struct {

	// 请求路径，如/apis/policy/v1beta1/podsecuritypolicies
	Url *string `json:"url,omitempty"`

	// 客户端信息
	UserAgent *string `json:"userAgent,omitempty"`
}

func (o DeprecatedApiRisks) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeprecatedApiRisks struct{}"
	}

	return strings.Join([]string{"DeprecatedApiRisks", string(data)}, " ")
}
