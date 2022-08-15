package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type OriginHostRequest struct {
	OriginHost *OriginHostBody `json:"origin_host"`
}

func (o OriginHostRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OriginHostRequest struct{}"
	}

	return strings.Join([]string{"OriginHostRequest", string(data)}, " ")
}
