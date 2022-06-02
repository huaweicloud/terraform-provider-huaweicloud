package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type ThumbnailPara struct {

	// 采样类型。  取值如下： - \"TIME\"：根据时间间隔采样截图。 - \"DOTS\"：指定时间点截图。选择同步截图时，需指定此类型。  默认值：\"TIME\"
	Type *ThumbnailParaType `json:"type,omitempty"`

	// 采样截图的时间间隔值。  默认值：12。  单位：秒
	Time *int32 `json:"time,omitempty"`

	// 采样类型为“TIME”模式的开始时间，和“time”配合使用。  默认值：0。  单位：秒。
	StartTime *int32 `json:"start_time,omitempty"`

	// 采样类型为“TIME”模式的持续时间，和“time”、“start_time”配合使用，表示从视频文件的第“start_time”开始，持续时间为“duration”，每间隔“time”生成一张截图。 取值范围：[数字，ToEND]。“ToEND”表示持续到视频结束。  默认值： ToEND。  单位：秒。 > “duration”必须大于等0，若设置为0，则截图持续时间从“start_time”到视频结束。
	Duration *int32 `json:"duration,omitempty"`

	// 指定时间截图时的时间点数组，最多支持10个。
	Dots *[]int32 `json:"dots,omitempty"`

	// 截图输出文件名。  - 如果只抽一张图（即：按DOTS方式，指定1个时间点）则按该指定文件名输出图片。  - 如果抽多张图（即：按DOTS方式指定多个时间点或按TIME间隔截图）则输出图片名在该指定文件名基础上在增加时间点（示例：output_filename_10.jpg）。  - 如果指定了压缩抽帧图片生成tar包，则tar包按该指定文件名输出。
	OutputFilename *string `json:"output_filename,omitempty"`

	// 截图文件格式。  取值如下：  1：表示jpg格式
	Format *int32 `json:"format,omitempty"`

	// 纵横比。
	AspectRatio *int32 `json:"aspect_ratio,omitempty"`

	// 图片宽度  取值范围：(96,3840]  单位：px
	Width *int32 `json:"width,omitempty"`

	// 图片高度  取值范围：(96,2160]  单位：px
	Height *int32 `json:"height,omitempty"`

	// 截图最长边的尺寸。宽边尺寸按照该尺寸与原始视频像素等比缩放计算。   取值范围：[240,3840]  默认值：480  单位：像素  > 该参数和width/height选择使用，以width/height优先，若width/height都不等于0，则图片尺寸按width/height得出；反之，则图片尺寸按 max_length 得出。  > 若该参数和width/height都未选择，则按源片源宽高输出截图
	MaxLength *int32 `json:"max_length,omitempty"`
}

func (o ThumbnailPara) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ThumbnailPara struct{}"
	}

	return strings.Join([]string{"ThumbnailPara", string(data)}, " ")
}

type ThumbnailParaType struct {
	value string
}

type ThumbnailParaTypeEnum struct {
	PERCENT ThumbnailParaType
	TIME    ThumbnailParaType
	DOTS    ThumbnailParaType
}

func GetThumbnailParaTypeEnum() ThumbnailParaTypeEnum {
	return ThumbnailParaTypeEnum{
		PERCENT: ThumbnailParaType{
			value: "PERCENT",
		},
		TIME: ThumbnailParaType{
			value: "TIME",
		},
		DOTS: ThumbnailParaType{
			value: "DOTS",
		},
	}
}

func (c ThumbnailParaType) Value() string {
	return c.value
}

func (c ThumbnailParaType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ThumbnailParaType) UnmarshalJSON(b []byte) error {
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
