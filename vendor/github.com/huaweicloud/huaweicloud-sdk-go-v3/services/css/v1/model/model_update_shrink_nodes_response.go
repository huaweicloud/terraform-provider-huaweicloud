package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateShrinkNodesResponse Response Object
type UpdateShrinkNodesResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateShrinkNodesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateShrinkNodesResponse struct{}"
	}

	return strings.Join([]string{"UpdateShrinkNodesResponse", string(data)}, " ")
}
