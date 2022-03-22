package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowGaussMySqlFlavorsRequest struct {
	// 语言。

	XLanguage *string `json:"X-Language,omitempty"`
	// 数据库引擎名称。

	DatabaseName string `json:"database_name"`
	// 数据库版本号，目前仅支持兼容MySQL 8.0。

	VersionName *string `json:"version_name,omitempty"`
	// 规格的可用区模式，现在仅支持\"single\"、\"multi\"，不区分大小写。

	AvailabilityZoneMode string `json:"availability_zone_mode"`
	// 规格编码。

	SpecCode *string `json:"spec_code,omitempty"`
}

func (o ShowGaussMySqlFlavorsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlFlavorsRequest struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlFlavorsRequest", string(data)}, " ")
}
