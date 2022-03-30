package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowGaussMySqlEngineVersionResponse struct {
	// 数据库版本信息列表

	Datastores     *[]MysqlEngineVersionInfo `json:"datastores,omitempty"`
	HttpStatusCode int                       `json:"-"`
}

func (o ShowGaussMySqlEngineVersionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlEngineVersionResponse struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlEngineVersionResponse", string(data)}, " ")
}
