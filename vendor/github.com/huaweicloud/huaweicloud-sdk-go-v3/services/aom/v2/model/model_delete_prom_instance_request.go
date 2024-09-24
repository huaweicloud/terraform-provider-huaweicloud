package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePromInstanceRequest Request Object
type DeletePromInstanceRequest struct {

	// Prometheus实例id。
	PromId string `json:"prom_id"`

	// 企业项目id。 - 删除单个企业项目下实例，填写企业项目id。
	EnterpriseProjectId string `json:"Enterprise-Project-Id"`
}

func (o DeletePromInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePromInstanceRequest struct{}"
	}

	return strings.Join([]string{"DeletePromInstanceRequest", string(data)}, " ")
}
