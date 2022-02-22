package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateSubNetworkInterfaceRequest struct {
	Body *CreateSubNetworkInterfaceRequestBody `json:"body,omitempty"`
}

func (o CreateSubNetworkInterfaceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSubNetworkInterfaceRequest struct{}"
	}

	return strings.Join([]string{"CreateSubNetworkInterfaceRequest", string(data)}, " ")
}
