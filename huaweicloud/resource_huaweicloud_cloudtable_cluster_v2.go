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
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceCloudtableClusterV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudtableClusterV2Create,
		Read:   resourceCloudtableClusterV2Read,
		Delete: resourceCloudtableClusterV2Delete,

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

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"rs_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enable_iam_auth": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"lemon_num": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"opentsdb_num": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hbase_public_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"lemon_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"open_tsdb_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"opentsdb_public_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_quota": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"used_storage_size": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zookeeper_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudtableClusterV2UserInputParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"terraform_resource_data": d,
		"availability_zone":       d.Get("availability_zone"),
		"enable_iam_auth":         d.Get("enable_iam_auth"),
		"lemon_num":               d.Get("lemon_num"),
		"name":                    d.Get("name"),
		"opentsdb_num":            d.Get("opentsdb_num"),
		"rs_num":                  d.Get("rs_num"),
		"security_group_id":       d.Get("security_group_id"),
		"storage_type":            d.Get("storage_type"),
		"subnet_id":               d.Get("subnet_id"),
		"tags":                    d.Get("tags"),
		"vpc_id":                  d.Get("vpc_id"),
	}
}

func resourceCloudtableClusterV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CloudtableV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	opts := resourceCloudtableClusterV2UserInputParams(d)

	params, err := buildCloudtableClusterV2CreateParameters(opts, nil)
	if err != nil {
		return fmtp.Errorf("Error building the request body of api(create), err=%s", err)
	}
	r, err := sendCloudtableClusterV2CreateRequest(d, params, client)
	if err != nil {
		return fmtp.Errorf("Error creating CloudtableClusterV2, err=%s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	obj, err := asyncWaitCloudtableClusterV2Create(d, config, r, client, timeout)
	if err != nil {
		return err
	}
	id, err := navigateValue(obj, []string{"cluster_id"}, nil)
	if err != nil {
		return fmtp.Errorf("Error constructing id, err=%s", err)
	}
	d.SetId(id.(string))

	return resourceCloudtableClusterV2Read(d, meta)
}

func resourceCloudtableClusterV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CloudtableV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	res := make(map[string]interface{})

	v, err := sendCloudtableClusterV2ReadRequest(d, client)
	if err != nil {
		return err
	}
	res["read"] = fillCloudtableClusterV2ReadRespBody(v)

	return setCloudtableClusterV2Properties(d, res)
}

func resourceCloudtableClusterV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CloudtableV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	url, err := replaceVars(d, "clusters/{id}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	logp.Printf("[DEBUG] Deleting Cluster %q", d.Id())
	r := golangsdk.Result{}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:      successHTTPCodes,
		JSONBody:     nil,
		JSONResponse: nil,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	})
	if r.Err != nil {
		return fmtp.Errorf("Error deleting Cluster %q, err=%s", d.Id(), r.Err)
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

func buildCloudtableClusterV2CreateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	v, err := navigateValue(opts, []string{"enable_iam_auth"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["auth_mode"] = v
	}

	v, err = expandCloudtableClusterV2CreateDatastore(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["datastore"] = v
	}

	v, err = expandCloudtableClusterV2CreateEnableLemon(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["enable_lemon"] = v
	}

	v, err = expandCloudtableClusterV2CreateEnableOpenTSDB(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["enable_openTSDB"] = v
	}

	v, err = expandCloudtableClusterV2CreateInstance(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["instance"] = v
	}

	v, err = navigateValue(opts, []string{"name"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["name"] = v
	}

	v, err = navigateValue(opts, []string{"storage_type"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["storage_type"] = v
	}

	v, err = expandCloudtableClusterV2CreateSysTags(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["sys_tags"] = v
	}

	v, err = navigateValue(opts, []string{"vpc_id"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["vpc_id"] = v
	}

	if len(params) == 0 {
		return params, nil
	}

	params = map[string]interface{}{"cluster": params}

	return params, nil
}

func expandCloudtableClusterV2CreateDatastore(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	req["type"] = "hbase"

	req["version"] = "1.0.6"

	return req, nil
}

func expandCloudtableClusterV2CreateEnableLemon(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"lemon_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if v1, ok := v.(int); ok && v1 > 0 {
		return true, nil
	}
	return false, nil
}

func expandCloudtableClusterV2CreateEnableOpenTSDB(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	v, err := navigateValue(d, []string{"opentsdb_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if v1, ok := v.(int); ok && v1 > 0 {
		return true, nil
	}
	return false, nil
}

func expandCloudtableClusterV2CreateInstance(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	v, err := navigateValue(d, []string{"availability_zone"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["availability_zone"] = v
	}

	v, err = navigateValue(d, []string{"rs_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["cu_num"] = v
	}

	v, err = navigateValue(d, []string{"lemon_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["lemon_num"] = v
	}

	v, err = expandCloudtableClusterV2CreateInstanceNics(d, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["nics"] = v
	}

	v, err = navigateValue(d, []string{"opentsdb_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["tsd_num"] = v
	}

	return req, nil
}

func expandCloudtableClusterV2CreateInstanceNics(d interface{}, arrayIndex map[string]int) (interface{}, error) {
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
			transformed["net_id"] = v
		}

		v, err = navigateValue(d, []string{"security_group_id"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["security_group_id"] = v
		}

		if len(transformed) > 0 {
			req = append(req, transformed)
		}
	}

	return req, nil
}

func expandCloudtableClusterV2CreateSysTags(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	val, err := navigateValue(d, []string{"tags"}, newArrayIndex)
	if err != nil {
		return nil, err
	}
	n := 0
	if val1, ok := val.([]interface{}); ok && len(val1) > 0 {
		n = len(val1)
	} else {
		return nil, nil
	}
	req := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		newArrayIndex["tags"] = i
		transformed := make(map[string]interface{})

		v, err := navigateValue(d, []string{"tags", "key"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["key"] = v
		}

		v, err = navigateValue(d, []string{"tags", "value"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["value"] = v
		}

		if len(transformed) > 0 {
			req = append(req, transformed)
		}
	}

	return req, nil
}

func sendCloudtableClusterV2CreateRequest(d *schema.ResourceData, params interface{},
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
		return nil, fmtp.Errorf("Error running api(create), err=%s", r.Err)
	}
	return r.Body, nil
}

func asyncWaitCloudtableClusterV2Create(d *schema.ResourceData, config *config.Config, result interface{},
	client *golangsdk.ServiceClient, timeout time.Duration) (interface{}, error) {

	data := make(map[string]interface{})
	pathParameters := map[string][]string{
		"cluster_id": []string{"cluster_id"},
	}
	for key, path := range pathParameters {
		value, err := navigateValue(result, path, nil)
		if err != nil {
			return nil, fmtp.Errorf("Error retrieving async operation path parameter, err=%s", err)
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

func sendCloudtableClusterV2ReadRequest(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
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
		return nil, fmtp.Errorf("Error running api(read) for resource(CloudtableClusterV2), err=%s", r.Err)
	}

	return r.Body, nil
}

func fillCloudtableClusterV2ReadRespBody(body interface{}) interface{} {
	result := make(map[string]interface{})
	val, ok := body.(map[string]interface{})
	if !ok {
		val = make(map[string]interface{})
	}

	if v, ok := val["actions"]; ok {
		result["actions"] = v
	} else {
		result["actions"] = nil
	}

	if v, ok := val["auth_mode"]; ok {
		result["auth_mode"] = v
	} else {
		result["auth_mode"] = nil
	}

	if v, ok := val["cluster_id"]; ok {
		result["cluster_id"] = v
	} else {
		result["cluster_id"] = nil
	}

	if v, ok := val["cluster_name"]; ok {
		result["cluster_name"] = v
	} else {
		result["cluster_name"] = nil
	}

	if v, ok := val["created"]; ok {
		result["created"] = v
	} else {
		result["created"] = nil
	}

	if v, ok := val["cu_num"]; ok {
		result["cu_num"] = v
	} else {
		result["cu_num"] = nil
	}

	if v, ok := val["enable_lemon"]; ok {
		result["enable_lemon"] = v
	} else {
		result["enable_lemon"] = nil
	}

	if v, ok := val["enable_openTSDB"]; ok {
		result["enable_openTSDB"] = v
	} else {
		result["enable_openTSDB"] = nil
	}

	if v, ok := val["hbase_public_endpoint"]; ok {
		result["hbase_public_endpoint"] = v
	} else {
		result["hbase_public_endpoint"] = nil
	}

	if v, ok := val["is_frozen"]; ok {
		result["is_frozen"] = v
	} else {
		result["is_frozen"] = nil
	}

	if v, ok := val["lemon_link"]; ok {
		result["lemon_link"] = v
	} else {
		result["lemon_link"] = nil
	}

	if v, ok := val["lemon_num"]; ok {
		result["lemon_num"] = v
	} else {
		result["lemon_num"] = nil
	}

	if v, ok := val["openTSDB_link"]; ok {
		result["openTSDB_link"] = v
	} else {
		result["openTSDB_link"] = nil
	}

	if v, ok := val["security_group_id"]; ok {
		result["security_group_id"] = v
	} else {
		result["security_group_id"] = nil
	}

	if v, ok := val["status"]; ok {
		result["status"] = v
	} else {
		result["status"] = nil
	}

	if v, ok := val["storage_quota"]; ok {
		result["storage_quota"] = v
	} else {
		result["storage_quota"] = nil
	}

	if v, ok := val["storage_type"]; ok {
		result["storage_type"] = v
	} else {
		result["storage_type"] = nil
	}

	if v, ok := val["sub_net_id"]; ok {
		result["sub_net_id"] = v
	} else {
		result["sub_net_id"] = nil
	}

	if v, ok := val["tsd_num"]; ok {
		result["tsd_num"] = v
	} else {
		result["tsd_num"] = nil
	}

	if v, ok := val["tsd_public_endpoint"]; ok {
		result["tsd_public_endpoint"] = v
	} else {
		result["tsd_public_endpoint"] = nil
	}

	if v, ok := val["updated"]; ok {
		result["updated"] = v
	} else {
		result["updated"] = nil
	}

	if v, ok := val["used_storage_size"]; ok {
		result["used_storage_size"] = v
	} else {
		result["used_storage_size"] = nil
	}

	if v, ok := val["vpc_id"]; ok {
		result["vpc_id"] = v
	} else {
		result["vpc_id"] = nil
	}

	if v, ok := val["zookeeper_link"]; ok {
		result["zookeeper_link"] = v
	} else {
		result["zookeeper_link"] = nil
	}

	return result
}

func setCloudtableClusterV2Properties(d *schema.ResourceData, response map[string]interface{}) error {
	opts := resourceCloudtableClusterV2UserInputParams(d)

	v, err := navigateValue(response, []string{"read", "created"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:created, err: %s", err)
	}
	if err = d.Set("created", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:created, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "auth_mode"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:enable_iam_auth, err: %s", err)
	}
	if err = d.Set("enable_iam_auth", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:enable_iam_auth, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "hbase_public_endpoint"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:hbase_public_endpoint, err: %s", err)
	}
	if err = d.Set("hbase_public_endpoint", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:hbase_public_endpoint, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "lemon_link"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:lemon_link, err: %s", err)
	}
	if err = d.Set("lemon_link", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:lemon_link, err: %s", err)
	}

	v, _ = opts["lemon_num"]
	v, err = flattenCloudtableClusterV2LemonNum(response, nil, v)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:lemon_num, err: %s", err)
	}
	if err = d.Set("lemon_num", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:lemon_num, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "cluster_name"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:name, err: %s", err)
	}
	if err = d.Set("name", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:name, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "openTSDB_link"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:open_tsdb_link, err: %s", err)
	}
	if err = d.Set("open_tsdb_link", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:open_tsdb_link, err: %s", err)
	}

	v, _ = opts["opentsdb_num"]
	v, err = flattenCloudtableClusterV2OpentsdbNum(response, nil, v)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:opentsdb_num, err: %s", err)
	}
	if err = d.Set("opentsdb_num", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:opentsdb_num, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "tsd_public_endpoint"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:opentsdb_public_endpoint, err: %s", err)
	}
	if err = d.Set("opentsdb_public_endpoint", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:opentsdb_public_endpoint, err: %s", err)
	}

	v, _ = opts["rs_num"]
	v, err = flattenCloudtableClusterV2RsNum(response, nil, v)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:rs_num, err: %s", err)
	}
	if err = d.Set("rs_num", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:rs_num, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "security_group_id"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:security_group_id, err: %s", err)
	}
	if err = d.Set("security_group_id", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:security_group_id, err: %s", err)
	}

	v, _ = opts["storage_quota"]
	v, err = flattenCloudtableClusterV2StorageQuota(response, nil, v)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:storage_quota, err: %s", err)
	}
	if err = d.Set("storage_quota", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:storage_quota, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "storage_type"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:storage_type, err: %s", err)
	}
	if err = d.Set("storage_type", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:storage_type, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "sub_net_id"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:subnet_id, err: %s", err)
	}
	if err = d.Set("subnet_id", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:subnet_id, err: %s", err)
	}

	v, _ = opts["used_storage_size"]
	v, err = flattenCloudtableClusterV2UsedStorageSize(response, nil, v)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:used_storage_size, err: %s", err)
	}
	if err = d.Set("used_storage_size", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:used_storage_size, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "vpc_id"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:vpc_id, err: %s", err)
	}
	if err = d.Set("vpc_id", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:vpc_id, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "zookeeper_link"}, nil)
	if err != nil {
		return fmtp.Errorf("Error reading Cluster:zookeeper_link, err: %s", err)
	}
	if err = d.Set("zookeeper_link", v); err != nil {
		return fmtp.Errorf("Error setting Cluster:zookeeper_link, err: %s", err)
	}

	return nil
}

func flattenCloudtableClusterV2LemonNum(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	v, err := navigateValue(d, []string{"read", "lemon_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	return convertToInt(v)
}

func flattenCloudtableClusterV2OpentsdbNum(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	v, err := navigateValue(d, []string{"read", "tsd_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	return convertToInt(v)
}

func flattenCloudtableClusterV2RsNum(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	v, err := navigateValue(d, []string{"read", "cu_num"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	return convertToInt(v)
}

func flattenCloudtableClusterV2StorageQuota(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	v, err := navigateValue(d, []string{"read", "storage_quota"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	return convertToStr(v) + " GB", nil
}

func flattenCloudtableClusterV2UsedStorageSize(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	v, err := navigateValue(d, []string{"read", "used_storage_size"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	return convertToStr(v) + " GB", nil
}
