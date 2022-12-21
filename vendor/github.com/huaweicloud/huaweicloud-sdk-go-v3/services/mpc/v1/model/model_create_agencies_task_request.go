package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateAgenciesTaskRequest struct {
	Body *AgenciesTaskReq `json:"body,omitempty"`
}

func (o CreateAgenciesTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAgenciesTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateAgenciesTaskRequest", string(data)}, " ")
}
