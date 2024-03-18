package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Content struct {

	// 用例卡片类型（1：思考时间；2：报文；3：检查点；4：变量提取）
	ContentType *int32 `json:"content_type,omitempty"`

	Content *ContentInfo `json:"content,omitempty"`
}

func (o Content) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Content struct{}"
	}

	return strings.Join([]string{"Content", string(data)}, " ")
}
