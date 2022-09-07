package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateTranscodeTemplateResponse struct {

	// 模板组ID<br/>
	GroupId        *string `json:"group_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateTranscodeTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTranscodeTemplateResponse struct{}"
	}

	return strings.Join([]string{"CreateTranscodeTemplateResponse", string(data)}, " ")
}
