package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 截图参数
type Thumbnail struct {

	// 截图类型。  取值如下： - time：每次进行截图的间隔时间。 - dots: 按照指定的时间点截图。
	Type ThumbnailType `json:"type"`

	// **type**取值为time时必填。根据时间间隔采样时的时间间隔值。  取值范围：[1,12]之间的整数。  单位：秒。
	Time *int32 `json:"time,omitempty"`

	// **type**取值为dots时必填。指定时间截图时的时间点数组。
	Dots *[]int32 `json:"dots,omitempty"`

	// 该值表示指定第几张截图作为封面。  默认值：1。
	CoverPosition *int32 `json:"cover_position,omitempty"`

	// 截图文件格式。  取值如下： - 1：jpg。  默认值：1 。
	Format *int32 `json:"format,omitempty"`

	// 纵横比，图像缩放方式。  取值如下： - 0：自适应（保持原有宽高比）。 - 1：16:9。  默认值：0。
	AspectRatio *int32 `json:"aspect_ratio,omitempty"`

	// 截图最长边的尺寸。  单位：像素。  宽边尺寸按照该尺寸与原始视频像素等比缩放计算。
	MaxLength *int32 `json:"max_length,omitempty"`
}

func (o Thumbnail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Thumbnail struct{}"
	}

	return strings.Join([]string{"Thumbnail", string(data)}, " ")
}

type ThumbnailType struct {
	value string
}

type ThumbnailTypeEnum struct {
	TIME ThumbnailType
	DOTS ThumbnailType
}

func GetThumbnailTypeEnum() ThumbnailTypeEnum {
	return ThumbnailTypeEnum{
		TIME: ThumbnailType{
			value: "time",
		},
		DOTS: ThumbnailType{
			value: "dots",
		},
	}
}

func (c ThumbnailType) Value() string {
	return c.value
}

func (c ThumbnailType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ThumbnailType) UnmarshalJSON(b []byte) error {
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
