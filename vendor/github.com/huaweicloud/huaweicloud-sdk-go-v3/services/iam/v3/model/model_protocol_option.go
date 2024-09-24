package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ProtocolOption
type ProtocolOption struct {

	// 映射ID。身份提供商类型为iam_user_sso时，不需要绑定映射ID，无需传入此字段；否则此字段必填。
	MappingId string `json:"mapping_id"`
}

func (o ProtocolOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtocolOption struct{}"
	}

	return strings.Join([]string{"ProtocolOption", string(data)}, " ")
}
