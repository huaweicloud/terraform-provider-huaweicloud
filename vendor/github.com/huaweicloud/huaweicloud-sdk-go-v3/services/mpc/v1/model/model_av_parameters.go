package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AvParameters struct {
	Video *VideoParameters `json:"video,omitempty"`

	Audio *Audio `json:"audio,omitempty"`

	Common *Common `json:"common"`
}

func (o AvParameters) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AvParameters struct{}"
	}

	return strings.Join([]string{"AvParameters", string(data)}, " ")
}
