package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowStorageUsedSpaceResponse Response Object
type ShowStorageUsedSpaceResponse struct {

	// 节点id。
	NodeId *string `json:"node_id,omitempty"`

	// 磁盘空间使用量。
	Used           *string `json:"used,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowStorageUsedSpaceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowStorageUsedSpaceResponse struct{}"
	}

	return strings.Join([]string{"ShowStorageUsedSpaceResponse", string(data)}, " ")
}
