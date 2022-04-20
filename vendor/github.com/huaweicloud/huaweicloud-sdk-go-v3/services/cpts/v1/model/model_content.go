package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Content struct {
	// content_type

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
