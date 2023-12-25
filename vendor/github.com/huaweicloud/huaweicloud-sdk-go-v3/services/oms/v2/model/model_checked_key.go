package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CheckedKey 检查键
type CheckedKey struct {

	// 键
	Key *string `json:"key,omitempty"`

	// 是否电子标签匹配
	IsEtagMatching *bool `json:"is_etag_matching,omitempty"`

	// 是否存在对象
	IsObjectExisting *bool `json:"is_object_existing,omitempty"`
}

func (o CheckedKey) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckedKey struct{}"
	}

	return strings.Join([]string{"CheckedKey", string(data)}, " ")
}
