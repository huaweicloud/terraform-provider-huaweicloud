package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// IsolationStatus 隔离状态，包含如下:   - isolated : 已隔离   - restored : 已恢复   - isolating : 已下发隔离任务   - restoring : 已下发恢复任务
type IsolationStatus struct {
}

func (o IsolationStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IsolationStatus struct{}"
	}

	return strings.Join([]string{"IsolationStatus", string(data)}, " ")
}
