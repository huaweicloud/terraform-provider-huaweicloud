package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PrivateIpInfo 端口私有IP信息
type PrivateIpInfo struct {

	// 端口所属子网ID
	SubnetCidrId string `json:"subnet_cidr_id"`

	// 端口私有IP地址
	IpAddress string `json:"ip_address"`
}

func (o PrivateIpInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrivateIpInfo struct{}"
	}

	return strings.Join([]string{"PrivateIpInfo", string(data)}, " ")
}
