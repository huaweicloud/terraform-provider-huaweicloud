package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBucketRegionReq 查询桶对应的region
type ShowBucketRegionReq struct {

	// 云类型 AWS：亚马逊 Aliyun：阿里云 Qiniu：七牛云 QingCloud：青云 Tencent：腾讯云 Baidu：百度云 KingsoftCloud：金山云 Azure：微软云 UCloud：优刻得 HuaweiCloud：华为云 Google: 谷歌云 URLSource：URL HEC：HEC
	CloudType string `json:"cloud_type"`

	// 源端桶的AK（最大长度100个字符），task_type为非url_list时，本参数为必选。
	Ak *string `json:"ak,omitempty"`

	// 源端桶的SK（最大长度100个字符），task_type为非url_list时，本参数为必选。
	Sk *string `json:"sk,omitempty"`

	// 用于谷歌云Cloud Storage鉴权
	JsonAuthFile *string `json:"json_auth_file,omitempty"`

	// 当源端为腾讯云时，会返回此参数。
	AppId *string `json:"app_id,omitempty"`

	// 桶名
	BucketName string `json:"bucket_name"`
}

func (o ShowBucketRegionReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBucketRegionReq struct{}"
	}

	return strings.Join([]string{"ShowBucketRegionReq", string(data)}, " ")
}
