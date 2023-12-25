package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetInstancesDbShrinkResponse Response Object
type SetInstancesDbShrinkResponse struct {

	// 收缩结果。successful:成功 failed:失败
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SetInstancesDbShrinkResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetInstancesDbShrinkResponse struct{}"
	}

	return strings.Join([]string{"SetInstancesDbShrinkResponse", string(data)}, " ")
}
