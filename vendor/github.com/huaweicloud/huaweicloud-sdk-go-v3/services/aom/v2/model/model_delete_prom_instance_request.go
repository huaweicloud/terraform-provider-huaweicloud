package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePromInstanceRequest Request Object
type DeletePromInstanceRequest struct {

	// Prometheus实例id。
	PromId string `json:"prom_id"`
}

func (o DeletePromInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePromInstanceRequest struct{}"
	}

	return strings.Join([]string{"DeletePromInstanceRequest", string(data)}, " ")
}
