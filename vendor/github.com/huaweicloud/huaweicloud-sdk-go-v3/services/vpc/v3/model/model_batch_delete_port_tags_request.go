package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeletePortTagsRequest Request Object
type BatchDeletePortTagsRequest struct {

	// 功能说明：端口唯一标识 取值范围：合法UUID 约束：ID对应的端口必须存在
	PortId string `json:"port_id"`

	Body *BatchDeletePortTagsRequestBody `json:"body,omitempty"`
}

func (o BatchDeletePortTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeletePortTagsRequest struct{}"
	}

	return strings.Join([]string{"BatchDeletePortTagsRequest", string(data)}, " ")
}
