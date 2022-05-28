package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateWatermarkTemplateResponse struct {

	// 水印配置模板id<br/>
	Id *string `json:"id,omitempty"`

	// 水印图片上传地址<br/>
	UploadUrl      *string `json:"upload_url,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateWatermarkTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateWatermarkTemplateResponse struct{}"
	}

	return strings.Join([]string{"CreateWatermarkTemplateResponse", string(data)}, " ")
}
