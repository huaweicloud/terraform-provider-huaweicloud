package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TagEntity struct {

	// 标签键。  - 不能为空。  - 对于同一个实例，Key值唯一。  [- 长度不超过36个字符。](tag:cmcc,hws_ocb,ocb)  [- 长度为1~128个字符（中文也可以输入128个字符）。](tag:hws,hws_hk,ctc,hk_sbc,sbc,dt,g42,hk_g42,hk_tm,hws_eu,tm)  [- 不能包含“=”，“*”，“<”，“>”，“\\”，“,”，“|”，“/”](tag:cmcc,hws_ocb,ocb)  [- 由任意语种字母、数字、空格和字符组成，字符仅支持_ . : = + - @](tag:hws,hws_hk,ctc,hk_sbc,sbc,dt,g42,hk_g42,hk_tm,hws_eu,tm)  - 首尾字符不能为空格。
	Key *string `json:"key,omitempty"`

	// 标签值。  [- 不能为空。](tag:cmcc,hws_ocb,ocb)  [- 长度不超过43个字符。](tag:cmcc,hws_ocb,ocb)  [- 长度为0~255个字符（中文也可以输入255个字符）。](tag:hws,hws_hk,ctc,hk_sbc,sbc,dt,g42,hk_g42,hk_tm,hws_eu,tm)  [- 不能包含“=”，“*”，“<”，“>”，“\\”，“,”，“|”，“/”](tag:cmcc,hws_ocb,ocb)  [- 由任意语种字母、数字、空格和字符组成，字符仅支持_ . : = + - @](tag:hws,hws_hk,ctc,hk_sbc,sbc,dt,g42,hk_g42,hk_tm,hws_eu,tm)  - 首尾字符不能为空格。
	Value *string `json:"value,omitempty"`
}

func (o TagEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagEntity struct{}"
	}

	return strings.Join([]string{"TagEntity", string(data)}, " ")
}
