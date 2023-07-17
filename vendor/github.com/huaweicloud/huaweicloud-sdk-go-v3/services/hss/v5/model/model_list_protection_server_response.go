package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListProtectionServerResponse Response Object
type ListProtectionServerResponse struct {

	// 总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 查询勒索防护服务器列表
	DataList       *[]ProtectionServerInfo `json:"data_list,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ListProtectionServerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProtectionServerResponse struct{}"
	}

	return strings.Join([]string{"ListProtectionServerResponse", string(data)}, " ")
}
