package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ContentHeader struct {

	// 键
	Key *string `json:"key,omitempty"`

	// 值
	Value *string `json:"value,omitempty"`
}

func (o ContentHeader) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ContentHeader struct{}"
	}

	return strings.Join([]string{"ContentHeader", string(data)}, " ")
}
