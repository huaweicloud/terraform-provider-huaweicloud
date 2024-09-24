package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePortTagRequestBody This is a auto create Body Object
type CreatePortTagRequestBody struct {
	Tag *ResourceTag `json:"tag"`
}

func (o CreatePortTagRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePortTagRequestBody struct{}"
	}

	return strings.Join([]string{"CreatePortTagRequestBody", string(data)}, " ")
}
