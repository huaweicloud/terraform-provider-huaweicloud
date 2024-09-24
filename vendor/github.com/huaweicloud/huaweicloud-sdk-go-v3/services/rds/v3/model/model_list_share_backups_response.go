package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListShareBackupsResponse Response Object
type ListShareBackupsResponse struct {

	// 共享备份列表。
	Backups *[]ShareBackups `json:"backups,omitempty"`

	// 总记录数。
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListShareBackupsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListShareBackupsResponse struct{}"
	}

	return strings.Join([]string{"ListShareBackupsResponse", string(data)}, " ")
}
