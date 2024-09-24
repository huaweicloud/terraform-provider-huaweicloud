package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePortTagRequest Request Object
type DeletePortTagRequest struct {

	// 功能说明：端口唯一标识 取值范围：合法UUID 约束：ID对应的端口必须存在
	PortId string `json:"port_id"`

	// 功能说明：标签键
	TagKey string `json:"tag_key"`
}

func (o DeletePortTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePortTagRequest struct{}"
	}

	return strings.Join([]string{"DeletePortTagRequest", string(data)}, " ")
}
