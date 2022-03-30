package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlChangeSpecificationRequest struct {
	ResizeFlavor *MysqlResizeFlavor `json:"resize_flavor"`
	// 变更包周期实例规格时可指定，表示是否自动从客户的账户中支付。true，为自动支付，默认该方式。false，为手动支付。

	IsAutoPay *string `json:"is_auto_pay,omitempty"`
}

func (o MysqlChangeSpecificationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlChangeSpecificationRequest struct{}"
	}

	return strings.Join([]string{"MysqlChangeSpecificationRequest", string(data)}, " ")
}
