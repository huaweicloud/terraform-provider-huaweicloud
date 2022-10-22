package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type StartPipelineResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o StartPipelineResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartPipelineResponse struct{}"
	}

	return strings.Join([]string{"StartPipelineResponse", string(data)}, " ")
}
