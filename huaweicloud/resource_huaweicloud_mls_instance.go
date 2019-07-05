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
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
)

func resourceMlsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceMlsInstanceCreate,
		Read:   resourceMlsInstanceRead,
		Delete: resourceMlsInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"mrs_cluster": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"user_password": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"available_zone": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"public_ip": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bind_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"eip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"agency": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"current_task": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"inner_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMlsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.sdkClient(GetRegion(d, config), "mls", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	opts := make(map[string]interface{})
	agencyProp := d.Get("agency")
	e, err := isEmptyValue(reflect.ValueOf(agencyProp))
	if err != nil {
		return err
	}
	if !e {
		opts["agency"] = agencyProp
	}

	flavorProp := d.Get("flavor")
	e, err = isEmptyValue(reflect.ValueOf(flavorProp))
	if err != nil {
		return err
	}
	if !e {
		opts["flavorRef"] = flavorProp
	}

	mrsClusterProp, err := expandMlsInstanceMrsCluster(d.Get("mrs_cluster"))
	if err != nil {
		return err
	}
	e, err = isEmptyValue(reflect.ValueOf(mrsClusterProp))
	if err != nil {
		return err
	}
	if !e {
		opts["mrsCluster"] = mrsClusterProp
	}

	nameProp := d.Get("name")
	e, err = isEmptyValue(reflect.ValueOf(nameProp))
	if err != nil {
		return err
	}
	if !e {
		opts["name"] = nameProp
	}

	networkProp, err := expandMlsInstanceNetwork(d.Get("network"))
	if err != nil {
		return err
	}
	e, err = isEmptyValue(reflect.ValueOf(networkProp))
	if err != nil {
		return err
	}
	if !e {
		opts["network"] = networkProp
	}

	versionProp := d.Get("version")
	e, err = isEmptyValue(reflect.ValueOf(versionProp))
	if err != nil {
		return err
	}
	if !e {
		opts["version"] = versionProp
	}

	url, err := replaceVars(d, "instances", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	log.Printf("[DEBUG] Creating new Instance: %#v", opts)
	r := golangsdk.Result{}
	_, r.Err = client.Post(
		url,
		&map[string]interface{}{"instance": opts},
		&r.Body,
		&golangsdk.RequestOpts{OkCodes: successHTTPCodes})
	if r.Err != nil {
		return fmt.Errorf("Error creating Instance: %s", r.Err)
	}

	pathParameters := map[string][]string{
		"id": {"instance", "id"},
	}
	var data = make(map[string]interface{})
	for key, path := range pathParameters {
		value, err := navigateMap(r.Body, path)
		if err != nil {
			return fmt.Errorf("Error retrieving async operation path parameter: %s", err)
		}
		data[key] = value
	}
	url, err = replaceVars(d, "instances/{id}", data)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	obj, err := waitToFinish(
		[]string{"AVAILABLE"},
		[]string{"CREATING", "Pending"},
		d.Timeout(schema.TimeoutCreate),
		1*time.Second,
		func() (interface{}, string, error) {
			r := golangsdk.Result{}
			_, r.Err = client.Get(
				url, &r.Body,
				&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
			if r.Err != nil {
				return nil, "", nil
			}

			status, err := navigateMap(r.Body, []string{"instance", "status"})
			if err != nil {
				return nil, "", nil
			}
			return r.Body, status.(string), nil
		},
	)

	if err != nil {
		return err
	}
	id, err := navigateMap(obj, []string{"instance", "id"})
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id.(string))

	return resourceMlsInstanceRead(d, meta)
}

func resourceMlsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.sdkClient(GetRegion(d, config), "mls", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	url, err := replaceVars(d, "instances/{id}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	r := golangsdk.Result{}
	_, r.Err = client.Get(
		url, &r.Body,
		&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
	if r.Err != nil {
		return fmt.Errorf("Error reading %s: %s", fmt.Sprintf("MlsInstance %q", d.Id()), r.Err)
	}
	v, err := navigateMap(r.Body, []string{"instance"})
	if err != nil {
		return fmt.Errorf("Error reading %s: the result does not contain instance", fmt.Sprintf("MlsInstance %q", d.Id()))
	}
	res := v.(map[string]interface{})

	if v, ok := res["created"]; ok {
		if err := d.Set("created", v); err != nil {
			return fmt.Errorf("Error reading Instance:created, err: %s", err)
		}
	}

	if v, ok := res["currentTask"]; ok {
		if err := d.Set("current_task", v); err != nil {
			return fmt.Errorf("Error reading Instance:current_task, err: %s", err)
		}
	}

	if v, ok := res["innerEndPoint"]; ok {
		if err := d.Set("inner_endpoint", v); err != nil {
			return fmt.Errorf("Error reading Instance:inner_endpoint, err: %s", err)
		}
	}

	if v, ok := res["publicEndPoint"]; ok {
		if err := d.Set("public_endpoint", v); err != nil {
			return fmt.Errorf("Error reading Instance:public_endpoint, err: %s", err)
		}
	}

	if v, ok := res["status"]; ok {
		if err := d.Set("status", v); err != nil {
			return fmt.Errorf("Error reading Instance:status, err: %s", err)
		}
	}

	if v, ok := res["updated"]; ok {
		if err := d.Set("updated", v); err != nil {
			return fmt.Errorf("Error reading Instance:updated, err: %s", err)
		}
	}

	if v, ok := res["network"]; ok {
		networkProp, err := flattenMlsInstanceNetwork(v)
		if err != nil {
			return fmt.Errorf("Error reading Instance:network, err: %s", err)
		}
		if err := d.Set("network", networkProp); err != nil {
			return fmt.Errorf("Error reading Instance:network, err: %s", err)
		}
	}

	return nil
}

func resourceMlsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.sdkClient(GetRegion(d, config), "mls", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	url, err := replaceVars(d, "instances/{id}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	log.Printf("[DEBUG] Deleting Instance %q", d.Id())
	r := golangsdk.Result{}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:      successHTTPCodes,
		JSONResponse: &r.Body,
		MoreHeaders:  map[string]string{"Content-Type": "application/json"},
		JSONBody:     map[string]interface{}{},
	})
	if r.Err != nil {
		return fmt.Errorf("Error deleting Instance %q: %s", d.Id(), r.Err)
	}

	_, err = waitToFinish(
		[]string{"Done"}, []string{"Pending"},
		d.Timeout(schema.TimeoutDelete),
		1*time.Second,
		func() (interface{}, string, error) {
			resp, err := client.Get(
				url, nil,
				&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return resp, "Done", nil
				}
				return nil, "", nil
			}
			return resp, "Pending", nil
		},
	)
	return err
}

func flattenMlsInstanceNetwork(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	if val, ok := original["availableZone"]; ok {
		transformed["available_zone"] = val
	}

	if val, ok := original["subnetId"]; ok {
		transformed["network_id"] = val
	}

	if val, ok := original["publicIP"]; ok {
		publicIPProp, err := flattenMlsInstanceNetworkPublicIP(val)
		if err != nil {
			return nil, fmt.Errorf("Error reading network:public_ip, err: %s", err)
		}
		transformed["public_ip"] = publicIPProp
	}

	if val, ok := original["securityGroupId"]; ok {
		transformed["security_group_id"] = val
	}

	if val, ok := original["vpcId"]; ok {
		transformed["vpc_id"] = val
	}

	return []interface{}{transformed}, nil
}

func flattenMlsInstanceNetworkPublicIP(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	if val, ok := original["bindType"]; ok {
		transformed["bind_type"] = val
	}

	if val, ok := original["eipId"]; ok {
		transformed["eip_id"] = val
	}

	return []interface{}{transformed}, nil
}

func expandMlsInstanceMrsCluster(v interface{}) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	idProp := original["id"]
	e, err := isEmptyValue(reflect.ValueOf(idProp))
	if err != nil {
		return nil, err
	}
	if !e {
		transformed["id"] = idProp
	}

	userNameProp := original["user_name"]
	e, err = isEmptyValue(reflect.ValueOf(userNameProp))
	if err != nil {
		return nil, err
	}
	if !e {
		transformed["userName"] = userNameProp
	}

	userPasswordProp := original["user_password"]
	e, err = isEmptyValue(reflect.ValueOf(userPasswordProp))
	if err != nil {
		return nil, err
	}
	if !e {
		transformed["userPassword"] = userPasswordProp
	}

	return transformed, nil
}

func expandMlsInstanceNetwork(v interface{}) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	availableZoneProp := original["available_zone"]
	e, err := isEmptyValue(reflect.ValueOf(availableZoneProp))
	if err != nil {
		return nil, err
	}
	if !e {
		transformed["availableZone"] = availableZoneProp
	}

	networkIDProp := original["network_id"]
	e, err = isEmptyValue(reflect.ValueOf(networkIDProp))
	if err != nil {
		return nil, err
	}
	if !e {
		transformed["subnetId"] = networkIDProp
	}

	publicIPProp, err := expandMlsInstanceNetworkPublicIP(original["public_ip"])
	if err != nil {
		return nil, err
	}
	e, err = isEmptyValue(reflect.ValueOf(publicIPProp))
	if err != nil {
		return nil, err
	}
	if !e {
		transformed["publicIP"] = publicIPProp
	}

	securityGroupIDProp := original["security_group_id"]
	e, err = isEmptyValue(reflect.ValueOf(securityGroupIDProp))
	if err != nil {
		return nil, err
	}
	if !e {
		transformed["securityGroupId"] = securityGroupIDProp
	}

	vpcIDProp := original["vpc_id"]
	e, err = isEmptyValue(reflect.ValueOf(vpcIDProp))
	if err != nil {
		return nil, err
	}
	if !e {
		transformed["vpcId"] = vpcIDProp
	}

	return transformed, nil
}

func expandMlsInstanceNetworkPublicIP(v interface{}) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	bindTypeProp := original["bind_type"]
	e, err := isEmptyValue(reflect.ValueOf(bindTypeProp))
	if err != nil {
		return nil, err
	}
	if !e {
		transformed["bindType"] = bindTypeProp
	}

	return transformed, nil
}
