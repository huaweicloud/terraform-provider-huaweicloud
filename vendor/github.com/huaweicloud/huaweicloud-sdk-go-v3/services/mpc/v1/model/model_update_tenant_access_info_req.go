package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateTenantAccessInfoReq struct {

	// 是否已开通服务 - false：未开通 - true：已开通
	IsOpen *bool `json:"is_open,omitempty"`

	// 服务协议版本
	AgreementVersion *int32 `json:"agreement_version,omitempty"`
}

func (o UpdateTenantAccessInfoReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTenantAccessInfoReq struct{}"
	}

	return strings.Join([]string{"UpdateTenantAccessInfoReq", string(data)}, " ")
}
