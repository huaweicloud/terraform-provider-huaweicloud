package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteTranscodeTemplateRequest struct {

	// 模板id
	GroupId string `json:"group_id"`
}

func (o DeleteTranscodeTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTranscodeTemplateRequest struct{}"
	}

	return strings.Join([]string{"DeleteTranscodeTemplateRequest", string(data)}, " ")
}
