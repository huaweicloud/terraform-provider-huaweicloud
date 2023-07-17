package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePostPaidInstanceResponse Response Object
type CreatePostPaidInstanceResponse struct {

	// 实例ID
	InstanceId     *string `json:"instance_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreatePostPaidInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePostPaidInstanceResponse struct{}"
	}

	return strings.Join([]string{"CreatePostPaidInstanceResponse", string(data)}, " ")
}
