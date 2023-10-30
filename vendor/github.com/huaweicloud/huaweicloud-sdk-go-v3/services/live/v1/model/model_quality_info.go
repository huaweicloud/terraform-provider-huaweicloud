package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type QualityInfo struct {

	// 自定义模板名称。 - 若需要自定义模板名称，请将quality参数设置为userdefine； - 多个自定义模板名称之间不能重复； - 自定义模板名称不能与其他模板的quality参数重复； - 若quality不为userdefine，请勿填写此字段。
	TemplateName *string `json:"templateName,omitempty"`

	// 包含如下取值： - lud： 超高清，系统缺省名称； - lhd： 高清，系统缺省名称； - lsd： 标清，系统缺省名称； - lld： 流畅，系统缺省名称； - userdefine： 视频质量自定义。填写userdefine时，templateName字段不能为空。
	Quality string `json:"quality"`

	// 是否使用窄带高清转码。默认值：off。  注意：该字段已不再维护，建议使用hdlb。  包含如下取值： - off：不启用。 - on：启用。
	Pvc *QualityInfoPvc `json:"PVC,omitempty"`

	// 是否启用高清低码，较PVC相比画质增强。默认值：off。  提示：使用hdlb字段开启高清低码时，PVC字段不生效。  包含如下取值： - off：不开启高清低码； - on：开启高清低码。
	Hdlb *QualityInfoHdlb `json:"hdlb,omitempty"`

	// 视频编码格式。默认为H264。 - H264：使用H.264。 - H265：使用H.265。
	Codec *QualityInfoCodec `json:"codec,omitempty"`

	// 视频长边（横屏的宽，竖屏的高）  单位：像素；默认值：0 - H264 建议取值范围：32-3840，必须为2的倍数 。 - H265 建议取值范围：320-3840 ，必须为2的倍数。  注意：width和height全为0，则输出分辨率和源一致；width和height只有一个为0， 则分辨率按非0项的比例缩放。
	Width *int32 `json:"width,omitempty"`

	// 视频短边（横屏的高，竖屏的宽）  单位：像素；默认值：0 - H264 建议取值范围：32-2160，必须为2的倍数。 - H265 建议取值范围：240-2160，必须为2的倍数。  注意：width和height全为0，则输出分辨率和源一致；width和height只有一个为0， 则分辨率按非0项的比例缩放。
	Height *int32 `json:"height,omitempty"`

	// 转码视频的码率  单位：Kbps  取值范围：40-30000
	Bitrate int32 `json:"bitrate"`

	// 转码视频帧率  单位：fps  默认值：0  取值范围：0-60，0表示保持帧率不变。
	VideoFrameRate *int32 `json:"video_frame_rate,omitempty"`

	// 转码输出支持的协议类型。默认为RTMP。当前只支持RTMP。  包含如下取值： - RTMP
	Protocol *QualityInfoProtocol `json:"protocol,omitempty"`

	// 最大I帧间隔  单位：帧数  取值范围：[0, 500]，默认值：50  注意：若希望通过iFrameInterval设置i帧间隔，请将gop设为0，或不传gop参数。
	IFrameInterval *int32 `json:"iFrameInterval,omitempty"`

	// 按时间设置I帧间隔  单位：秒  取值范围：[0,10]，默认值：2  注意：gop不为0时，则以gop设置i帧间隔，iFrameInterval字段不生效。
	Gop *int32 `json:"gop,omitempty"`

	// 自适应码率参数，默认值：off。  包含如下取值： - off：关闭码率自适应，目标码率按设定的码率输出； - minimum：目标码率按设定码率和源文件码率最小值输出（即码率不上扬）； - adaptive：目标码率按源文件码率自适应输出。
	BitrateAdaptive *QualityInfoBitrateAdaptive `json:"bitrate_adaptive,omitempty"`

	// 编码输出I帧策略，默认值：auto。  包含如下取值： - auto：I帧按设置的gop时长输出； - strictSync：编码输出I帧完全和源保持一致（源是I帧则编码输出I帧，源不是I帧则编码非I帧），设置该参数后gop时长设置无效。
	IFramePolicy *QualityInfoIFramePolicy `json:"i_frame_policy,omitempty"`
}

func (o QualityInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QualityInfo struct{}"
	}

	return strings.Join([]string{"QualityInfo", string(data)}, " ")
}

type QualityInfoPvc struct {
	value string
}

type QualityInfoPvcEnum struct {
	ON  QualityInfoPvc
	OFF QualityInfoPvc
}

func GetQualityInfoPvcEnum() QualityInfoPvcEnum {
	return QualityInfoPvcEnum{
		ON: QualityInfoPvc{
			value: "on",
		},
		OFF: QualityInfoPvc{
			value: "off",
		},
	}
}

func (c QualityInfoPvc) Value() string {
	return c.value
}

func (c QualityInfoPvc) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *QualityInfoPvc) UnmarshalJSON(b []byte) error {
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

type QualityInfoHdlb struct {
	value string
}

type QualityInfoHdlbEnum struct {
	ON  QualityInfoHdlb
	OFF QualityInfoHdlb
}

func GetQualityInfoHdlbEnum() QualityInfoHdlbEnum {
	return QualityInfoHdlbEnum{
		ON: QualityInfoHdlb{
			value: "on",
		},
		OFF: QualityInfoHdlb{
			value: "off",
		},
	}
}

func (c QualityInfoHdlb) Value() string {
	return c.value
}

func (c QualityInfoHdlb) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *QualityInfoHdlb) UnmarshalJSON(b []byte) error {
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

type QualityInfoCodec struct {
	value string
}

type QualityInfoCodecEnum struct {
	H264 QualityInfoCodec
	H265 QualityInfoCodec
}

func GetQualityInfoCodecEnum() QualityInfoCodecEnum {
	return QualityInfoCodecEnum{
		H264: QualityInfoCodec{
			value: "H264",
		},
		H265: QualityInfoCodec{
			value: "H265",
		},
	}
}

func (c QualityInfoCodec) Value() string {
	return c.value
}

func (c QualityInfoCodec) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *QualityInfoCodec) UnmarshalJSON(b []byte) error {
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

type QualityInfoProtocol struct {
	value string
}

type QualityInfoProtocolEnum struct {
	RTMP QualityInfoProtocol
	HLS  QualityInfoProtocol
	DASH QualityInfoProtocol
}

func GetQualityInfoProtocolEnum() QualityInfoProtocolEnum {
	return QualityInfoProtocolEnum{
		RTMP: QualityInfoProtocol{
			value: "RTMP",
		},
		HLS: QualityInfoProtocol{
			value: "HLS",
		},
		DASH: QualityInfoProtocol{
			value: "DASH",
		},
	}
}

func (c QualityInfoProtocol) Value() string {
	return c.value
}

func (c QualityInfoProtocol) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *QualityInfoProtocol) UnmarshalJSON(b []byte) error {
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

type QualityInfoBitrateAdaptive struct {
	value string
}

type QualityInfoBitrateAdaptiveEnum struct {
	MINIMUM  QualityInfoBitrateAdaptive
	ADAPTIVE QualityInfoBitrateAdaptive
}

func GetQualityInfoBitrateAdaptiveEnum() QualityInfoBitrateAdaptiveEnum {
	return QualityInfoBitrateAdaptiveEnum{
		MINIMUM: QualityInfoBitrateAdaptive{
			value: "minimum",
		},
		ADAPTIVE: QualityInfoBitrateAdaptive{
			value: "adaptive",
		},
	}
}

func (c QualityInfoBitrateAdaptive) Value() string {
	return c.value
}

func (c QualityInfoBitrateAdaptive) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *QualityInfoBitrateAdaptive) UnmarshalJSON(b []byte) error {
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

type QualityInfoIFramePolicy struct {
	value string
}

type QualityInfoIFramePolicyEnum struct {
	AUTO        QualityInfoIFramePolicy
	STRICT_SYNC QualityInfoIFramePolicy
}

func GetQualityInfoIFramePolicyEnum() QualityInfoIFramePolicyEnum {
	return QualityInfoIFramePolicyEnum{
		AUTO: QualityInfoIFramePolicy{
			value: "auto",
		},
		STRICT_SYNC: QualityInfoIFramePolicy{
			value: "strictSync",
		},
	}
}

func (c QualityInfoIFramePolicy) Value() string {
	return c.value
}

func (c QualityInfoIFramePolicy) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *QualityInfoIFramePolicy) UnmarshalJSON(b []byte) error {
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
