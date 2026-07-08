package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC POST /v1/{project_id}/devices
// @API DSC GET /v1/{project_id}/devices
// @API DSC PUT /v1/{project_id}/devices/{device_id}
// @API DSC DELETE /v1/{project_id}/devices/{device_id}
func ResourceDscDevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDscDeviceCreate,
		ReadContext:   resourceDscDeviceRead,
		UpdateContext: resourceDscDeviceUpdate,
		DeleteContext: resourceDscDeviceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"manage_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"related_datasource_policy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscRelatedResourceInfoSchema(),
			},
		},
	}
}

func dscRelatedResourceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datasource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datasource_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datasource_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"proxy_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ddm_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscResourceDdmPolicySchema(),
			},
			"gde_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscResourceGdePolicySchema(),
			},
			"sdm_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscResourceSdmPolicySchema(),
			},
		},
	}
}

func dscResourceDdmPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscResourceColumnSchema(),
			},
		},
	}
}

func dscResourceGdePolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"alg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscResourceColumnSchema(),
			},
		},
	}
}

func dscResourceSdmPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"do_mask": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"do_move": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscResourceColumnSchema(),
			},
		},
	}
}

func dscResourceColumnSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mask": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCreateDscDeviceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"mode":        d.Get("mode"),
		"vpc_id":      d.Get("vpc_id"),
		"subnet_id":   d.Get("subnet_id"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"manage_url":  utils.ValueIgnoreEmpty(d.Get("manage_url")),
	}
}

// GetDeviceIdByName is used to obtain the device ID by device name.
// Pagination parameters for the list API are invalid.
func GetDeviceIdByName(client *golangsdk.ServiceClient, deviceName string) (string, error) {
	requestPath := client.Endpoint + "v1/{project_id}/devices"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = fmt.Sprintf("%s?limit=1000&offset=0", requestPath)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	return utils.PathSearch(fmt.Sprintf("devices[?name=='%s']|[0].id", deviceName), respBody, "").(string), nil
}

func resourceDscDeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/devices"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDscDeviceBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating DSC device: %s", err)
	}

	deviceId, err := GetDeviceIdByName(client, d.Get("name").(string))
	if err != nil || deviceId == "" {
		return diag.Errorf("error getting DSC device ID after creation: %s", err)
	}

	d.SetId(deviceId)

	return resourceDscDeviceRead(ctx, d, meta)
}

func GetDeviceById(client *golangsdk.ServiceClient, deviceId string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/devices"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = fmt.Sprintf("%s?limit=1000&offset=0", requestPath)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	device := utils.PathSearch(fmt.Sprintf("devices[?id=='%s']|[0]", deviceId), respBody, nil)
	if device == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return device, nil
}

func resourceDscDeviceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	device, err := GetDeviceById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DSC device")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", device, nil)),
		d.Set("type", utils.PathSearch("type", device, nil)),
		d.Set("mode", utils.PathSearch("mode", device, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", device, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", device, nil)),
		d.Set("description", utils.PathSearch("description", device, nil)),
		d.Set("manage_url", utils.PathSearch("manage_url", device, nil)),
		d.Set("ip", utils.PathSearch("ip", device, nil)),
		d.Set("status", utils.PathSearch("status", device, nil)),
		d.Set("version", utils.PathSearch("version", device, nil)),
		d.Set("create_time", utils.PathSearch("create_time", device, nil)),
		d.Set("update_time", utils.PathSearch("update_time", device, nil)),
		d.Set("related_datasource_policy_list", flattenDscResourceRelatedDatasourcePolicyList(
			utils.PathSearch("related_datasource_policy_list", device, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscResourceRelatedDatasourcePolicyList(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"datasource_id":      utils.PathSearch("datasource_id", v, nil),
			"datasource_address": utils.PathSearch("datasource_address", v, nil),
			"datasource_port":    utils.PathSearch("datasource_port", v, nil),
			"proxy_port":         utils.PathSearch("proxy_port", v, nil),
			"ddm_policies": flattenDscResourceDdmPolicies(
				utils.PathSearch("ddm_policies", v, make([]interface{}, 0)).([]interface{})),
			"gde_policies": flattenDscResourceGdePolicies(
				utils.PathSearch("gde_policies", v, make([]interface{}, 0)).([]interface{})),
			"sdm_policies": flattenDscResourceSdmPolicies(
				utils.PathSearch("sdm_policies", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscResourceDdmPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"namespace": utils.PathSearch("namespace", v, nil),
			"table":     utils.PathSearch("table", v, nil),
			"columns": flattenDscResourceColumns(
				utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscResourceGdePolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"action":  utils.PathSearch("action", v, nil),
			"alg":     utils.PathSearch("alg", v, nil),
			"table":   utils.PathSearch("table", v, nil),
			"columns": flattenDscResourceColumns(utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscResourceSdmPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"table":     utils.PathSearch("table", v, nil),
			"namespace": utils.PathSearch("namespace", v, nil),
			"do_mask":   utils.PathSearch("do_mask", v, nil),
			"do_move":   utils.PathSearch("do_move", v, nil),
			"columns":   flattenDscResourceColumns(utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscResourceColumns(columns []interface{}) []interface{} {
	if len(columns) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(columns))
	for _, v := range columns {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"mask": utils.PathSearch("mask", v, nil),
		})
	}

	return rst
}

func buildUpdateDscDeviceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"mode":        d.Get("mode"),
		"vpc_id":      d.Get("vpc_id"),
		"subnet_id":   d.Get("subnet_id"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"manage_url":  utils.ValueIgnoreEmpty(d.Get("manage_url")),
	}
}

func resourceDscDeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/devices/{device_id}"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateDscDeviceBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating DSC device: %s", err)
	}

	return resourceDscDeviceRead(ctx, d, meta)
}

func resourceDscDeviceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/devices/{device_id}"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DSC device: %s", err)
	}

	return nil
}
