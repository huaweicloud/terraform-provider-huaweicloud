package coc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCocScriptOrderBatchDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocScriptOrderBatchDetailsRead,

		Schema: map[string]*schema.Schema{
			"batch_index": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"execute_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execute_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target_instance": executeInstancesTargetInstances(),
						"gmt_created": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gmt_finished": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"execute_costs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func executeInstancesTargetInstances() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"resource_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"provider": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"region_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"custom_attributes": executeInstancesTargetInstancesCustomAttributes(),
				"agent_sn": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"agent_status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"properties": executeInstancesTargetInstancesProperties(),
			},
		},
	}
}

func executeInstancesTargetInstancesCustomAttributes() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func executeInstancesTargetInstancesProperties() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"host_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"fixed_ip": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"floating_ip": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"region_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"zone_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"application": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"group": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"project_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// @API COC GET /v1/job/script/orders/{execute_uuid}/batches/{batch_index}
func dataSourceCocScriptOrderBatchDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl     = "v1/job/script/orders/{execute_uuid}/batches/{batch_index}"
		product     = "coc"
		executeUUID = d.Get("execute_uuid").(string)
		batchIndex  = d.Get("batch_index").(int)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{execute_uuid}", executeUUID)
	basePath = strings.ReplaceAll(basePath, "{batch_index}", strconv.Itoa(batchIndex))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	marker := 0.0
	res := make([]map[string]interface{}, 0)
	for {
		getPath := basePath + buildGetScriptOrderBatchDetailsParams(d, marker)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving COC script order batch details: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		batchDetails, nextMarker := flattenCocGetScriptOrderBatchDetails(getRespBody)
		if len(batchDetails) < 1 {
			break
		}
		res = append(res, batchDetails...)
		marker = nextMarker
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("execute_instances", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetScriptOrderBatchDetailsParams(d *schema.ResourceData, marker float64) string {
	res := "?limit=100"
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if marker != 0 {
		res = fmt.Sprintf("%s&marker=%v", res, int(marker))
	}

	return res
}

func flattenCocGetScriptOrderBatchDetails(resp interface{}) ([]map[string]interface{}, float64) {
	batchDetailsJson := utils.PathSearch("data.execute_instances", resp, make([]interface{}, 0))
	batchDetailsArray := batchDetailsJson.([]interface{})
	if len(batchDetailsArray) == 0 {
		return nil, 0
	}

	result := make([]map[string]interface{}, 0, len(batchDetailsArray))
	var marker float64
	for _, batchDetail := range batchDetailsArray {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", batchDetail, nil),
			"gmt_created":     utils.PathSearch("gmt_created", batchDetail, nil),
			"gmt_finished":    utils.PathSearch("gmt_finished", batchDetail, nil),
			"execute_costs":   utils.PathSearch("execute_costs", batchDetail, nil),
			"status":          utils.PathSearch("status", batchDetail, nil),
			"message":         utils.PathSearch("message", batchDetail, nil),
			"target_instance": flattenCocGetScriptOrderBatchDetailTargetInstances(utils.PathSearch("target_instance", batchDetail, nil)),
		})
		marker = utils.PathSearch("id", batchDetail, float64(0)).(float64)
	}
	return result, marker
}

func flattenCocGetScriptOrderBatchDetailTargetInstances(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"resource_id":  utils.PathSearch("resource_id", param, nil),
			"provider":     utils.PathSearch("provider", param, nil),
			"region_id":    utils.PathSearch("region_id", param, nil),
			"type":         utils.PathSearch("type", param, nil),
			"agent_sn":     utils.PathSearch("agent_sn", param, nil),
			"agent_status": utils.PathSearch("agent_status", param, nil),
			"custom_attributes": flattenScriptOrderBatchDetailCustomAttributes(
				utils.PathSearch("custom_attributes", param, make([]interface{}, 0)).([]interface{})),
			"properties": flattenScriptOrderBatchDetailProperties(utils.PathSearch("properties", param, nil)),
		},
	}

	return rst
}

func flattenScriptOrderBatchDetailCustomAttributes(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"key":   utils.PathSearch("key", params, nil),
			"value": utils.PathSearch("value", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenScriptOrderBatchDetailProperties(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"host_name":   utils.PathSearch("host_name", param, nil),
			"fixed_ip":    utils.PathSearch("fixed_ip", param, nil),
			"floating_ip": utils.PathSearch("floating_ip", param, nil),
			"region_id":   utils.PathSearch("region_id", param, nil),
			"zone_id":     utils.PathSearch("zone_id", param, nil),
			"application": utils.PathSearch("application", param, nil),
			"group":       utils.PathSearch("group", param, nil),
			"project_id":  utils.PathSearch("project_id", param, nil),
		},
	}

	return rst
}
