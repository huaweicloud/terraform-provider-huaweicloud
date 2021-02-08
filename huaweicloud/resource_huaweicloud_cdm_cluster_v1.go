// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package huaweicloud

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
)

func resourceCdmClusterV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCdmClusterV1Create,
		Read:   resourceCdmClusterV1Read,
		Delete: resourceCdmClusterV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"email": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"is_auto_off": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"phone_num": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"schedule_boot_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"schedule_off_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"publid_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCdmClusterV1UserInputParams(d *schema.ResourceData, config *Config) map[string]interface{} {
	return map[string]interface{}{
		"terraform_resource_data": d,
		"availability_zone":       d.Get("availability_zone"),
		"email":                   d.Get("email"),
		"enterprise_project_id":   GetEnterpriseProjectID(d, config),
		"flavor_id":               d.Get("flavor_id"),
		"is_auto_off":             d.Get("is_auto_off"),
		"name":                    d.Get("name"),
		"phone_num":               d.Get("phone_num"),
		"schedule_boot_time":      d.Get("schedule_boot_time"),
		"schedule_off_time":       d.Get("schedule_off_time"),
		"security_group_id":       d.Get("security_group_id"),
		"subnet_id":               d.Get("subnet_id"),
		"version":                 d.Get("version"),
		"vpc_id":                  d.Get("vpc_id"),
	}
}

func resourceCdmClusterV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cdmV11Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	opts := resourceCdmClusterV1UserInputParams(d, config)
	params, err := buildCdmClusterV1CreateParameters(opts, nil)
	if err != nil {
		return fmt.Errorf("Error building the request body of api(create), err=%s", err)
	}
	r, err := sendCdmClusterV1CreateRequest(d, params, client)
	if err != nil {
		return fmt.Errorf("Error creating CdmClusterV1, err=%s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	obj, err := asyncWaitCdmClusterV1Create(d, config, r, client, timeout)
	if err != nil {
		return err
	}
	id, err := navigateValue(obj, []string{"id"}, nil)
	if err != nil {
		return fmt.Errorf("Error constructing id, err=%s", err)
	}
	d.SetId(id.(string))

	return resourceCdmClusterV1Read(d, meta)
}

func resourceCdmClusterV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cdmV11Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	res := make(map[string]interface{})

	v, err := sendCdmClusterV1ReadRequest(d, client)
	if err != nil {
		return err
	}
	res["read"] = fillCdmClusterV1ReadRespBody(v)

	return setCdmClusterV1Properties(d, config, res)
}

func resourceCdmClusterV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cdmV11Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	url, err := replaceVars(d, "clusters/{id}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	opts := resourceCdmClusterV1UserInputParams(d, config)
	params, err := buildCdmClusterV1DeleteParameters(opts, nil)
	if err != nil {
		return fmt.Errorf("Error building the request body of api(delete), err=%s", err)
	}

	log.Printf("[DEBUG] Deleting Cluster %q", d.Id())
	r := golangsdk.Result{}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:      successHTTPCodes,
		JSONBody:     params,
		JSONResponse: nil,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	})
	if r.Err != nil {
		return fmt.Errorf("Error deleting Cluster %q, err=%s", d.Id(), r.Err)
	}

	_, err = waitToFinish(
		[]string{"Done"}, []string{"Pending"},
		d.Timeout(schema.TimeoutCreate),
		1*time.Second,
		func() (interface{}, string, error) {
			_, err := client.Get(url, nil, &golangsdk.RequestOpts{
				MoreHeaders: map[string]string{
					"Content-Type": "application/json",
					"X-Language":   "en-us",
				}})
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return true, "Done", nil
				}
				return nil, "", nil
			}
			return true, "Pending", nil
		},
	)
	return err
}

func buildCdmClusterV1CreateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	v, err := expandCdmClusterV1CreateAutoRemind(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["autoRemind"] = v
	}

	v, err = expandCdmClusterV1CreateCluster(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["cluster"] = v
	}

	v, err = expandCdmClusterV1CreateEmail(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["email"] = v
	}

	v, err = expandCdmClusterV1CreatePhoneNum(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["phoneNum"] = v
	}

	return params, nil
}

func expandCdmClusterV1CreateCluster(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	v, err := expandCdmClusterV1CreateClusterDatastore(d, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["datastore"] = v
	}

	v, err = expandCdmClusterV1CreateClusterInstances(d, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["instances"] = v
	}

	v, err = navigateValue(d, []string{"is_auto_off"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["isAutoOff"] = v
	}

	v, err = expandCdmClusterV1CreateClusterIsScheduleBootOff(d, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["isScheduleBootOff"] = v
	}

	v, err = navigateValue(d, []string{"name"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["name"] = v
	}

	v, err = navigateValue(d, []string{"schedule_boot_time"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["scheduleBootTime"] = v
	}

	v, err = navigateValue(d, []string{"schedule_off_time"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["scheduleOffTime"] = v
	}

	v, err = expandCdmClusterV1CreateClusterSysTags(d, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["sys_tags"] = v
	}

	v, err = navigateValue(d, []string{"vpc_id"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["vpcId"] = v
	}

	return req, nil
}

func expandCdmClusterV1CreateClusterDatastore(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	req["type"] = "cdm"

	v, err := navigateValue(d, []string{"version"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["version"] = v
	}

	return req, nil
}

func expandCdmClusterV1CreateClusterInstances(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	n := 1
	req := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		transformed := make(map[string]interface{})

		v, err := navigateValue(d, []string{"availability_zone"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["availability_zone"] = v
		}

		v, err = navigateValue(d, []string{"flavor_id"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["flavorRef"] = v
		}

		v, err = expandCdmClusterV1CreateClusterInstancesNics(d, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["nics"] = v
		}

		transformed["type"] = "cdm"

		if len(transformed) > 0 {
			req = append(req, transformed)
		}
	}

	return req, nil
}

func expandCdmClusterV1CreateClusterInstancesNics(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	n := 1
	req := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		transformed := make(map[string]interface{})

		v, err := navigateValue(d, []string{"subnet_id"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["net-id"] = v
		}

		v, err = navigateValue(d, []string{"security_group_id"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["securityGroupId"] = v
		}

		if len(transformed) > 0 {
			req = append(req, transformed)
		}
	}

	return req, nil
}

func expandCdmClusterV1CreateClusterSysTags(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"enterprise_project_id"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	sysTags := make([]interface{}, 1, 1)
	sysTags[0] = map[string]string{
		"key":   "_sys_enterprise_project_id",
		"value": v.(string),
	}
	return sysTags, nil
}

func expandCdmClusterV1CreateEmail(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"email"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if v1, ok := v.([]interface{}); ok && len(v1) > 0 {
		v2 := make([]string, len(v1), len(v1))
		for i, j := range v1 {
			v2[i] = j.(string)
		}
		return strings.Join(v2, ","), nil
	}
	return "", nil
}

func expandCdmClusterV1CreatePhoneNum(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"phone_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if v1, ok := v.([]interface{}); ok && len(v1) > 0 {
		v2 := make([]string, len(v1), len(v1))
		for i, j := range v1 {
			v2[i] = j.(string)
		}
		return strings.Join(v2, ","), nil
	}
	return "", nil
}

func sendCdmClusterV1CreateRequest(d *schema.ResourceData, params interface{},
	client *golangsdk.ServiceClient) (interface{}, error) {
	url := client.ServiceURL("clusters")

	r := golangsdk.Result{}
	_, r.Err = client.Post(url, params, &r.Body, &golangsdk.RequestOpts{
		OkCodes: successHTTPCodes,
		MoreHeaders: map[string]string{
			"X-Language": "en-us",
		},
	})
	if r.Err != nil {
		return nil, fmt.Errorf("Error running api(create), err=%s", r.Err)
	}
	return r.Body, nil
}

func asyncWaitCdmClusterV1Create(d *schema.ResourceData, config *Config, result interface{},
	client *golangsdk.ServiceClient, timeout time.Duration) (interface{}, error) {

	data := make(map[string]interface{})
	pathParameters := map[string][]string{
		"cluster_id": []string{"id"},
	}
	for key, path := range pathParameters {
		value, err := navigateValue(result, path, nil)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving async operation path parameter, err=%s", err)
		}
		data[key] = value
	}

	url, err := replaceVars(d, "clusters/{cluster_id}", data)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	return waitToFinish(
		[]string{"200"},
		[]string{"100"},
		timeout, 1*time.Second,
		func() (interface{}, string, error) {
			r := golangsdk.Result{}
			_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
				MoreHeaders: map[string]string{
					"Content-Type": "application/json",
					"X-Language":   "en-us",
				}})
			if r.Err != nil {
				return nil, "", nil
			}

			status, err := navigateValue(r.Body, []string{"status"}, nil)
			if err != nil {
				return nil, "", nil
			}
			return r.Body, status.(string), nil
		},
	)
}

func buildCdmClusterV1DeleteParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	params["keepLastManualBackup"] = 0

	return params, nil
}

func sendCdmClusterV1ReadRequest(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	url, err := replaceVars(d, "clusters/{id}", nil)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		}})
	if r.Err != nil {
		return nil, fmt.Errorf("Error running api(read) for resource(CdmClusterV1), err=%s", r.Err)
	}

	return r.Body, nil
}

func fillCdmClusterV1ReadRespBody(body interface{}) interface{} {
	result := make(map[string]interface{})
	val, ok := body.(map[string]interface{})
	if !ok {
		val = make(map[string]interface{})
	}

	if v, ok := val["actionProgress"]; ok {
		result["actionProgress"] = v
	} else {
		result["actionProgress"] = nil
	}

	if v, ok := val["actions"]; ok {
		result["actions"] = v
	} else {
		result["actions"] = nil
	}

	if v, ok := val["azName"]; ok {
		result["azName"] = v
	} else {
		result["azName"] = nil
	}

	if v, ok := val["clusterMode"]; ok {
		result["clusterMode"] = v
	} else {
		result["clusterMode"] = nil
	}

	if v, ok := val["config_status"]; ok {
		result["config_status"] = v
	} else {
		result["config_status"] = nil
	}

	if v, ok := val["created"]; ok {
		result["created"] = v
	} else {
		result["created"] = nil
	}

	if v, ok := val["customerConfig"]; ok {
		result["customerConfig"] = fillCdmClusterV1ReadRespCustomerConfig(v)
	} else {
		result["customerConfig"] = nil
	}

	if v, ok := val["datastore"]; ok {
		result["datastore"] = fillCdmClusterV1ReadRespDatastore(v)
	} else {
		result["datastore"] = nil
	}

	if v, ok := val["dbuser"]; ok {
		result["dbuser"] = v
	} else {
		result["dbuser"] = nil
	}

	if v, ok := val["eipId"]; ok {
		result["eipId"] = v
	} else {
		result["eipId"] = nil
	}

	if v, ok := val["endpointDomainName"]; ok {
		result["endpointDomainName"] = v
	} else {
		result["endpointDomainName"] = nil
	}

	if v, ok := val["flavorName"]; ok {
		result["flavorName"] = v
	} else {
		result["flavorName"] = nil
	}

	if v, ok := val["id"]; ok {
		result["id"] = v
	} else {
		result["id"] = nil
	}

	if v, ok := val["instances"]; ok {
		result["instances"] = fillCdmClusterV1ReadRespInstances(v)
	} else {
		result["instances"] = nil
	}

	if v, ok := val["isAutoOff"]; ok {
		result["isAutoOff"] = v
	} else {
		result["isAutoOff"] = nil
	}

	if v, ok := val["isFrozen"]; ok {
		result["isFrozen"] = v
	} else {
		result["isFrozen"] = nil
	}

	if v, ok := val["isScheduleBootOff"]; ok {
		result["isScheduleBootOff"] = v
	} else {
		result["isScheduleBootOff"] = nil
	}

	if v, ok := val["links"]; ok {
		result["links"] = fillCdmClusterV1ReadRespLinks(v)
	} else {
		result["links"] = nil
	}

	if v, ok := val["maintainWindow"]; ok {
		result["maintainWindow"] = fillCdmClusterV1ReadRespMaintainWindow(v)
	} else {
		result["maintainWindow"] = nil
	}

	if v, ok := val["name"]; ok {
		result["name"] = v
	} else {
		result["name"] = nil
	}

	if v, ok := val["publicEndpoint"]; ok {
		result["publicEndpoint"] = v
	} else {
		result["publicEndpoint"] = nil
	}

	if v, ok := val["publicEndpointDomainName"]; ok {
		result["publicEndpointDomainName"] = v
	} else {
		result["publicEndpointDomainName"] = nil
	}

	if v, ok := val["publicEndpointStatus"]; ok {
		result["publicEndpointStatus"] = fillCdmClusterV1ReadRespPublicEndpointStatus(v)
	} else {
		result["publicEndpointStatus"] = nil
	}

	if v, ok := val["recentEvent"]; ok {
		result["recentEvent"] = v
	} else {
		result["recentEvent"] = nil
	}

	if v, ok := val["status"]; ok {
		result["status"] = v
	} else {
		result["status"] = nil
	}

	if v, ok := val["statusDetail"]; ok {
		result["statusDetail"] = v
	} else {
		result["statusDetail"] = nil
	}

	if v, ok := val["task"]; ok {
		result["task"] = fillCdmClusterV1ReadRespTask(v)
	} else {
		result["task"] = nil
	}

	if v, ok := val["updated"]; ok {
		result["updated"] = v
	} else {
		result["updated"] = nil
	}

	return result
}

func fillCdmClusterV1ReadRespCustomerConfig(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.(map[string]interface{})
	if !ok {
		value1 = make(map[string]interface{})
	}
	result := make(map[string]interface{})

	if v, ok := value1["clusterName"]; ok {
		result["clusterName"] = v
	} else {
		result["clusterName"] = nil
	}

	if v, ok := value1["failureRemind"]; ok {
		result["failureRemind"] = v
	} else {
		result["failureRemind"] = nil
	}

	if v, ok := value1["localDisk"]; ok {
		result["localDisk"] = v
	} else {
		result["localDisk"] = nil
	}

	if v, ok := value1["serviceProvider"]; ok {
		result["serviceProvider"] = v
	} else {
		result["serviceProvider"] = nil
	}

	if v, ok := value1["ssl"]; ok {
		result["ssl"] = v
	} else {
		result["ssl"] = nil
	}

	return result
}

func fillCdmClusterV1ReadRespDatastore(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.(map[string]interface{})
	if !ok {
		value1 = make(map[string]interface{})
	}
	result := make(map[string]interface{})

	if v, ok := value1["type"]; ok {
		result["type"] = v
	} else {
		result["type"] = nil
	}

	if v, ok := value1["version"]; ok {
		result["version"] = v
	} else {
		result["version"] = nil
	}

	return result
}

func fillCdmClusterV1ReadRespInstances(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.([]interface{})
	if !ok || len(value1) == 0 {
		return nil
	}

	n := len(value1)
	result := make([]interface{}, n, n)
	for i := 0; i < n; i++ {
		val := make(map[string]interface{})
		item := value1[i].(map[string]interface{})

		if v, ok := item["config_status"]; ok {
			val["config_status"] = v
		} else {
			val["config_status"] = nil
		}

		if v, ok := item["flavor"]; ok {
			val["flavor"] = fillCdmClusterV1ReadRespInstancesFlavor(v)
		} else {
			val["flavor"] = nil
		}

		if v, ok := item["group"]; ok {
			val["group"] = v
		} else {
			val["group"] = nil
		}

		if v, ok := item["id"]; ok {
			val["id"] = v
		} else {
			val["id"] = nil
		}

		if v, ok := item["isFrozen"]; ok {
			val["isFrozen"] = v
		} else {
			val["isFrozen"] = nil
		}

		if v, ok := item["links_2"]; ok {
			val["links_2"] = fillCdmClusterV1ReadRespInstancesLinks2(v)
		} else {
			val["links_2"] = nil
		}

		if v, ok := item["managerIp"]; ok {
			val["managerIp"] = v
		} else {
			val["managerIp"] = nil
		}

		if v, ok := item["name"]; ok {
			val["name"] = v
		} else {
			val["name"] = nil
		}

		if v, ok := item["paramsGroupId"]; ok {
			val["paramsGroupId"] = v
		} else {
			val["paramsGroupId"] = nil
		}

		if v, ok := item["publicIp"]; ok {
			val["publicIp"] = v
		} else {
			val["publicIp"] = nil
		}

		if v, ok := item["role"]; ok {
			val["role"] = v
		} else {
			val["role"] = nil
		}

		if v, ok := item["shard_id"]; ok {
			val["shard_id"] = v
		} else {
			val["shard_id"] = nil
		}

		if v, ok := item["status"]; ok {
			val["status"] = v
		} else {
			val["status"] = nil
		}

		if v, ok := item["trafficIp"]; ok {
			val["trafficIp"] = v
		} else {
			val["trafficIp"] = nil
		}

		if v, ok := item["type"]; ok {
			val["type"] = v
		} else {
			val["type"] = nil
		}

		if v, ok := item["volume"]; ok {
			val["volume"] = fillCdmClusterV1ReadRespInstancesVolume(v)
		} else {
			val["volume"] = nil
		}

		result[i] = val
	}

	return result
}

func fillCdmClusterV1ReadRespInstancesFlavor(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.(map[string]interface{})
	if !ok {
		value1 = make(map[string]interface{})
	}
	result := make(map[string]interface{})

	if v, ok := value1["id"]; ok {
		result["id"] = v
	} else {
		result["id"] = nil
	}

	if v, ok := value1["links_1"]; ok {
		result["links_1"] = fillCdmClusterV1ReadRespInstancesFlavorLinks1(v)
	} else {
		result["links_1"] = nil
	}

	return result
}

func fillCdmClusterV1ReadRespInstancesFlavorLinks1(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.([]interface{})
	if !ok || len(value1) == 0 {
		return nil
	}

	n := len(value1)
	result := make([]interface{}, n, n)
	for i := 0; i < n; i++ {
		val := make(map[string]interface{})
		item := value1[i].(map[string]interface{})

		if v, ok := item["href"]; ok {
			val["href"] = v
		} else {
			val["href"] = nil
		}

		if v, ok := item["rel"]; ok {
			val["rel"] = v
		} else {
			val["rel"] = nil
		}

		result[i] = val
	}

	return result
}

func fillCdmClusterV1ReadRespInstancesLinks2(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.([]interface{})
	if !ok || len(value1) == 0 {
		return nil
	}

	n := len(value1)
	result := make([]interface{}, n, n)
	for i := 0; i < n; i++ {
		val := make(map[string]interface{})
		item := value1[i].(map[string]interface{})

		if v, ok := item["href"]; ok {
			val["href"] = v
		} else {
			val["href"] = nil
		}

		if v, ok := item["rel"]; ok {
			val["rel"] = v
		} else {
			val["rel"] = nil
		}

		result[i] = val
	}

	return result
}

func fillCdmClusterV1ReadRespInstancesVolume(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.(map[string]interface{})
	if !ok {
		value1 = make(map[string]interface{})
	}
	result := make(map[string]interface{})

	if v, ok := value1["size"]; ok {
		result["size"] = v
	} else {
		result["size"] = nil
	}

	if v, ok := value1["type"]; ok {
		result["type"] = v
	} else {
		result["type"] = nil
	}

	return result
}

func fillCdmClusterV1ReadRespLinks(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.([]interface{})
	if !ok || len(value1) == 0 {
		return nil
	}

	n := len(value1)
	result := make([]interface{}, n, n)
	for i := 0; i < n; i++ {
		val := make(map[string]interface{})
		item := value1[i].(map[string]interface{})

		if v, ok := item["href"]; ok {
			val["href"] = v
		} else {
			val["href"] = nil
		}

		if v, ok := item["rel"]; ok {
			val["rel"] = v
		} else {
			val["rel"] = nil
		}

		result[i] = val
	}

	return result
}

func fillCdmClusterV1ReadRespMaintainWindow(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.(map[string]interface{})
	if !ok {
		value1 = make(map[string]interface{})
	}
	result := make(map[string]interface{})

	if v, ok := value1["day"]; ok {
		result["day"] = v
	} else {
		result["day"] = nil
	}

	if v, ok := value1["endTime"]; ok {
		result["endTime"] = v
	} else {
		result["endTime"] = nil
	}

	if v, ok := value1["startTime"]; ok {
		result["startTime"] = v
	} else {
		result["startTime"] = nil
	}

	return result
}

func fillCdmClusterV1ReadRespPublicEndpointStatus(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.(map[string]interface{})
	if !ok {
		value1 = make(map[string]interface{})
	}
	result := make(map[string]interface{})

	if v, ok := value1["errorMessage"]; ok {
		result["errorMessage"] = v
	} else {
		result["errorMessage"] = nil
	}

	if v, ok := value1["status"]; ok {
		result["status"] = v
	} else {
		result["status"] = nil
	}

	return result
}

func fillCdmClusterV1ReadRespTask(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.(map[string]interface{})
	if !ok {
		value1 = make(map[string]interface{})
	}
	result := make(map[string]interface{})

	if v, ok := value1["description"]; ok {
		result["description"] = v
	} else {
		result["description"] = nil
	}

	if v, ok := value1["id"]; ok {
		result["id"] = v
	} else {
		result["id"] = nil
	}

	if v, ok := value1["name"]; ok {
		result["name"] = v
	} else {
		result["name"] = nil
	}

	return result
}

func setCdmClusterV1Properties(d *schema.ResourceData, config *Config, response map[string]interface{}) error {
	opts := resourceCdmClusterV1UserInputParams(d, config)

	v, err := navigateValue(response, []string{"read", "created"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Cluster:created, err: %s", err)
	}
	if err = d.Set("created", v); err != nil {
		return fmt.Errorf("Error setting Cluster:created, err: %s", err)
	}

	v, _ = opts["instances"]
	v, err = flattenCdmClusterV1Instances(response, nil, v)
	if err != nil {
		return fmt.Errorf("Error reading Cluster:instances, err: %s", err)
	}
	if err = d.Set("instances", v); err != nil {
		return fmt.Errorf("Error setting Cluster:instances, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "isAutoOff"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Cluster:is_auto_off, err: %s", err)
	}
	if err = d.Set("is_auto_off", v); err != nil {
		return fmt.Errorf("Error setting Cluster:is_auto_off, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "name"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Cluster:name, err: %s", err)
	}
	if err = d.Set("name", v); err != nil {
		return fmt.Errorf("Error setting Cluster:name, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "publicEndpoint"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Cluster:publid_ip, err: %s", err)
	}
	if err = d.Set("publid_ip", v); err != nil {
		return fmt.Errorf("Error setting Cluster:publid_ip, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "datastore", "version"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Cluster:version, err: %s", err)
	}
	if err = d.Set("version", v); err != nil {
		return fmt.Errorf("Error setting Cluster:version, err: %s", err)
	}

	return nil
}

func flattenCdmClusterV1Instances(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	n := 0
	hasInitValue := true
	result, ok := currentValue.([]interface{})
	if !ok || len(result) == 0 {
		v, err := navigateValue(d, []string{"read", "instances"}, arrayIndex)
		if err != nil {
			return nil, err
		}
		if v1, ok := v.([]interface{}); ok && len(v1) > 0 {
			n = len(v1)
		} else {
			return currentValue, nil
		}
		result = make([]interface{}, 0, n)
		hasInitValue = false
	} else {
		n = len(result)
	}

	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	for i := 0; i < n; i++ {
		newArrayIndex["read.instances"] = i

		var r map[string]interface{}
		if len(result) >= (i+1) && result[i] != nil {
			r = result[i].(map[string]interface{})
		} else {
			r = make(map[string]interface{})
		}

		v, err := navigateValue(d, []string{"read", "instances", "id"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Cluster:id, err: %s", err)
		}
		r["id"] = v

		v, err = navigateValue(d, []string{"read", "instances", "name"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Cluster:name, err: %s", err)
		}
		r["name"] = v

		v, err = navigateValue(d, []string{"read", "instances", "publicIp"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Cluster:public_ip, err: %s", err)
		}
		r["public_ip"] = v

		v, err = navigateValue(d, []string{"read", "instances", "role"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Cluster:role, err: %s", err)
		}
		r["role"] = v

		v, err = navigateValue(d, []string{"read", "instances", "trafficIp"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Cluster:traffic_ip, err: %s", err)
		}
		r["traffic_ip"] = v

		v, err = navigateValue(d, []string{"read", "instances", "type"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Cluster:type, err: %s", err)
		}
		r["type"] = v

		if len(result) >= (i + 1) {
			if result[i] == nil {
				result[i] = r
			}
		} else {
			for _, v := range r {
				if v != nil {
					result = append(result, r)
					break
				}
			}
		}
	}

	if hasInitValue || len(result) > 0 {
		return result, nil
	}
	return currentValue, nil
}
