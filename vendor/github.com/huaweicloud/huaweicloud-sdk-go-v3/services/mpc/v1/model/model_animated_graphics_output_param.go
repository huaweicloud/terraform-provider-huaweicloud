package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type AnimatedGraphicsOutputParam struct {

	// 动图格式，目前仅支持取值 gif
	Format AnimatedGraphicsOutputParamFormat `json:"format"`

	// 输出动图的宽。  取值范围：0，-1或[32,3840]之间2的倍数。  >- 若设置为-1， 则宽根据高来自适应，此时“height”不能取-1或0。 >- 若设置为0，则取原始视频的宽，此时“height”只能取0。
	Width int32 `json:"width"`

	// 输出动图的高。  取值范围：0，-1或[32,2160]之间2的倍数。  >- 若设置为-1， 则高根据宽来自适应，此时“width”不能取-1或0。 >- 若设置为0，则取原始视频的高，此时“width”只能取0。
	Height int32 `json:"height"`

	// 起始时间，单位：毫秒
	Start int32 `json:"start"`

	// 结束时间。  单位：毫秒。  end、start差值最多60秒。
	End int32 `json:"end"`

	// 动图帧率。  取值范围：[1,75]
	FrameRate *int32 `json:"frame_rate,omitempty"`
}

func (o AnimatedGraphicsOutputParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AnimatedGraphicsOutputParam struct{}"
	}

	return strings.Join([]string{"AnimatedGraphicsOutputParam", string(data)}, " ")
}

type AnimatedGraphicsOutputParamFormat struct {
	value string
}

type AnimatedGraphicsOutputParamFormatEnum struct {
	GIF AnimatedGraphicsOutputParamFormat
}

func GetAnimatedGraphicsOutputParamFormatEnum() AnimatedGraphicsOutputParamFormatEnum {
	return AnimatedGraphicsOutputParamFormatEnum{
		GIF: AnimatedGraphicsOutputParamFormat{
			value: "gif",
		},
	}
}

func (c AnimatedGraphicsOutputParamFormat) Value() string {
	return c.value
}

func (c AnimatedGraphicsOutputParamFormat) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AnimatedGraphicsOutputParamFormat) UnmarshalJSON(b []byte) error {
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
