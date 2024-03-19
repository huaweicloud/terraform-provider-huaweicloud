package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResumeConnectorTaskResponse Response Object
type ResumeConnectorTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ResumeConnectorTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResumeConnectorTaskResponse struct{}"
	}

	return strings.Join([]string{"ResumeConnectorTaskResponse", string(data)}, " ")
}
