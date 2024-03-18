package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VulScanTaskHostInfo struct {

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 弹性公网IP地址
	PublicIp *string `json:"public_ip,omitempty"`

	// 私有IP地址
	PrivateIp *string `json:"private_ip,omitempty"`

	// 资产重要性，包含如下:   - important ：重要资产   - common ：一般资产   - test ：测试资产
	AssetValue *string `json:"asset_value,omitempty"`

	// 主机的扫描状态，包含如下：   -scanning : 扫描中   -success : 扫描成功   -failed : 扫描失败
	ScanStatus *string `json:"scan_status,omitempty"`

	// 扫描失败的原因列表
	FailedReasons *[]VulScanTaskHostInfoFailedReasons `json:"failed_reasons,omitempty"`
}

func (o VulScanTaskHostInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulScanTaskHostInfo struct{}"
	}

	return strings.Join([]string{"VulScanTaskHostInfo", string(data)}, " ")
}
