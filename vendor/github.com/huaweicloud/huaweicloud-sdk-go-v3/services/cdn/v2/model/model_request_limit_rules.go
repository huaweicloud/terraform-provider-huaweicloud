package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RequestLimitRules 请求限速配置。
type RequestLimitRules struct {

	// status只支持on，off无效。  > request_limit_rules字段置空时代表关闭请求限速功能。  > 旧接口中的参数，后续将下线。
	Status *string `json:"status,omitempty"`

	// 优先级，值越大，优先级越高,取值范围：1-100。
	Priority int32 `json:"priority"`

	// 匹配类型，all：所有文件，catalog：目录。
	MatchType string `json:"match_type"`

	// 匹配类型值。 当match_type为all时传空值，例如：\"\"； 当match_type为catalog时传目录地址，以“/”作为首字符,例如：\"/test\"。  > 值为catalog的时候必填
	MatchValue *string `json:"match_value,omitempty"`

	// 限速方式，当前仅支持按流量大小限速，取值为size。
	Type string `json:"type"`

	// 限速条件,type=size,limit_rate_after=50表示从传输表示传输50个字节后开始限速且限速值为limit_rate_value， 单位byte，取值范围：0-1073741824。
	LimitRateAfter int64 `json:"limit_rate_after"`

	// 限速值,单位Bps，取值范围 0-104857600。
	LimitRateValue int32 `json:"limit_rate_value"`
}

func (o RequestLimitRules) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RequestLimitRules struct{}"
	}

	return strings.Join([]string{"RequestLimitRules", string(data)}, " ")
}
