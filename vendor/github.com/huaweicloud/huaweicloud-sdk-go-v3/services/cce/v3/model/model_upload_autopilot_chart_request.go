package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UploadAutopilotChartRequest Request Object
type UploadAutopilotChartRequest struct {
	Body *UploadAutopilotChartRequestBody `json:"body,omitempty" type:"multipart"`
}

func (o UploadAutopilotChartRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UploadAutopilotChartRequest struct{}"
	}

	return strings.Join([]string{"UploadAutopilotChartRequest", string(data)}, " ")
}
