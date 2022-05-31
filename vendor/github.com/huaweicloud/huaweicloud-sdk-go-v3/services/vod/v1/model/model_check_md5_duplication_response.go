package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CheckMd5DuplicationResponse struct {

	// 是否重复。  取值如下： - 0：表示不重复。 - 1：表示重复。
	IsDuplicated *int32 `json:"is_duplicated,omitempty"`

	// 重复的媒资ID
	AssetIds       *[]string `json:"asset_ids,omitempty"`
	HttpStatusCode int       `json:"-"`
}

func (o CheckMd5DuplicationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckMd5DuplicationResponse struct{}"
	}

	return strings.Join([]string{"CheckMd5DuplicationResponse", string(data)}, " ")
}
