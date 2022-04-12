package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CodeMessageResq struct {
	// code

	Code *string `json:"code,omitempty"`
	// message

	Message *string `json:"message,omitempty"`
}

func (o CodeMessageResq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CodeMessageResq struct{}"
	}

	return strings.Join([]string{"CodeMessageResq", string(data)}, " ")
}
