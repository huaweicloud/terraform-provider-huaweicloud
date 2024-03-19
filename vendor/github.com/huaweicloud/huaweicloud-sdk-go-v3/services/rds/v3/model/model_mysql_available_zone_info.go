package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// MysqlAvailableZoneInfo 可用区信息。
type MysqlAvailableZoneInfo struct {

	// 可用区编码。
	Code *string `json:"code,omitempty"`

	// 可用区描述。
	Description *string `json:"description,omitempty"`
}

func (o MysqlAvailableZoneInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlAvailableZoneInfo struct{}"
	}

	return strings.Join([]string{"MysqlAvailableZoneInfo", string(data)}, " ")
}
