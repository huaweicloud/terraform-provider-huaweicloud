package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// TaskGroupSrcNodeResp 源端迁移任务组节点信息。
type TaskGroupSrcNodeResp struct {

	// 源端桶的名称。
	Bucket *string `json:"bucket,omitempty"`

	// 源端云服务提供商。  可选值有AWS、AZURE、ALIYUN、TENCENT、HUAWEICLOUD、QINGCLOUD、KINGSOFTCLOUD、BAIDU、QINIU、GOOGLE、URLSOURCE或者UCLOUD。默认值为ALIYUN。
	CloudType *TaskGroupSrcNodeRespCloudType `json:"cloud_type,omitempty"`

	// 源端桶所处的区域。
	Region *string `json:"region,omitempty"`

	// 当源端为腾讯云时，会返回此参数。
	AppId *string `json:"app_id,omitempty"`

	// 任务组类型为前缀迁移任务时，表示待迁移前缀。
	ObjectKey *[]string `json:"object_key,omitempty"`

	ListFile *ListFile `json:"list_file,omitempty"`
}

func (o TaskGroupSrcNodeResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskGroupSrcNodeResp struct{}"
	}

	return strings.Join([]string{"TaskGroupSrcNodeResp", string(data)}, " ")
}

type TaskGroupSrcNodeRespCloudType struct {
	value string
}

type TaskGroupSrcNodeRespCloudTypeEnum struct {
	AWS           TaskGroupSrcNodeRespCloudType
	AZURE         TaskGroupSrcNodeRespCloudType
	ALIYUN        TaskGroupSrcNodeRespCloudType
	TENCENT       TaskGroupSrcNodeRespCloudType
	HUAWEICLOUD   TaskGroupSrcNodeRespCloudType
	QINGCLOUD     TaskGroupSrcNodeRespCloudType
	KINGSOFTCLOUD TaskGroupSrcNodeRespCloudType
	BAIDU         TaskGroupSrcNodeRespCloudType
	QINIU         TaskGroupSrcNodeRespCloudType
	URLSOURCE     TaskGroupSrcNodeRespCloudType
	UCLOUD        TaskGroupSrcNodeRespCloudType
	GOOGLE        TaskGroupSrcNodeRespCloudType
}

func GetTaskGroupSrcNodeRespCloudTypeEnum() TaskGroupSrcNodeRespCloudTypeEnum {
	return TaskGroupSrcNodeRespCloudTypeEnum{
		AWS: TaskGroupSrcNodeRespCloudType{
			value: "AWS",
		},
		AZURE: TaskGroupSrcNodeRespCloudType{
			value: "AZURE",
		},
		ALIYUN: TaskGroupSrcNodeRespCloudType{
			value: "ALIYUN",
		},
		TENCENT: TaskGroupSrcNodeRespCloudType{
			value: "TENCENT",
		},
		HUAWEICLOUD: TaskGroupSrcNodeRespCloudType{
			value: "HUAWEICLOUD",
		},
		QINGCLOUD: TaskGroupSrcNodeRespCloudType{
			value: "QINGCLOUD",
		},
		KINGSOFTCLOUD: TaskGroupSrcNodeRespCloudType{
			value: "KINGSOFTCLOUD",
		},
		BAIDU: TaskGroupSrcNodeRespCloudType{
			value: "BAIDU",
		},
		QINIU: TaskGroupSrcNodeRespCloudType{
			value: "QINIU",
		},
		URLSOURCE: TaskGroupSrcNodeRespCloudType{
			value: "URLSOURCE",
		},
		UCLOUD: TaskGroupSrcNodeRespCloudType{
			value: "UCLOUD",
		},
		GOOGLE: TaskGroupSrcNodeRespCloudType{
			value: "GOOGLE",
		},
	}
}

func (c TaskGroupSrcNodeRespCloudType) Value() string {
	return c.value
}

func (c TaskGroupSrcNodeRespCloudType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskGroupSrcNodeRespCloudType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
