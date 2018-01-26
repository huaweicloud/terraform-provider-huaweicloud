package huaweicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elb"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elb/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elb/loadbalancers"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func waitForELBJobSuccess(networkingClient *gophercloud.ServiceClient1, j *elb.Job, timeout time.Duration) (*elb.JobInfo, error) {
	jobId := j.JobId
	target := "SUCCESS"
	pending := []string{"INIT", "RUNNING"}

	log.Printf("[DEBUG] Waiting for elb job %s to become %s.", jobId, target)

	ji, err := waitForELBResource(networkingClient, "job", j.Uri, target, pending, timeout, getELBJobInfo)
	if err == nil {
		return ji.(*elb.JobInfo), nil
	}
	return nil, err
}

func getELBJobInfo(networkingClient *gophercloud.ServiceClient1, uri string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		info, err := elb.QueryJobInfo(networkingClient, uri).Extract()
		if err != nil {
			return nil, "", err
		}

		return info, info.Status, nil
	}
}

func waitForELBLoadBalancerActive(networkingClient *gophercloud.ServiceClient1, id string, timeout time.Duration) error {
	target := "ACTIVE"

	log.Printf("[DEBUG] Waiting for elb %s to become %s.", id, target)

	_, err := waitForELBResource(networkingClient, "loadbalancer", id, target, []string{"PENDING_CREATE"}, timeout, getELBLoadBalancer)
	return err
}

func getELBLoadBalancer(networkingClient *gophercloud.ServiceClient1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := loadbalancers.Get(networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}
		return lb, lb.Status, nil
	}
}

func waitForELBListenerActive(networkingClient *gophercloud.ServiceClient1, id string, timeout time.Duration) error {
	target := "ACTIVE"

	log.Printf("[DEBUG] Waiting for elb-listener %s to become %s.", id, target)

	_, err := waitForELBResource(networkingClient, "listener", id, target, []string{"PENDING_CREATE"}, timeout, getELBListener)
	return err
}

func getELBListener(networkingClient *gophercloud.ServiceClient1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		l, err := listeners.Get(networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}
		return l, l.Status, nil
	}
}

type getELBResource func(networkingClient *gophercloud.ServiceClient1, id string) resource.StateRefreshFunc

func waitForELBResource(networkingClient *gophercloud.ServiceClient1, name string, id string, target string, pending []string, timeout time.Duration, f getELBResource) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    f(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	o, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			return nil, fmt.Errorf("Error: elb %s %s not found: %s", name, id, err)
		}
		return nil, fmt.Errorf("Error waiting for elb %s %s to become %s: %s", name, id, target, err)
	}

	return o, nil
}

func chooseELBClient(d *schema.ResourceData, config *Config) (*gophercloud.ServiceClient1, error) {
	return config.loadElasticLoadBalancerClient(GetRegion(d, config))
}

func isResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(gophercloud.ErrDefault404)
	return ok
}

// The result may be not correct when the type of param is string and user config it to 'param=""'
// but, there is no other way.
func hasFilledParam(d *schema.ResourceData, param string) bool {
	_, b := d.GetOkExists(param)
	return b
}

func getParamTag(key string, tag reflect.StructTag) string {
	v, ok := tag.Lookup(key)
	if ok {
		return v
	}
	return "tag_not_set"
}

func get_param_name_from_tag(name string) string {
	return strings.ToLower(name)
}

type skipParamParsing func(param string, jsonTags []string, tag reflect.StructTag) bool

func buildCreateParam(opts interface{}, d *schema.ResourceData) (error, []string) {
	var not_pass_params []string
	var h skipParamParsing
	h = func(param string, jsonTags []string, tag reflect.StructTag) bool {
		if getParamTag("required", tag) == "true" {
			return false
		}

		// For Create operation, it should not pass the parameter in the request, which match all the following situations.
		// a. Parameter is optional, which means it is not set 'required' in the tag.
		// b. Parameter's default value is allowed, which menas it is not set 'omitempty' in the tag of 'json'. The default value is like this, '0' for int and 'false' for bool
		// c. Parameter is not set default value in schema. It did not find a way to check whether it was set default value in the schema. so, add a new tag of "no_default" to mark it.
		// d. User did not set that parameter in the configuration file, which means the return value of 'hasFilledParam' is false.
		if (len(jsonTags) == 1 || jsonTags[1] == "-") && getParamTag("no_default", tag) == "y" && !hasFilledParam(d, param) {
			not_pass_params = append(not_pass_params, param)
			return true
		}

		return false
	}

	return buildCUParam(opts, d, h), not_pass_params
}

func buildUpdateParam(opts interface{}, d *schema.ResourceData) (error, []string) {
	hasUpdatedItems := false
	var not_pass_params []string

	var h skipParamParsing
	h = func(param string, jsonTags []string, tag reflect.StructTag) bool {
		// filter the unchanged parameters
		if !d.HasChange(param) {
			not_pass_params = append(not_pass_params, param)
			return true
		} else if !hasUpdatedItems {
			hasUpdatedItems = true
		}

		return false
	}
	err := buildCUParam(opts, d, h)
	if err != nil {
		return err, not_pass_params
	}
	if !hasUpdatedItems {
		return fmt.Errorf("no changes happened"), not_pass_params
	}
	return nil, not_pass_params
}

func buildCUParam(opts interface{}, d *schema.ResourceData, skip skipParamParsing) error {
	optsValue := reflect.ValueOf(opts)
	if optsValue.Kind() != reflect.Ptr {
		return fmt.Errorf("parameter of opts should be a pointer")
	}
	optsValue = optsValue.Elem()
	if optsValue.Kind() != reflect.Struct {
		return fmt.Errorf("parameter must be a pointer to a struct")
	}

	optsType := reflect.TypeOf(opts)
	optsType = optsType.Elem()
	value := make(map[string]interface{})

	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)
		tag := getParamTag("json", f.Tag)
		if tag == "" {
			return fmt.Errorf("can not convert for item %v: without of json tag", v)
		}
		tags := strings.Split(tag, ",")
		param := get_param_name_from_tag(tags[0])
		// Only check the parameters in top struct.
		// If the parameters in sub-struct need skip, it will miss.
		// If it happens, need refactor here.
		if skip(param, tags, f.Tag) {
			continue
		}
		pv := d.Get(param)
		if pv == nil {
			log.Printf("[DEBUG] param:%s is not set", param)
			continue
		}
		value[param] = pv
	}
	if len(value) == 0 {
		log.Printf("[WARN]no parameter was set")
		return nil
	}
	return buildStruct(&optsValue, optsType, value)
}

func buildStruct(optsValue *reflect.Value, optsType reflect.Type, value map[string]interface{}) error {
	log.Printf("[DEBUG] buildStruct:: optsValue=%v, optsType=%v, value=%#v\n", optsValue, optsType, value)

	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" {
			return fmt.Errorf("can not convert for item %v: without of json tag", v)
		}
		param := get_param_name_from_tag(strings.Split(tag, ",")[0])
		log.Printf("[DEBUG] buildStruct:: convert for param:%s", param)
		if _, e := value[param]; !e {
			log.Printf("[DEBUG] param:%s was not supplied", param)
			continue
		}

		switch v.Kind() {
		case reflect.String:
			v.SetString(value[param].(string))
		case reflect.Int:
			v.SetInt(int64(value[param].(int)))
		case reflect.Int64:
			v.SetInt(value[param].(int64))
		case reflect.Bool:
			v.SetBool(value[param].(bool))
		case reflect.Slice:
			s := value[param].([]interface{})

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
					e := buildStruct(&rv, f.Type.Elem(), iv.(map[string]interface{}))
					if e != nil {
						return e
					}
				}
				v.Set(t)

			default:
				return fmt.Errorf("unknown type of item %v: %v", v, v.Type().Elem().Kind())
			}
		case reflect.Struct:
			log.Printf("[DEBUG] buildStruct:: convert struct for param %s: %#v", param, value[param])
			var p map[string]interface{}

			// If the type of parameter is Struct, then the corresponding type in Schema is TypeList
			v0, ok := value[param].([]interface{})
			if ok {
				p, ok = v0[0].(map[string]interface{})
			} else {
				p, ok = value[param].(map[string]interface{})
			}
			if !ok {
				return fmt.Errorf("can not convert to (map[string]interface{}) for param %s: %#v", param, value[param])
			}

			e := buildStruct(&v, f.Type, p)
			if e != nil {
				return e
			}

		default:
			return fmt.Errorf("unknown type of item %v: %v", v, v.Kind())
		}
	}
	return nil
}

func changeKeyToLowercase(b []byte) {
	m := regexp.MustCompile(`"([a-z0-9_]*[A-Z]+[a-z0-9_]*)+":`)
	for {
		bs := fmt.Sprintf("%s", b)
		index := m.FindStringIndex(bs)
		if index == nil {
			break
		}
		for i := index[0] + 1; i < index[1]-1; i++ {
			if b[i] > 0x40 && b[i] < 0x5B {
				b[i] = b[i] + 0x20
			}
		}
	}
}

func refreshResourceData(resource interface{}, d *schema.ResourceData) error {
	b, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("refreshResourceData:: marshal failed:%v", err)
	}
	changeKeyToLowercase(b)

	p := make(map[string]interface{})
	err = json.Unmarshal(b, &p)
	if err != nil {
		return fmt.Errorf("refreshResourceData:: unmarshal failed:%v", err)
	}
	log.Printf("[DEBUG]refreshResourceData:: raw data = %#v\n", p)
	return readStruct(resource, p, d)
}

func readStruct(resource interface{}, value map[string]interface{}, d *schema.ResourceData) error {

	optsValue := reflect.ValueOf(resource)
	if optsValue.Kind() == reflect.Ptr {
		optsValue = optsValue.Elem()
	}

	optsType := reflect.TypeOf(resource)
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
		param := get_param_name_from_tag(strings.Split(tag, ",")[0])
		if param == "id" {
			continue
		}
		log.Printf("[DEBUG readStruct:: convert for param:%s", param)

		switch v.Kind() {
		default:
			e := d.Set(param, value[param])
			if e != nil {
				return e
			}
		case reflect.Struct:
			//The corresponding schema of Struct is TypeList in Terrafrom
			e := d.Set(param, []interface{}{value[param]})
			if e != nil {
				return e
			}
		}
	}
	return nil
}
