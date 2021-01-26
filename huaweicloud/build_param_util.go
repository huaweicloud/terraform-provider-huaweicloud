package huaweicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type funcSkipOpt func(optn string, jsonTags []string, tag reflect.StructTag) bool

type exchangeParam struct {
	OptNameMap *map[string]string
	SkipOpt    funcSkipOpt
}

func (e *exchangeParam) toSchemaOptName(name string) string {
	if e.OptNameMap != nil {
		if n, ok := (*e.OptNameMap)[name]; ok {
			return n
		}
	}
	return strings.ToLower(name)
}

func (e *exchangeParam) BuildCUParam(opts interface{}, d *schema.ResourceData) ([]string, error) {
	var skippedFields []string

	optsValue := reflect.ValueOf(opts)
	if optsValue.Kind() != reflect.Ptr {
		return skippedFields, fmt.Errorf("parameter of opts should be a pointer")
	}
	optsValue = optsValue.Elem()
	if optsValue.Kind() != reflect.Struct {
		return skippedFields, fmt.Errorf("parameter must be a pointer to a struct")
	}

	optsType := reflect.TypeOf(opts)
	optsType = optsType.Elem()
	value := make(map[string]interface{})

	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)
		tag := getParamTag("json", f.Tag)
		if tag == "" {
			return skippedFields, fmt.Errorf("can not convert for item %v: without of json tag", v)
		}
		tags := strings.Split(tag, ",")
		fieldn := tags[0]
		optn := e.toSchemaOptName(fieldn)

		// Only check the parameters in top struct.
		// If the parameters in sub-struct need skip, it will miss.
		// If it happens, need refactor here.
		if e.SkipOpt(optn, tags, f.Tag) {
			skippedFields = append(skippedFields, fieldn)
			continue
		}
		optv := d.Get(optn)
		if optv == nil {
			log.Printf("[DEBUG] opt:%s is not set", optn)
			continue
		}
		value[optn] = optv
	}
	if len(value) == 0 {
		log.Printf("[WARN]no parameter was set")
		return skippedFields, nil
	}
	return skippedFields, e.buildStruct(&optsValue, optsType, value)
}

func (e *exchangeParam) buildStruct(optsValue *reflect.Value, optsType reflect.Type, value map[string]interface{}) error {
	log.Printf("[DEBUG] buildStruct:: optsValue=%v, optsType=%v, value=%#v\n", optsValue, optsType, value)

	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" {
			return fmt.Errorf("can not convert for item %v: without of json tag", v)
		}
		fieldn := strings.Split(tag, ",")[0]
		optn := e.toSchemaOptName(fieldn)
		log.Printf("[DEBUG] buildStruct:: convert for opt:%s", optn)

		optv, ok := value[optn]
		if !ok {
			log.Printf("[DEBUG] field:%s was not supplied", fieldn)
			continue
		}

		switch v.Kind() {
		case reflect.String:
			v.SetString(optv.(string))
		case reflect.Int:
			v.SetInt(int64(optv.(int)))
		case reflect.Int64:
			v.SetInt(optv.(int64))
		case reflect.Bool:
			v.SetBool(optv.(bool))
		case reflect.Slice:
			s := optv.([]interface{})

			switch v.Type().Elem().Kind() {
			case reflect.String:
				t := make([]string, len(s))
				for i, iv := range s {
					t[i] = iv.(string)
				}
				v.Set(reflect.ValueOf(t))
			case reflect.Struct:
				t := reflect.MakeSlice(f.Type, len(s), len(s))
				for i, iv := range s {
					rv := t.Index(i)
					err := e.buildStruct(&rv, f.Type.Elem(), iv.(map[string]interface{}))
					if err != nil {
						return err
					}
				}
				v.Set(t)

			default:
				return fmt.Errorf("unknown type of item %v: %v", v, v.Type().Elem().Kind())
			}
		case reflect.Struct:
			log.Printf("[DEBUG] buildStruct:: convert struct for opt:%s, value:%#v", optn, optv)
			var p map[string]interface{}
			ok := true

			// If the type of parameter is Struct, then the corresponding type in Schema is TypeList
			if v0, ok0 := optv.([]interface{}); ok0 {
				p, ok = v0[0].(map[string]interface{})
			} else {
				p, ok = optv.(map[string]interface{})
			}
			if !ok {
				return fmt.Errorf("can not convert to (map[string]interface{}) for opt:%s, value:%#v", optn, optv)
			}

			err := e.buildStruct(&v, f.Type, p)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown type of item %v: %v", v, v.Kind())
		}
	}
	return nil
}

func (e *exchangeParam) BuildResourceData(resp interface{}, d *schema.ResourceData) error {
	value, err := e.convertToMap(resp)
	if err != nil {
		return err
	}

	optsValue := reflect.ValueOf(resp)
	if optsValue.Kind() == reflect.Ptr {
		optsValue = optsValue.Elem()
	}

	optsType := reflect.TypeOf(resp)
	if optsType.Kind() == reflect.Ptr {
		optsType = optsType.Elem()
	}

	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" {
			return fmt.Errorf("can not convert for item %v: without of json tag", v)
		}
		fieldn := strings.Split(tag, ",")[0]
		optn := e.toSchemaOptName(fieldn)
		if optn == "id" {
			continue
		}
		optv := value[optn]
		log.Printf("[DEBUG BuildResourceData: set for opt:%s, value:%#v", optn, optv)

		switch v.Kind() {
		default:
			err := d.Set(optn, optv) //lintignore:R001
			if err != nil {
				return err
			}
		case reflect.Struct:
			//The corresponding schema of Struct is TypeList in Terrafrom
			err := d.Set(optn, []interface{}{optv}) //lintignore:R001
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *exchangeParam) convertToMap(resp interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("refreshResourceData:: marshal failed:%v", err)
	}

	m := regexp.MustCompile(`"[a-z0-9A-Z_]+":`)
	nb := m.ReplaceAllFunc(
		b,
		func(src []byte) []byte {
			k := fmt.Sprintf("%s", src[1:len(src)-2])
			return []byte(fmt.Sprintf("\"%s\":", e.toSchemaOptName(k)))
		},
	)
	log.Printf("[DEBUG]befor change b =%s, b=%#v", b, b)
	log.Printf("[DEBUG] after change nb=%s, nb=%#v", nb, nb)

	p := make(map[string]interface{})
	err = json.Unmarshal(nb, &p)
	if err != nil {
		return nil, fmt.Errorf("refreshResourceData:: unmarshal failed:%v", err)
	}
	log.Printf("[DEBUG]refreshResourceData:: raw data = %#v\n", p)
	return p, nil
}

// The result may be not correct when the type of param is string and user config it to 'param=""'
// but, there is no other way.
func hasFilledOpt(d *schema.ResourceData, param string) bool {
	_, b := d.GetOkExists(param)
	return b
}

func getParamTag(key string, tag reflect.StructTag) string {
	if v, ok := tag.Lookup(key); ok {
		return v
	}
	return "tag_not_set"
}

func buildCreateParam(opts interface{}, d *schema.ResourceData, nameMap *map[string]string) ([]string, error) {
	var f funcSkipOpt

	f = func(optn string, jsonTags []string, tag reflect.StructTag) bool {
		if getParamTag("required", tag) == "true" {
			return false
		}

		// For Create operation, it should not pass the parameter in the request, which match all the following situations.
		// a. Parameter is optional, which means it is not set 'required' in the tag.
		// b. Parameter's default value is allowed, which menas it is not set 'omitempty' in the tag of 'json'. The default value is like this, '0' for int and 'false' for bool
		// c. Parameter is not set default value in schema. It did not find a way to check whether it was set default value in the schema. so, add a new tag of "no_default" to mark it.
		// d. User did not set that parameter in the configuration file, which means the return value of 'hasFilledOpt' is false.
		if (len(jsonTags) == 1 || jsonTags[1] == "-") && getParamTag("no_default", tag) == "y" && !hasFilledOpt(d, optn) {
			return true
		}

		return false
	}

	e := &exchangeParam{
		OptNameMap: nameMap,
		SkipOpt:    f,
	}
	return e.BuildCUParam(opts, d)
}

func buildUpdateParam(opts interface{}, d *schema.ResourceData, nameMap *map[string]string) ([]string, error) {
	hasUpdatedItems := false

	var f funcSkipOpt
	f = func(optn string, jsonTags []string, tag reflect.StructTag) bool {
		v := d.HasChange(optn)
		if !hasUpdatedItems && v {
			hasUpdatedItems = true
		}

		// filter the unchanged parameters
		return !v
	}

	e := &exchangeParam{
		OptNameMap: nameMap,
		SkipOpt:    f,
	}
	notPassFileds, err := e.BuildCUParam(opts, d)
	if err != nil {
		return notPassFileds, err
	}
	if !hasUpdatedItems {
		return notPassFileds, fmt.Errorf("no changes happened")
	}
	return notPassFileds, nil
}

func refreshResourceData(resource interface{}, d *schema.ResourceData, nameMap *map[string]string) error {
	e := &exchangeParam{
		OptNameMap: nameMap,
	}
	return e.BuildResourceData(resource, d)
}
