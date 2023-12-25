package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowCdnInfoReq 查桶对应的CDN信息
type ShowCdnInfoReq struct {

	// 源端桶的AK（最大长度100个字符），task_type为非url_list时，本参数为必选。
	Ak string `json:"ak"`

	// 源端桶的SK（最大长度100个字符），task_type为非url_list时，本参数为必选。
	Sk string `json:"sk"`

	// 云类型 AWS：亚马逊 Aliyun：阿里云 Qiniu：七牛云 QingCloud：青云 Tencent：腾讯云 Baidu：百度云 KingsoftCloud：金山云 Azure：微软云 UCloud：优刻得 HuaweiCloud：华为云 URLSource：URL HEC：HEC
	CloudType string `json:"cloud_type"`

	// 区域
	Region string `json:"region"`

	// 当源端为腾讯云时，会返回此参数。
	AppId *string `json:"app_id,omitempty"`

	// 桶名
	Bucket string `json:"bucket"`

	Prefix *PrefixKeyInfo `json:"prefix,omitempty"`

	SourceCdn *SourceCdnReq `json:"source_cdn"`
}

func (o ShowCdnInfoReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCdnInfoReq struct{}"
	}

	return strings.Join([]string{"ShowCdnInfoReq", string(data)}, " ")
}
