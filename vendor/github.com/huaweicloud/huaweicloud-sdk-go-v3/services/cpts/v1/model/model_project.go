package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Project struct {
	// create_time

	CreateTime *string `json:"create_time,omitempty"`
	// description

	Description *string `json:"description,omitempty"`
	// group

	Group *string `json:"group,omitempty"`
	// id

	Id *int32 `json:"id,omitempty"`
	// name

	Name *string `json:"name,omitempty"`
	// source

	Source *int32 `json:"source,omitempty"`
	// update_time

	UpdateTime *string `json:"update_time,omitempty"`
}

func (o Project) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Project struct{}"
	}

	return strings.Join([]string{"Project", string(data)}, " ")
}
