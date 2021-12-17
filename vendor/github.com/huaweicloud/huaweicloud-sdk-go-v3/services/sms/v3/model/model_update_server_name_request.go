package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateServerNameRequest struct {
	// 源端服务器在主机迁移服务中的id

	SourceId string `json:"source_id"`

	Body *PutSourceServerBody `json:"body,omitempty"`
}

func (o UpdateServerNameRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateServerNameRequest struct{}"
	}

	return strings.Join([]string{"UpdateServerNameRequest", string(data)}, " ")
}
