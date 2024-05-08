package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// InputStreamInfo 频道入流信息
type InputStreamInfo struct {

	// 频道入流协议 - FLV_PULL - RTMP_PUSH - RTMP_PULL - HLS_PULL - SRT_PULL - SRT_PUSH
	InputProtocol InputStreamInfoInputProtocol `json:"input_protocol"`

	// 频道主源流信息。入流协议为RTMP_PUSH和SRT_PUSH时，非必填项。其他情况下，均为必填项。
	Sources *[]SourcesInfo `json:"sources,omitempty"`

	// 备入流数组，非必填项。如果有备入流，则主备入流必须保证路数、codec和分辨率均一致。入流协议为RTMP_PUSH时，无需填写。
	SecondarySources *[]SecondarySourcesInfo `json:"secondary_sources,omitempty"`

	FailoverConditions *FailoverConditions `json:"failover_conditions,omitempty"`

	// 当入流协议为HLS_PULL时，最大带宽限制。 未配置会默认选择BANDWIDTH最高的流
	MaxBandwidthLimit *int32 `json:"max_bandwidth_limit,omitempty"`

	// 当推流协议为SRT_PUSH时，如果配置了直推源站，编码器不支持输入streamid，需要打开设置为true
	IpPortMode *bool `json:"ip_port_mode,omitempty"`
}

func (o InputStreamInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InputStreamInfo struct{}"
	}

	return strings.Join([]string{"InputStreamInfo", string(data)}, " ")
}

type InputStreamInfoInputProtocol struct {
	value string
}

type InputStreamInfoInputProtocolEnum struct {
	FLV_PULL  InputStreamInfoInputProtocol
	RTMP_PUSH InputStreamInfoInputProtocol
	RTMP_PULL InputStreamInfoInputProtocol
	HLS_PULL  InputStreamInfoInputProtocol
	SRT_PULL  InputStreamInfoInputProtocol
	SRT_PUSH  InputStreamInfoInputProtocol
}

func GetInputStreamInfoInputProtocolEnum() InputStreamInfoInputProtocolEnum {
	return InputStreamInfoInputProtocolEnum{
		FLV_PULL: InputStreamInfoInputProtocol{
			value: "FLV_PULL",
		},
		RTMP_PUSH: InputStreamInfoInputProtocol{
			value: "RTMP_PUSH",
		},
		RTMP_PULL: InputStreamInfoInputProtocol{
			value: "RTMP_PULL",
		},
		HLS_PULL: InputStreamInfoInputProtocol{
			value: "HLS_PULL",
		},
		SRT_PULL: InputStreamInfoInputProtocol{
			value: "SRT_PULL",
		},
		SRT_PUSH: InputStreamInfoInputProtocol{
			value: "SRT_PUSH",
		},
	}
}

func (c InputStreamInfoInputProtocol) Value() string {
	return c.value
}

func (c InputStreamInfoInputProtocol) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *InputStreamInfoInputProtocol) UnmarshalJSON(b []byte) error {
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
