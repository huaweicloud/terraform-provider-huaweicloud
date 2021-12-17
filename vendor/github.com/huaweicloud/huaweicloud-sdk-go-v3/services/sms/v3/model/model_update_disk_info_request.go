package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateDiskInfoRequest struct {
	// 源端服务器id

	SourceId string `json:"source_id"`

	Body *PutDiskInfoReq `json:"body,omitempty"`
}

func (o UpdateDiskInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDiskInfoRequest struct{}"
	}

	return strings.Join([]string{"UpdateDiskInfoRequest", string(data)}, " ")
}
