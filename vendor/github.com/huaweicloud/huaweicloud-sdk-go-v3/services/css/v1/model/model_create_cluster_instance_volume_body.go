package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClusterInstanceVolumeBody volume信息。当flavorRef选择的是本地盘规格时不需要填写,目前支持的本地盘规格有： - ess.spec-i3small - ess.spec-i3medium - ess.spec-i3.8xlarge.8 - ess.spec-ds.xlarge.8 - ess.spec-ds.2xlarge.8 - ess.spec-ds.4xlarge.8
type CreateClusterInstanceVolumeBody struct {

	// 卷类型。  - COMMON：普通I/O。 - HIGH：高I/O。 - ULTRAHIGH：超高I/O。
	VolumeType string `json:"volume_type"`

	// 卷大小，必须大于0且为4和10的公倍数，磁盘规格大小可以通过[获取实例规格列表](ListFlavors.xml)中diskrange属性获得。 单位：GB。
	Size int32 `json:"size"`
}

func (o CreateClusterInstanceVolumeBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterInstanceVolumeBody struct{}"
	}

	return strings.Join([]string{"CreateClusterInstanceVolumeBody", string(data)}, " ")
}
