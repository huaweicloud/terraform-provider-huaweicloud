package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateInstanceByEngineResponse Response Object
type CreateInstanceByEngineResponse struct {

	// 实例ID
	InstanceId     *string `json:"instance_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateInstanceByEngineResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateInstanceByEngineResponse struct{}"
	}

	return strings.Join([]string{"CreateInstanceByEngineResponse", string(data)}, " ")
}
