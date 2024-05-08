package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// IsolateSource 隔离来源，包含如下:   - event : 安全告警事件   - antivirus : 病毒查杀
type IsolateSource struct {
}

func (o IsolateSource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IsolateSource struct{}"
	}

	return strings.Join([]string{"IsolateSource", string(data)}, " ")
}
