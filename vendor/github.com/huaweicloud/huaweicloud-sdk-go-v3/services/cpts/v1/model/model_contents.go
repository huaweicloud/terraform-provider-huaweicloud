package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Contents struct {
	// content_id

	ContentId *int32 `json:"content_id,omitempty"`
	// content

	Content *[]Content `json:"content,omitempty"`
	// index

	Index *int32 `json:"index,omitempty"`
	// selected_temp_name

	SelectedTempName *string `json:"selected_temp_name,omitempty"`
	// data

	Data *string `json:"data,omitempty"`
	// data_type

	DataType *int32 `json:"data_type,omitempty"`
}

func (o Contents) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Contents struct{}"
	}

	return strings.Join([]string{"Contents", string(data)}, " ")
}
