package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 只读节点信息
type MysqlCreateReadonlyNodeRequest struct {
	// 指定创建的只读节点故障倒换优先级。倒换优先级列表个数即为只读节点格式。 故障倒换优先级的取值范围为1~16，数字越小，优先级越大，即故障倒换时，主节点会优先倒换到优先级高的备节点上，优先级相同的备节点选为主节点的概率相同。

	Priorities []int32 `json:"priorities"`
	// 创建包周期时可指定，表示是否自动从客户的账户中支付，此字段不影响自动续订的支付方式。  - true，为自动支付，默认该方式。 - false，为手动支付。

	IsAutoPay *string `json:"is_auto_pay,omitempty"`
}

func (o MysqlCreateReadonlyNodeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlCreateReadonlyNodeRequest struct{}"
	}

	return strings.Join([]string{"MysqlCreateReadonlyNodeRequest", string(data)}, " ")
}
