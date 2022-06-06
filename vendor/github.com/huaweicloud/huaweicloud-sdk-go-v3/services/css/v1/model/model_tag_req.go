package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 标签对象列表。
type TagReq struct {
	Tag *Tag `json:"tag"`
}

func (o TagReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagReq struct{}"
	}

	return strings.Join([]string{"TagReq", string(data)}, " ")
}
