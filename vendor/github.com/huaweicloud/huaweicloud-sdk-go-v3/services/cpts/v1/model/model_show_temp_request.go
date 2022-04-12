package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowTempRequest struct {
	// 事务id

	TemplateId int32 `json:"template_id"`
}

func (o ShowTempRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTempRequest struct{}"
	}

	return strings.Join([]string{"ShowTempRequest", string(data)}, " ")
}
