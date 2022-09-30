package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTemplatesRequest struct {

	// 模板类型。custom为自定义模板，system为系统模板。不指定查询模板类型默认查找自定义模板和系统模板。
	Type *string `json:"type,omitempty"`
}

func (o ListTemplatesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTemplatesRequest struct{}"
	}

	return strings.Join([]string{"ListTemplatesRequest", string(data)}, " ")
}
