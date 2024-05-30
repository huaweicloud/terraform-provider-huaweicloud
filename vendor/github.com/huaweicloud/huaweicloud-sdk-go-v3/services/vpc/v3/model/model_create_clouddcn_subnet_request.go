package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClouddcnSubnetRequest Request Object
type CreateClouddcnSubnetRequest struct {
	Body *CreateClouddcnSubnetRequestBody `json:"body,omitempty"`
}

func (o CreateClouddcnSubnetRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClouddcnSubnetRequest struct{}"
	}

	return strings.Join([]string{"CreateClouddcnSubnetRequest", string(data)}, " ")
}
