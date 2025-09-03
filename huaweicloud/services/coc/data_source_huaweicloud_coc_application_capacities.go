package coc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC POST /v1/capacity
func DataSourceCocApplicationCapacities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocApplicationCapacitiesRead,

		Schema: map[string]*schema.Schema{
			"provider_obj": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_service_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sum_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sum_cpu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sum_mem": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_service_name": {
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
		},
	}
}

func buildApplicationCapacitiesCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"provider_obj":   buildApplicationCapacitiesProviderObjCreateOpts(d.Get("provider_obj")),
		"group_id":       utils.ValueIgnoreEmpty(d.Get("group_id")),
		"component_id":   utils.ValueIgnoreEmpty(d.Get("component_id")),
		"application_id": utils.ValueIgnoreEmpty(d.Get("application_id")),
	}

	return bodyParams
}

func buildApplicationCapacitiesProviderObjCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"provider": raw["cloud_service_name"],
				"type":     raw["type"],
			}
		}
		return params
	}

	return nil
}

func dataSourceCocApplicationCapacitiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/capacity"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildApplicationCapacitiesCreateOpts(d)),
	}

	getResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error retrieving COC application capacities: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		nil,
		d.Set("data", flattenCocGetApplicationCapacities(
			utils.PathSearch("data", getRespBody, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCocGetApplicationCapacities(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"sum_size":           utils.PathSearch("sum_size", raw, nil),
				"sum_cpu":            utils.PathSearch("sum_cpu", raw, nil),
				"sum_mem":            utils.PathSearch("sum_mem", raw, nil),
				"cloud_service_name": utils.PathSearch("provider", raw, nil),
				"type":               utils.PathSearch("type", raw, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}
