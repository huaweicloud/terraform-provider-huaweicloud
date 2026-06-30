package drs

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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/configuration-histories
func DataSourceDrsConfigurationHistories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsConfigurationHistoriesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameter_history_config_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"old_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"new_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_update_success": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_applied": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apply_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildConfigurationHistoriesQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := "?limit=1000"
	if v, ok := d.GetOk("begin_time"); ok {
		queryParams += fmt.Sprintf("&begin_time=%s", v.(string))
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams += fmt.Sprintf("&end_time=%s", v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams += fmt.Sprintf("&name=%s", v.(string))
	}
	if offset > 0 {
		queryParams += fmt.Sprintf("&offset=%d", offset)
	}

	return queryParams
}

func dataSourceDrsConfigurationHistoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/configuration-histories"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPathWithQuery := requestPath + buildConfigurationHistoriesQueryParams(d, offset)
		resp, err := client.Request("GET", requestPathWithQuery, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS configuration histories: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		parameterHistoryConfigList := utils.PathSearch(
			"parameter_history_config_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(parameterHistoryConfigList) == 0 {
			break
		}

		result = append(result, parameterHistoryConfigList...)
		offset += len(parameterHistoryConfigList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("parameter_history_config_list", flattenParameterHistoryConfigList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenParameterHistoryConfigList(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		result = append(result, map[string]interface{}{
			"name":              utils.PathSearch("name", item, nil),
			"old_value":         utils.PathSearch("old_value", item, nil),
			"new_value":         utils.PathSearch("new_value", item, nil),
			"is_update_success": utils.PathSearch("is_update_success", item, nil),
			"is_applied":        utils.PathSearch("is_applied", item, nil),
			"update_time":       utils.PathSearch("update_time", item, nil),
			"apply_time":        utils.PathSearch("apply_time", item, nil),
		})
	}

	return result
}
