package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateProductRequest struct {

	// **参数说明**：实例ID。物理多租下各实例的唯一标识，一般华为云租户无需携带该参数，仅在物理多租场景下从管理面访问API时需要携带该参数。
	InstanceId *string `json:"Instance-Id,omitempty"`

	// **参数说明**：产品ID，用于唯一标识一个产品，在物联网平台创建产品后由平台分配获得。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	ProductId string `json:"product_id"`

	Body *UpdateProduct `json:"body,omitempty"`
}

func (o UpdateProductRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProductRequest struct{}"
	}

	return strings.Join([]string{"UpdateProductRequest", string(data)}, " ")
}
