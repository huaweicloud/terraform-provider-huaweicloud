package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateInstanceIam5Response Response Object
type CreateInstanceIam5Response struct {
	Instance *CreateInstanceRespItem `json:"instance,omitempty"`

	// 实例创建的任务id。  仅创建按需实例时会返回该参数。
	JobId *string `json:"job_id,omitempty"`

	// 订单号，创建包年包月时返回该参数。
	OrderId        *string `json:"order_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateInstanceIam5Response) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateInstanceIam5Response struct{}"
	}

	return strings.Join([]string{"CreateInstanceIam5Response", string(data)}, " ")
}
