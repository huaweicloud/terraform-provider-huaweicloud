package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ConfirmImageUploadResponse struct {

	// 水印配置模板id<br/>
	Id *string `json:"id,omitempty"`

	// 水印图片的下载url<br/>
	ImageUrl       *string `json:"image_url,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ConfirmImageUploadResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfirmImageUploadResponse struct{}"
	}

	return strings.Join([]string{"ConfirmImageUploadResponse", string(data)}, " ")
}
