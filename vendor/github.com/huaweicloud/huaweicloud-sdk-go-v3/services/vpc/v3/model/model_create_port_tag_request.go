package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePortTagRequest Request Object
type CreatePortTagRequest struct {

	// 功能说明：端口唯一标识 取值范围：合法UUID 约束：ID对应的端口必须存在
	PortId string `json:"port_id"`

	Body *CreatePortTagRequestBody `json:"body,omitempty"`
}

func (o CreatePortTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePortTagRequest struct{}"
	}

	return strings.Join([]string{"CreatePortTagRequest", string(data)}, " ")
}
