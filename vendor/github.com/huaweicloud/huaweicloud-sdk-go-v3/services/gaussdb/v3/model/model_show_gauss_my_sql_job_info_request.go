package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowGaussMySqlJobInfoRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 任务ID。

	Id string `json:"id"`
}

func (o ShowGaussMySqlJobInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlJobInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlJobInfoRequest", string(data)}, " ")
}
