package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ErrorInfoDto 异常信息
type ErrorInfoDto struct {

	// **参数说明**：异常信息错误码，包含IOTDA.014016和IOTDA.014112。IOTDA.014016表示设备不在线；IOTDA.014112表示设备没有订阅topic。
	ErrorCode *string `json:"error_code,omitempty"`

	// **参数说明**：异常信息说明，包含设备不在线和设备没有订阅topic说明。
	ErrorMsg *string `json:"error_msg,omitempty"`
}

func (o ErrorInfoDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ErrorInfoDto struct{}"
	}

	return strings.Join([]string{"ErrorInfoDto", string(data)}, " ")
}
