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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/configurations
func DataSourceDrsJobConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsJobConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DRS job.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the parameter name for filtering.",
			},
			"parameter_config_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the parameter.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the parameter.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default value of the parameter.",
						},
						"value_range": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value range of the parameter.",
						},
						"is_need_restart": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the parameter needs a restart.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the parameter.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the parameter.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time of the parameter.",
						},
					},
				},
			},
		},
	}
}

func buildJobConfigurationsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("&name=%s", v.(string))
	}

	return queryParams
}

func dataSourceDrsJobConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/configurations"
		result  = make([]interface{}, 0)
		limit   = 1000
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{job_id}", d.Get("job_id").(string))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		queryParams := buildJobConfigurationsQueryParams(d)

		currentListPath := fmt.Sprintf("%s?limit=%d&offset=%d%s", listPath, limit, offset, queryParams)

		listResp, err := client.Request("GET", currentListPath, &reqOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		parameterConfigs := utils.PathSearch("parameter_config_list", listRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, parameterConfigs...)

		if len(parameterConfigs) == 0 {
			break
		}

		offset += len(parameterConfigs)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("parameter_config_list", flattenJobConfigurations(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenJobConfigurations(jobConfigsResp []interface{}) []interface{} {
	if len(jobConfigsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(jobConfigsResp))
	for _, v := range jobConfigsResp {
		rst = append(rst, map[string]interface{}{
			"name":            utils.PathSearch("name", v, nil),
			"value":           utils.PathSearch("value", v, nil),
			"default_value":   utils.PathSearch("default_value", v, nil),
			"value_range":     utils.PathSearch("value_range", v, nil),
			"is_need_restart": utils.PathSearch("is_need_restart", v, nil),
			"description":     utils.PathSearch("description", v, nil),
			"created_at":      utils.PathSearch("created_at", v, nil),
			"updated_at":      utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
