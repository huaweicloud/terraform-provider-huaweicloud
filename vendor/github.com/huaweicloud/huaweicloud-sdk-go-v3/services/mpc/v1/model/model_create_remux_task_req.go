package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateRemuxTaskReq struct {
	Input *ObsObjInfo `json:"input,omitempty"`

	Output *ObsObjInfo `json:"output,omitempty"`

	// 用户自定义数据。
	UserData *string `json:"user_data,omitempty"`

	OutputParam *RemuxOutputParam `json:"output_param,omitempty"`
}

func (o CreateRemuxTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRemuxTaskReq struct{}"
	}

	return strings.Join([]string{"CreateRemuxTaskReq", string(data)}, " ")
}
