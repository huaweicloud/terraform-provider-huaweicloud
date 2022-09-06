package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListTranscodeTemplateResponse struct {

	// 模板组信息<br/>
	TemplateGroupList *[]TransTemplateRsp `json:"template_group_list,omitempty"`

	// 总记录条数<br/>
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListTranscodeTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTranscodeTemplateResponse struct{}"
	}

	return strings.Join([]string{"ListTranscodeTemplateResponse", string(data)}, " ")
}
