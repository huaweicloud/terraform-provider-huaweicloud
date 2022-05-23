package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 资源空间详情结构体。
type ApplicationDto struct {

	// 资源空间ID，唯一标识一个资源空间，由物联网平台在创建资源空间时分配。资源空间对应的是物联网平台原有的应用，在物联网平台的含义与应用一致，只是变更了名称。
	AppId *string `json:"app_id,omitempty"`

	// 资源空间名称。
	AppName *string `json:"app_name,omitempty"`

	// 资源空间创建时间，格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	CreateTime *string `json:"create_time,omitempty"`

	// 是否为默认资源空间
	DefaultApp *bool `json:"default_app,omitempty"`
}

func (o ApplicationDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ApplicationDto struct{}"
	}

	return strings.Join([]string{"ApplicationDto", string(data)}, " ")
}
