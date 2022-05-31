package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateCommandResponse struct {

	// 设备命令ID，用于唯一标识一条命令，在下发设备命令时由物联网平台分配获得。
	CommandId *string `json:"command_id,omitempty"`

	// 设备上报的命令执行结果。Json格式，具体格式需要应用和设备约定。
	Response       *interface{} `json:"response,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o CreateCommandResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCommandResponse struct{}"
	}

	return strings.Join([]string{"CreateCommandResponse", string(data)}, " ")
}
