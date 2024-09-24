package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SyncNodePoolResponse Response Object
type SyncNodePoolResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SyncNodePoolResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SyncNodePoolResponse struct{}"
	}

	return strings.Join([]string{"SyncNodePoolResponse", string(data)}, " ")
}
