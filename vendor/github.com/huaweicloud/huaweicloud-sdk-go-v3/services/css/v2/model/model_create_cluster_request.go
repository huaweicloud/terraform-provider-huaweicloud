package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClusterRequest Request Object
type CreateClusterRequest struct {
	Body *CreateClusterReq `json:"body,omitempty"`
}

func (o CreateClusterRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterRequest struct{}"
	}

	return strings.Join([]string{"CreateClusterRequest", string(data)}, " ")
}
