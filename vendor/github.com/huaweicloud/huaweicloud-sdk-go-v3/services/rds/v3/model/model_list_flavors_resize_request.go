package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListFlavorsResizeRequest Request Object
type ListFlavorsResizeRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ListFlavorsResizeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFlavorsResizeRequest struct{}"
	}

	return strings.Join([]string{"ListFlavorsResizeRequest", string(data)}, " ")
}
