package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EsListenersResource struct {

	// 监听器ID。
	Id *string `json:"id,omitempty"`
}

func (o EsListenersResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EsListenersResource struct{}"
	}

	return strings.Join([]string{"EsListenersResource", string(data)}, " ")
}
