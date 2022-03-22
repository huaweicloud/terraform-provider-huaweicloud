package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowGaussMySqlEngineVersionRequest struct {
	// 语言。

	XLanguage *string `json:"X-Language,omitempty"`
	// 数据库引擎。支持的引擎如下，不区分大小写：gaussdb-mysql

	DatabaseName string `json:"database_name"`
}

func (o ShowGaussMySqlEngineVersionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlEngineVersionRequest struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlEngineVersionRequest", string(data)}, " ")
}
