package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAllObsObjListRequest Request Object
type ListAllObsObjListRequest struct {

	// 桶名
	Bucket string `json:"bucket"`

	// 查询对象前缀
	Prefix *string `json:"prefix,omitempty"`

	// 查询对象文件类型
	Type *string `json:"type,omitempty"`
}

func (o ListAllObsObjListRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAllObsObjListRequest struct{}"
	}

	return strings.Join([]string{"ListAllObsObjListRequest", string(data)}, " ")
}
