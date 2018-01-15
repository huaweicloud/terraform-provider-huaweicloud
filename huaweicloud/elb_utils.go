package huaweicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
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

func isELBResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(gophercloud.ErrDefault404)
	return ok
}

func hasFilledParam(d *schema.ResourceData, param string) bool {
	_, b := d.GetOkExists(param)
	return b
}

func buildELBCreateParam(opts interface{}, d *schema.ResourceData) error {
	return buildELBCUParam(opts, d, false)
}

func buildELBUpdateParam(opts interface{}, d *schema.ResourceData) error {
	return buildELBCUParam(opts, d, true)
}

func buildELBCUParam(opts interface{}, d *schema.ResourceData, buildUpdate bool) error {
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

	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" {
			return fmt.Errorf("can not convert for item %v: without of json tag", v)
		}
		param := strings.Split(tag, ",")[0]
		if buildUpdate && !d.HasChange(param) {
			continue
		}

		if d.Get(param) == nil {
			continue
		}

		switch v.Kind() {
		case reflect.String:
			v.SetString(d.Get(param).(string))
		case reflect.Int:
			v.SetInt(int64(d.Get(param).(int)))
		case reflect.Bool:
			v.SetBool(d.Get(param).(bool))
		case reflect.Slice:
			s := d.Get(param).([]interface{})

			switch v.Type().Elem().Kind() {
			case reflect.String:
				t := make([]string, len(s))
				for i, iv := range s {
					t[i] = iv.(string)
				}
				v.Set(reflect.ValueOf(t))
			default:
				return fmt.Errorf("unknown type of item %v: %v", v, v.Type().Elem().Kind())
			}

		default:
			return fmt.Errorf("unknown type of item %v: %v", v, v.Kind())
		}
	}
	return nil
}

func refreshResourceData(resource interface{}, d *schema.ResourceData) error {
	b, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("refreshResourceData:: marshal failed:%v", err)
	}

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
		param := strings.Split(tag, ",")[0]
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
