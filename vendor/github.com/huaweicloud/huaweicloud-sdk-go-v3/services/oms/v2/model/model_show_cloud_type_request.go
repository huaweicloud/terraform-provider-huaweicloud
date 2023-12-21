package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCloudTypeRequest Request Object
type ShowCloudTypeRequest struct {

	// 连接端类型源端(src)，目的端(dst)
	Type string `json:"type"`
}

func (o ShowCloudTypeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCloudTypeRequest struct{}"
	}

	return strings.Join([]string{"ShowCloudTypeRequest", string(data)}, " ")
}
