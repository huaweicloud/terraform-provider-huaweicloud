package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/devices
func DataSourceDscDevices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscDevicesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"devices": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The device information list.",
				Elem:        dscDeviceSchema(),
			},
		},
	}
}

func dscDeviceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device name.",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device IP address.",
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The device type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device status.",
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deployment mode.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device version.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device description.",
			},
			"manage_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The management URL.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPC ID.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subnet ID.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The creation time.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update time.",
			},
			"related_datasource_policy_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The related datasource policy information list.",
				Elem:        dscRelatedDatasourceInfoSchema(),
			},
		},
	}
}

func dscRelatedDatasourceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datasource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The datasource asset ID.",
			},
			"datasource_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The datasource address.",
			},
			"datasource_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The datasource port.",
			},
			"proxy_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The proxy port.",
			},
			"ddm_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The dynamic masking policy information list.",
				Elem:        dscDdmPolicySchema(),
			},
			"gde_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The encryption policy information list.",
				Elem:        dscGdePolicySchema(),
			},
			"sdm_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The static masking policy information list.",
				Elem:        dscSdmPolicySchema(),
			},
		},
	}
}

func dscDdmPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The namespace name.",
			},
			"table": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The table name.",
			},
			"columns": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column information list.",
				Elem:        dscColumnSchema(),
			},
		},
	}
}

func dscGdePolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The action. 1 means encrypt, 2 means decrypt.",
			},
			"alg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The encryption algorithm.",
			},
			"table": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The table name.",
			},
			"columns": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column information list.",
				Elem:        dscColumnSchema(),
			},
		},
	}
}

func dscSdmPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The table name.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The namespace name.",
			},
			"do_mask": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to mask data.",
			},
			"do_move": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to move data.",
			},
			"columns": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column information list.",
				Elem:        dscColumnSchema(),
			},
		},
	}
}

func dscColumnSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The column name.",
			},
			"mask": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The masking algorithm name or ID.",
			},
		},
	}
}

func buildDscDevicesQueryParams(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func dataSourceDscDevicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/devices"
		offset  = 0
		limit   = 1000
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildDscDevicesQueryParams(limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC devices: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		devicesResp := utils.PathSearch("devices", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, devicesResp...)

		if len(devicesResp) < limit {
			break
		}

		offset += len(devicesResp)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("devices", flattenDscDevices(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscDevices(devices []interface{}) []interface{} {
	if len(devices) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(devices))
	for _, v := range devices {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"ip":          utils.PathSearch("ip", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"mode":        utils.PathSearch("mode", v, nil),
			"version":     utils.PathSearch("version", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"manage_url":  utils.PathSearch("manage_url", v, nil),
			"vpc_id":      utils.PathSearch("vpc_id", v, nil),
			"subnet_id":   utils.PathSearch("subnet_id", v, nil),
			"create_time": utils.PathSearch("create_time", v, nil),
			"update_time": utils.PathSearch("update_time", v, nil),
			"related_datasource_policy_list": flattenDscRelatedDatasourcePolicyList(
				utils.PathSearch("related_datasource_policy_list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscRelatedDatasourcePolicyList(policies []interface{}) []interface{} {
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
			"ddm_policies": flattenDscDdmPolicies(
				utils.PathSearch("ddm_policies", v, make([]interface{}, 0)).([]interface{})),
			"gde_policies": flattenDscGdePolicies(
				utils.PathSearch("gde_policies", v, make([]interface{}, 0)).([]interface{})),
			"sdm_policies": flattenDscSdmPolicies(
				utils.PathSearch("sdm_policies", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscDdmPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"namespace": utils.PathSearch("namespace", v, nil),
			"table":     utils.PathSearch("table", v, nil),
			"columns": flattenDscColumns(
				utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscGdePolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"action":  utils.PathSearch("action", v, nil),
			"alg":     utils.PathSearch("alg", v, nil),
			"table":   utils.PathSearch("table", v, nil),
			"columns": flattenDscColumns(utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscSdmPolicies(policies []interface{}) []interface{} {
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
			"columns":   flattenDscColumns(utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscColumns(columns []interface{}) []interface{} {
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
