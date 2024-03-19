package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetInstancesNewDbShrinkResponse Response Object
type SetInstancesNewDbShrinkResponse struct {

	// 收缩结果。successful:成功 failed:失败
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SetInstancesNewDbShrinkResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetInstancesNewDbShrinkResponse struct{}"
	}

	return strings.Join([]string{"SetInstancesNewDbShrinkResponse", string(data)}, " ")
}
