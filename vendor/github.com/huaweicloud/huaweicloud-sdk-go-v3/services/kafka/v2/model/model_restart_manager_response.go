package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type RestartManagerResponse struct {

	// 执行结果。
	Result *string `json:"result,omitempty"`

	// 实例ID。
	InstanceId     *string `json:"instance_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o RestartManagerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestartManagerResponse struct{}"
	}

	return strings.Join([]string{"RestartManagerResponse", string(data)}, " ")
}
