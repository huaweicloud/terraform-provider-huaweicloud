package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteTempRequest struct {
	// 事务id

	TemplateId int32 `json:"template_id"`
}

func (o DeleteTempRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTempRequest struct{}"
	}

	return strings.Join([]string{"DeleteTempRequest", string(data)}, " ")
}
