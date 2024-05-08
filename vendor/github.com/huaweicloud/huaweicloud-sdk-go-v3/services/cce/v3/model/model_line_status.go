package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type LineStatus struct {
	StartPoint *Point `json:"startPoint,omitempty"`

	EndPoint *Point `json:"endPoint,omitempty"`

	// 表示是否为关键线路（关键线路未执行无法取消升级流程）
	Critical *string `json:"critical,omitempty"`
}

func (o LineStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LineStatus struct{}"
	}

	return strings.Join([]string{"LineStatus", string(data)}, " ")
}
