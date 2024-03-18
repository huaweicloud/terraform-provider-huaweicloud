package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SmartConnectTaskReqSinkConfig struct {

	// Redis实例地址。（仅目标端类型为Redis时需要填写）
	RedisAddress *string `json:"redis_address,omitempty"`

	// Redis实例类型。（仅目标端类型为Redis时需要填写）
	RedisType *string `json:"redis_type,omitempty"`

	// DCS实例ID。（仅目标端类型为Redis时需要填写）
	DcsInstanceId *string `json:"dcs_instance_id,omitempty"`

	// Redis密码。（仅目标端类型为Redis时需要填写）
	RedisPassword *string `json:"redis_password,omitempty"`

	// 转储启动偏移量，latest为获取最新的数据，earliest为获取最早的数据。（仅目标端类型为OBS时需要填写）
	ConsumerStrategy *string `json:"consumer_strategy,omitempty"`

	// 转储文件格式。当前只支持TEXT。（仅目标端类型为OBS时需要填写）
	DestinationFileType *string `json:"destination_file_type,omitempty"`

	// 数据转储周期（秒），默认配置为300秒。（仅目标端类型为OBS时需要填写）
	DeliverTimeInterval *int32 `json:"deliver_time_interval,omitempty"`

	// AK，访问密钥ID。（仅目标端类型为OBS时需要填写）
	AccessKey *string `json:"access_key,omitempty"`

	// SK，与访问密钥ID结合使用的密钥。（仅目标端类型为OBS时需要填写）
	SecretKey *string `json:"secret_key,omitempty"`

	// 转储地址，即存储Topic数据的OBS桶的名称。（仅目标端类型为OBS时需要填写）
	ObsBucketName *string `json:"obs_bucket_name,omitempty"`

	// 转储目录，即OBS中存储Topic的目录，多级目录可以用“/”进行分隔。（仅目标端类型为OBS时需要填写）
	ObsPath *string `json:"obs_path,omitempty"`

	// 时间目录格式。（仅目标端类型为OBS时需要填写）   - yyyy：年   - yyyy/MM：年/月   - yyyy/MM/dd：年/月/日   - yyyy/MM/dd/HH：年/月/日/时   - yyyy/MM/dd/HH/mm：年/月/日/时/分
	PartitionFormat *string `json:"partition_format,omitempty"`

	//  记录分行符，用于分隔写入转储文件的用户数据。（仅目标端类型为OBS时需要填写）   取值范围：   - 逗号“,”   - 分号“;”   - 竖线“|”   - 换行符“\\n”   - NULL
	RecordDelimiter *string `json:"record_delimiter,omitempty"`

	// 是否转储Key，开启表示转储Key，关闭表示不转储Key。（仅目标端类型为OBS时需要填写）
	StoreKeys *bool `json:"store_keys,omitempty"`
}

func (o SmartConnectTaskReqSinkConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SmartConnectTaskReqSinkConfig struct{}"
	}

	return strings.Join([]string{"SmartConnectTaskReqSinkConfig", string(data)}, " ")
}
