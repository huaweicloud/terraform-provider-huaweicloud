package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SyncNodeResponse Response Object
type SyncNodeResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SyncNodeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SyncNodeResponse struct{}"
	}

	return strings.Join([]string{"SyncNodeResponse", string(data)}, " ")
}
