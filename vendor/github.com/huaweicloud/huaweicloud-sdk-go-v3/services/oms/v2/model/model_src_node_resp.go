package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 源端节点信息。
type SrcNodeResp struct {

	// 源端桶的名称。
	Bucket *string `json:"bucket,omitempty"`

	// 源端云服务提供商。  可选值有AWS、Azure、Aliyun、Tencent、HuaweiCloud、QingCloud、KingsoftCloud、Baidu、Qiniu、URLSource或者UCloud。默认值为Aliyun。
	CloudType *SrcNodeRespCloudType `json:"cloud_type,omitempty"`

	// 源端桶所处的区域。
	Region *string `json:"region,omitempty"`

	// 当源端为腾讯云时，会返回此参数。
	AppId *string `json:"app_id,omitempty"`

	// 任务类型为对象迁移任务时，表示待迁移对象名称； 任务类型为前缀迁移任务时，表示待迁移前缀。
	ObjectKey *[]string `json:"object_key,omitempty"`

	ListFile *ListFile `json:"list_file,omitempty"`
}

func (o SrcNodeResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SrcNodeResp struct{}"
	}

	return strings.Join([]string{"SrcNodeResp", string(data)}, " ")
}

type SrcNodeRespCloudType struct {
	value string
}

type SrcNodeRespCloudTypeEnum struct {
	AWS            SrcNodeRespCloudType
	AZURE          SrcNodeRespCloudType
	ALIYUN         SrcNodeRespCloudType
	TENCENT        SrcNodeRespCloudType
	HUAWEI_CLOUD   SrcNodeRespCloudType
	QING_CLOUD     SrcNodeRespCloudType
	KINGSOFT_CLOUD SrcNodeRespCloudType
	BAIDU          SrcNodeRespCloudType
	QINIU          SrcNodeRespCloudType
	URL_SOURCE     SrcNodeRespCloudType
	U_CLOUD        SrcNodeRespCloudType
}

func GetSrcNodeRespCloudTypeEnum() SrcNodeRespCloudTypeEnum {
	return SrcNodeRespCloudTypeEnum{
		AWS: SrcNodeRespCloudType{
			value: "AWS",
		},
		AZURE: SrcNodeRespCloudType{
			value: "Azure",
		},
		ALIYUN: SrcNodeRespCloudType{
			value: "Aliyun",
		},
		TENCENT: SrcNodeRespCloudType{
			value: "Tencent",
		},
		HUAWEI_CLOUD: SrcNodeRespCloudType{
			value: "HuaweiCloud",
		},
		QING_CLOUD: SrcNodeRespCloudType{
			value: "QingCloud",
		},
		KINGSOFT_CLOUD: SrcNodeRespCloudType{
			value: "KingsoftCloud",
		},
		BAIDU: SrcNodeRespCloudType{
			value: "Baidu",
		},
		QINIU: SrcNodeRespCloudType{
			value: "Qiniu",
		},
		URL_SOURCE: SrcNodeRespCloudType{
			value: "URLSource",
		},
		U_CLOUD: SrcNodeRespCloudType{
			value: "UCloud",
		},
	}
}

func (c SrcNodeRespCloudType) Value() string {
	return c.value
}

func (c SrcNodeRespCloudType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SrcNodeRespCloudType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
