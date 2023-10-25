package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopHotPipelineResponse Response Object
type StopHotPipelineResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StopHotPipelineResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopHotPipelineResponse struct{}"
	}

	return strings.Join([]string{"StopHotPipelineResponse", string(data)}, " ")
}
