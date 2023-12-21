package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePromInstanceRequest Request Object
type CreatePromInstanceRequest struct {
	Body *PromInstanceEpsModel `json:"body,omitempty"`
}

func (o CreatePromInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePromInstanceRequest struct{}"
	}

	return strings.Join([]string{"CreatePromInstanceRequest", string(data)}, " ")
}
