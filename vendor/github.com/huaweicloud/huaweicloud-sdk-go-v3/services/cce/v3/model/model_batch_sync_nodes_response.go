package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchSyncNodesResponse Response Object
type BatchSyncNodesResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o BatchSyncNodesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchSyncNodesResponse struct{}"
	}

	return strings.Join([]string{"BatchSyncNodesResponse", string(data)}, " ")
}
