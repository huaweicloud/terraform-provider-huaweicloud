package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 威胁等级，包含如下:   - Security : 安全   - Low : 低危   - Medium : 中危   - High : 高危   - Critical : 危急
type Severity struct {
}

func (o Severity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Severity struct{}"
	}

	return strings.Join([]string{"Severity", string(data)}, " ")
}
