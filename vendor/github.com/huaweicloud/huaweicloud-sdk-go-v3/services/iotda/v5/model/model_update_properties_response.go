package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePropertiesResponse Response Object
type UpdatePropertiesResponse struct {

	// 设备属性更新ID，用于唯一标识一条属性更新，在下发更新属性时由物联网平台分配获得。
	RequestId *string `json:"request_id,omitempty"`

	// 设备上报的属性执行结果。Json格式，具体格式需要应用和设备约定。
	Response *interface{} `json:"response,omitempty"`

	// 属性更新异常错误码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 属性更新异常错误信息。
	ErrorMsg       *interface{} `json:"error_msg,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o UpdatePropertiesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePropertiesResponse struct{}"
	}

	return strings.Join([]string{"UpdatePropertiesResponse", string(data)}, " ")
}
