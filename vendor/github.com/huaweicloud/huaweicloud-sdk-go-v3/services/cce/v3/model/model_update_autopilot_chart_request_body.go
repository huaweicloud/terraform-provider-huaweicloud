package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"encoding/json"
	"errors"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
	"os"
	"reflect"

	"strings"
)

type UpdateAutopilotChartRequestBody struct {

	// 上传模板的配置参数，示例如下：\"{\\\"override\\\":true,\\\"skip_lint\\\":true,\\\"source\\\":\\\"package\\\"}\" - skip_lint: whether lint uploaded chart - override: whether override existed chart - visible: update chart visible
	Parameters *def.MultiPart `json:"parameters,omitempty"`

	// 模板包文件
	Content *def.FilePart `json:"content"`
}

func (o UpdateAutopilotChartRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAutopilotChartRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateAutopilotChartRequestBody", string(data)}, " ")
}

func (o *UpdateAutopilotChartRequestBody) UnmarshalJSON(b []byte) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	t := reflect.TypeOf(o).Elem()
	v := reflect.ValueOf(o).Elem()
	count := v.NumField()
	for i := 0; i < count; i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		jsonName := strings.Split(jsonTag, ",")[0]
		if m[jsonName] == nil && strings.Contains(jsonTag, "omitempty") {
			continue
		}
		field := v.FieldByName(utils.UnderscoreToCamel(jsonName))
		switch v.Field(i).Interface().(type) {
		case *def.FilePart:
			filePath := m[jsonName].(string)
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(def.NewFilePart(file)))
		case *def.MultiPart:
			field.Set(reflect.ValueOf(def.NewMultiPart(m[jsonName])))
		default:
			return errors.New(fmt.Sprintf("unmarshal %s failed", m[jsonName]))
		}
	}
	return nil
}
