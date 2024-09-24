package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ScaleNodePoolResponse Response Object
type ScaleNodePoolResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ScaleNodePoolResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleNodePoolResponse struct{}"
	}

	return strings.Join([]string{"ScaleNodePoolResponse", string(data)}, " ")
}
