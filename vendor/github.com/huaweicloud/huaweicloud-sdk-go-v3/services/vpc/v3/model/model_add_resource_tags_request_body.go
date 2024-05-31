package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddResourceTagsRequestBody struct {
	Tag *AddResourceTagsRequestBodyTag `json:"tag"`
}

func (o AddResourceTagsRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddResourceTagsRequestBody struct{}"
	}

	return strings.Join([]string{"AddResourceTagsRequestBody", string(data)}, " ")
}
