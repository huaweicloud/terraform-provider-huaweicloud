package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Description 描述信息
type Description struct {
}

func (o Description) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Description struct{}"
	}

	return strings.Join([]string{"Description", string(data)}, " ")
}
