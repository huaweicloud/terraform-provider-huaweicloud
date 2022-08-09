package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

//
type UpdateWatermarkTemplateReq struct {

	// 水印模板配置id<br/>
	Id string `json:"id"`

	// 水印模板名称<br/>
	Name *string `json:"name,omitempty"`

	// 水印图片相对输出视频的水平偏移量，默认值是0<br/>
	Dx *string `json:"dx,omitempty"`

	// 水印图片相对输出视频的垂直偏移量，默认值是0<br/>
	Dy *string `json:"dy,omitempty"`

	// 水印的位置<br/>
	Position *UpdateWatermarkTemplateReqPosition `json:"position,omitempty"`

	// 水印图片宽<br/>
	Width *string `json:"width,omitempty"`

	// 水印图片高<br/>
	Height *string `json:"height,omitempty"`

	// 水印类型，当前只支持Image（图片水印）<br/>
	WatermarkType *UpdateWatermarkTemplateReqWatermarkType `json:"watermark_type,omitempty"`

	// type设置为Image时有效。  目前包括： - Original：只做简单缩放，不做其他处理 - Transparent：图片底色透明 - Grayed：彩色图片变灰
	ImageProcess *UpdateWatermarkTemplateReqImageProcess `json:"image_process,omitempty"`

	// 水印开始时间，与\"timeline_duration\"配合使用。 取值范围:[0, END)。\"END\"表示视频结束时间。 单位:秒。
	TimelineStart *string `json:"timeline_start,omitempty"`

	// 水印持续时间，与\"timeline_start\"配合使用。 取值范围:(0,END-开始时间]。\"END\"表示视频结束时间。 单位:秒。 默认:END。
	TimelineDuration *string `json:"timeline_duration,omitempty"`
}

func (o UpdateWatermarkTemplateReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateWatermarkTemplateReq struct{}"
	}

	return strings.Join([]string{"UpdateWatermarkTemplateReq", string(data)}, " ")
}

type UpdateWatermarkTemplateReqPosition struct {
	value string
}

type UpdateWatermarkTemplateReqPositionEnum struct {
	TOPRIGHT    UpdateWatermarkTemplateReqPosition
	TOPLEFT     UpdateWatermarkTemplateReqPosition
	BOTTOMRIGHT UpdateWatermarkTemplateReqPosition
	BOTTOMLEFT  UpdateWatermarkTemplateReqPosition
}

func GetUpdateWatermarkTemplateReqPositionEnum() UpdateWatermarkTemplateReqPositionEnum {
	return UpdateWatermarkTemplateReqPositionEnum{
		TOPRIGHT: UpdateWatermarkTemplateReqPosition{
			value: "TOPRIGHT",
		},
		TOPLEFT: UpdateWatermarkTemplateReqPosition{
			value: "TOPLEFT",
		},
		BOTTOMRIGHT: UpdateWatermarkTemplateReqPosition{
			value: "BOTTOMRIGHT",
		},
		BOTTOMLEFT: UpdateWatermarkTemplateReqPosition{
			value: "BOTTOMLEFT",
		},
	}
}

func (c UpdateWatermarkTemplateReqPosition) Value() string {
	return c.value
}

func (c UpdateWatermarkTemplateReqPosition) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateWatermarkTemplateReqPosition) UnmarshalJSON(b []byte) error {
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

type UpdateWatermarkTemplateReqWatermarkType struct {
	value string
}

type UpdateWatermarkTemplateReqWatermarkTypeEnum struct {
	IMAGE UpdateWatermarkTemplateReqWatermarkType
	TEXT  UpdateWatermarkTemplateReqWatermarkType
}

func GetUpdateWatermarkTemplateReqWatermarkTypeEnum() UpdateWatermarkTemplateReqWatermarkTypeEnum {
	return UpdateWatermarkTemplateReqWatermarkTypeEnum{
		IMAGE: UpdateWatermarkTemplateReqWatermarkType{
			value: "IMAGE",
		},
		TEXT: UpdateWatermarkTemplateReqWatermarkType{
			value: "TEXT",
		},
	}
}

func (c UpdateWatermarkTemplateReqWatermarkType) Value() string {
	return c.value
}

func (c UpdateWatermarkTemplateReqWatermarkType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateWatermarkTemplateReqWatermarkType) UnmarshalJSON(b []byte) error {
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

type UpdateWatermarkTemplateReqImageProcess struct {
	value string
}

type UpdateWatermarkTemplateReqImageProcessEnum struct {
	ORIGINAL    UpdateWatermarkTemplateReqImageProcess
	TRANSPARENT UpdateWatermarkTemplateReqImageProcess
	GRAYED      UpdateWatermarkTemplateReqImageProcess
}

func GetUpdateWatermarkTemplateReqImageProcessEnum() UpdateWatermarkTemplateReqImageProcessEnum {
	return UpdateWatermarkTemplateReqImageProcessEnum{
		ORIGINAL: UpdateWatermarkTemplateReqImageProcess{
			value: "ORIGINAL",
		},
		TRANSPARENT: UpdateWatermarkTemplateReqImageProcess{
			value: "TRANSPARENT",
		},
		GRAYED: UpdateWatermarkTemplateReqImageProcess{
			value: "GRAYED",
		},
	}
}

func (c UpdateWatermarkTemplateReqImageProcess) Value() string {
	return c.value
}

func (c UpdateWatermarkTemplateReqImageProcess) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateWatermarkTemplateReqImageProcess) UnmarshalJSON(b []byte) error {
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
