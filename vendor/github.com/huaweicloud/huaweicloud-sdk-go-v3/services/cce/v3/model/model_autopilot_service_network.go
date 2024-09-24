package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AutopilotServiceNetwork struct {

	// kubernetes clusterIP IPv4 CIDR取值范围。创建集群时若未传参，默认为\"10.247.0.0/16\"。
	IPv4CIDR *string `json:"IPv4CIDR,omitempty"`
}

func (o AutopilotServiceNetwork) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotServiceNetwork struct{}"
	}

	return strings.Join([]string{"AutopilotServiceNetwork", string(data)}, " ")
}
