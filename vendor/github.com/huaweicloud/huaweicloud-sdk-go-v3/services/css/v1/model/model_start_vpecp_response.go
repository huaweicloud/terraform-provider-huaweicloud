package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartVpecpResponse Response Object
type StartVpecpResponse struct {

	// 操作行为。固定为：createVpcepservice，表示已开启终端节点。
	Action         *string `json:"action,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartVpecpResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartVpecpResponse struct{}"
	}

	return strings.Join([]string{"StartVpecpResponse", string(data)}, " ")
}
