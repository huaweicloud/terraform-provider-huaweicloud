package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBucketRecord 查询桶返回的record信息
type ShowBucketRecord struct {

	// 对象名
	Name *string `json:"name,omitempty"`

	// 对象大小，若对象无size属性，则返回--
	Size *string `json:"size,omitempty"`
}

func (o ShowBucketRecord) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBucketRecord struct{}"
	}

	return strings.Join([]string{"ShowBucketRecord", string(data)}, " ")
}
