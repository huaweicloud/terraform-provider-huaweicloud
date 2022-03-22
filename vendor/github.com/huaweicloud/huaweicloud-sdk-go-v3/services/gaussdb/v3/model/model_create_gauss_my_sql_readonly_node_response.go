package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateGaussMySqlReadonlyNodeResponse struct {
	// 实例ID。

	InstanceId *string `json:"instance_id,omitempty"`
	// 节点名称列表。

	NodeNames *[]string `json:"node_names,omitempty"`
	// 实例创建的任务id。  仅创建按需实例时会返回该参数。

	JobId *string `json:"job_id,omitempty"`
	// 订单号，创建包年包月时返回该参数。

	OrderId        *string `json:"order_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateGaussMySqlReadonlyNodeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateGaussMySqlReadonlyNodeResponse struct{}"
	}

	return strings.Join([]string{"CreateGaussMySqlReadonlyNodeResponse", string(data)}, " ")
}
