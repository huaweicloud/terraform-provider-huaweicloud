package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateCommandResultRequest struct {
	// 上报命令执行结果的命令所对应的服务端id

	ServerId string `json:"server_id"`

	Body *CommandBody `json:"body,omitempty"`
}

func (o UpdateCommandResultRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCommandResultRequest struct{}"
	}

	return strings.Join([]string{"UpdateCommandResultRequest", string(data)}, " ")
}
