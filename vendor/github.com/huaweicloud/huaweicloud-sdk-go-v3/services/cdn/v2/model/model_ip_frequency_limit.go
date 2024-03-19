package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// IpFrequencyLimit IP访问限频，通过对单IP每秒访问单个节点的次数限制，实现CC攻击防御及恶意盗刷防护。
type IpFrequencyLimit struct {

	// 状态，on：打开，off：关闭。
	Status string `json:"status"`

	// 访问阈值,单位：次/秒，取值范围：1-100000。   > 当开启ip限频时，访问阈值必填。
	Qps *int32 `json:"qps,omitempty"`
}

func (o IpFrequencyLimit) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IpFrequencyLimit struct{}"
	}

	return strings.Join([]string{"IpFrequencyLimit", string(data)}, " ")
}
