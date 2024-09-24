package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProductUrlInfo struct {

	// 删除的产物URL
	Url *string `json:"url,omitempty"`
}

func (o ProductUrlInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProductUrlInfo struct{}"
	}

	return strings.Join([]string{"ProductUrlInfo", string(data)}, " ")
}
