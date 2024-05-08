package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListWeakPasswordUsersResponse Response Object
type ListWeakPasswordUsersResponse struct {

	// 弱口令总数
	TotalNum *int64 `json:"total_num,omitempty"`

	// 弱口令列表
	DataList       *[]WeakPwdListInfoResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                            `json:"-"`
}

func (o ListWeakPasswordUsersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListWeakPasswordUsersResponse struct{}"
	}

	return strings.Join([]string{"ListWeakPasswordUsersResponse", string(data)}, " ")
}
