package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreatePortTagsRequest Request Object
type BatchCreatePortTagsRequest struct {

	// 功能说明：端口唯一标识 取值范围：合法UUID 约束：ID对应的端口必须存在
	PortId string `json:"port_id"`

	Body *BatchCreatePortTagsRequestBody `json:"body,omitempty"`
}

func (o BatchCreatePortTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreatePortTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchCreatePortTagsRequest", string(data)}, " ")
}
