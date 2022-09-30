package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeleteTemplateReq struct {

	// 模板名称。
	Name string `json:"name"`
}

func (o DeleteTemplateReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTemplateReq struct{}"
	}

	return strings.Join([]string{"DeleteTemplateReq", string(data)}, " ")
}
