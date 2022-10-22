package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type BatchRestartOrDeleteInstancesResponse struct {

	// 修改实例的结果。
	Results        *[]BatchRestartOrDeleteInstanceRespResults `json:"results,omitempty"`
	HttpStatusCode int                                        `json:"-"`
}

func (o BatchRestartOrDeleteInstancesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchRestartOrDeleteInstancesResponse struct{}"
	}

	return strings.Join([]string{"BatchRestartOrDeleteInstancesResponse", string(data)}, " ")
}
