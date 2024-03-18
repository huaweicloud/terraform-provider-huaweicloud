package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyRdSforMySqlProxyRouteModeResponse Response Object
type ModifyRdSforMySqlProxyRouteModeResponse struct {

	// 修改数据库代理路由模式结果。 取值：     failed 失败     success 成功
	Result         *string `json:"result,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ModifyRdSforMySqlProxyRouteModeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyRdSforMySqlProxyRouteModeResponse struct{}"
	}

	return strings.Join([]string{"ModifyRdSforMySqlProxyRouteModeResponse", string(data)}, " ")
}
