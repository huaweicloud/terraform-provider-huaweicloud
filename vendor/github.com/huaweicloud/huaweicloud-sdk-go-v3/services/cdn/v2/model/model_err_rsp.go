package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ErrRsp 当北向接口报错时，按此格式返回到body体中
type ErrRsp struct {
	Error *ErrMsg `json:"error"`
}

func (o ErrRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ErrRsp struct{}"
	}

	return strings.Join([]string{"ErrRsp", string(data)}, " ")
}
