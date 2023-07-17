package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ConnectorOrderRequestBody struct {

	// 需要关闭connector的实例id，和请求路径上的一致。
	InstanceId string `json:"instance_id"`

	// 提交关闭connector订单后前端跳转的页面
	Url *string `json:"url,omitempty"`
}

func (o ConnectorOrderRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConnectorOrderRequestBody struct{}"
	}

	return strings.Join([]string{"ConnectorOrderRequestBody", string(data)}, " ")
}
