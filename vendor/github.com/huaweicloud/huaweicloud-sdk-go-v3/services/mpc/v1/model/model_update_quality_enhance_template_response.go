package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateQualityEnhanceTemplateResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateQualityEnhanceTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateQualityEnhanceTemplateResponse struct{}"
	}

	return strings.Join([]string{"UpdateQualityEnhanceTemplateResponse", string(data)}, " ")
}
