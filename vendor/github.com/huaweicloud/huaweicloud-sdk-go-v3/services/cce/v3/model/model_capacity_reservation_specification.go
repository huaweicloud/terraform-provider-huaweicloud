package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CapacityReservationSpecification 扩展伸缩组容量预留配置
type CapacityReservationSpecification struct {

	// 私有池id，preference为none时忽略该值
	Id *string `json:"id,omitempty"`

	// 私有池容量选项，为 none 时表示不指定容量预留，为 targeted 时表示指定容量预留，此时 id 不能为空
	Preference *string `json:"preference,omitempty"`
}

func (o CapacityReservationSpecification) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CapacityReservationSpecification struct{}"
	}

	return strings.Join([]string{"CapacityReservationSpecification", string(data)}, " ")
}
