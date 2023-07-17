package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAgenciesTaskRequest Request Object
type ShowAgenciesTaskRequest struct {
}

func (o ShowAgenciesTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAgenciesTaskRequest struct{}"
	}

	return strings.Join([]string{"ShowAgenciesTaskRequest", string(data)}, " ")
}
