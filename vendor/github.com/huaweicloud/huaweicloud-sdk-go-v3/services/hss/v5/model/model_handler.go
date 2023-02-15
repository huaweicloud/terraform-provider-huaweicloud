package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 手动处理的备注
type Handler struct {
}

func (o Handler) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Handler struct{}"
	}

	return strings.Join([]string{"Handler", string(data)}, " ")
}
