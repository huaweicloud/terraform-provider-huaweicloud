package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateAgencyRequest struct {
	Body *CreateAgencyRequestBody `json:"body,omitempty"`
}

func (o CreateAgencyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAgencyRequest struct{}"
	}

	return strings.Join([]string{"CreateAgencyRequest", string(data)}, " ")
}
