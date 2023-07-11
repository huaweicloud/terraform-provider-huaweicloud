package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type LiveSnapshotConfig struct {

	// 直播推流域名
	Domain string `json:"domain"`

	// 应用名称
	AppName string `json:"app_name"`

	// 回调鉴权密钥值  长度范围：[32-128]  若需要使用回调鉴权功能，请配置鉴权密钥，否则，留空即可。
	AuthKey *string `json:"auth_key,omitempty"`

	// 截图频率  取值范围：[5-3600]  单位：秒
	TimeInterval int32 `json:"time_interval"`

	// 在OBS桶存储截图的方式：  - 0：实时截图，以时间戳命名截图文件，保存所有截图文件到OBS桶。例：snapshot/{domain}/{app_name}/{stream_name}/{UnixTimestamp}.jpg  - 1：覆盖截图，只保存最新的截图文件，新的截图会覆盖原来的截图文件。例：snapshot/{domain}/{app_name}/{stream_name}.jpg
	ObjectWriteMode int32 `json:"object_write_mode"`

	ObsLocation *ObsFileAddr `json:"obs_location"`

	// 是否启用回调通知 - on：启用。 - off：不启用。
	CallBackEnable *LiveSnapshotConfigCallBackEnable `json:"call_back_enable,omitempty"`

	// 通知服务器地址，必须是合法的URL且携带协议，协议支持http和https。截图完成后直播服务会向此地址推送截图状态信息。
	CallBackUrl *string `json:"call_back_url,omitempty"`
}

func (o LiveSnapshotConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LiveSnapshotConfig struct{}"
	}

	return strings.Join([]string{"LiveSnapshotConfig", string(data)}, " ")
}

type LiveSnapshotConfigCallBackEnable struct {
	value string
}

type LiveSnapshotConfigCallBackEnableEnum struct {
	ON  LiveSnapshotConfigCallBackEnable
	OFF LiveSnapshotConfigCallBackEnable
}

func GetLiveSnapshotConfigCallBackEnableEnum() LiveSnapshotConfigCallBackEnableEnum {
	return LiveSnapshotConfigCallBackEnableEnum{
		ON: LiveSnapshotConfigCallBackEnable{
			value: "on",
		},
		OFF: LiveSnapshotConfigCallBackEnable{
			value: "off",
		},
	}
}

func (c LiveSnapshotConfigCallBackEnable) Value() string {
	return c.value
}

func (c LiveSnapshotConfigCallBackEnable) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *LiveSnapshotConfigCallBackEnable) UnmarshalJSON(b []byte) error {
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
