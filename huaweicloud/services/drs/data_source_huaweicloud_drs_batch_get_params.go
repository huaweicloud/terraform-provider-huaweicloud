package drs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v3/{project_id}/jobs/batch-get-params
func DataSourceDrsBatchGetParams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBatchGetParamsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"refresh": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"params_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"compare_result": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"data_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"need_restart": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildJobParamsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"jobs":    utils.ExpandToStringList(d.Get("job_ids").([]interface{})),
		"refresh": d.Get("refresh"),
	}
	return utils.RemoveNil(bodyParams)
}

func dataSourceBatchGetParamsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v3/{project_id}/jobs/batch-get-params"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildJobParamsBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("params_list", flattenParamsList(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenParamsList(respBody interface{}) []interface{} {
	paramsListRaw := utils.PathSearch("params_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(paramsListRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(paramsListRaw))
	for _, item := range paramsListRaw {
		result = append(result, map[string]interface{}{
			"job_id": utils.PathSearch("job_id", item, nil),
			"params": flattenParams(utils.PathSearch("params", item, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenParams(paramsRaw []interface{}) []interface{} {
	result := make([]interface{}, 0, len(paramsRaw))
	for _, param := range paramsRaw {
		result = append(result, map[string]interface{}{
			"group":          utils.PathSearch("group", param, nil),
			"key":            utils.PathSearch("key", param, nil),
			"source_value":   utils.PathSearch("source_value", param, nil),
			"target_value":   utils.PathSearch("target_value", param, nil),
			"compare_result": utils.PathSearch("compare_result", param, nil),
			"data_type":      utils.PathSearch("data_type", param, nil),
			"value_range":    utils.PathSearch("value_range", param, nil),
			"need_restart":   utils.PathSearch("need_restart", param, nil),
		})
	}
	return result
}
