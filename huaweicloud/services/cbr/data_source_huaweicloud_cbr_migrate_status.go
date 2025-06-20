package cbr

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

// @API CBR GET /v3/migrates
func DataSourceMigrateStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMigrateStatusRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to query the datasource.",
			},
			"all_regions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to query the migration results in other regions.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The overall migration status.",
			},
			"project_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of project migration status details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The migration status of the project.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project ID.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project name.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region ID.",
						},
						"progress": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The migration progress percentage.",
						},
						"fail_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The failure code when migration fails.",
						},
						"fail_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The failure reason when migration fails.",
						},
					},
				},
			},
		},
	}
}

func dataSourceMigrateStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/migrates"
		product = "cbr"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl

	queryParams := ""
	if d.Get("all_regions").(bool) {
		queryParams = "?all_regions=true"
	}

	requestPath += queryParams
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error querying CBR migrate status: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("project_status", flattenProjectStatus(utils.PathSearch("project_status", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the CBR migrate status: %s", mErr)
	}
	return nil
}

func flattenProjectStatus(projectStatusList []interface{}) []map[string]interface{} {
	if len(projectStatusList) < 1 {
		return nil
	}

	results := make([]map[string]interface{}, 0, len(projectStatusList))
	for _, projectStatus := range projectStatusList {
		projectStatusMap := map[string]interface{}{
			"status":       utils.PathSearch("status", projectStatus, nil),
			"project_id":   utils.PathSearch("project_id", projectStatus, nil),
			"project_name": utils.PathSearch("project_name", projectStatus, nil),
			"region_id":    utils.PathSearch("region_id", projectStatus, nil),
			"progress":     utils.PathSearch("progress", projectStatus, nil),
			"fail_code":    utils.PathSearch("fail_code", projectStatus, nil),
			"fail_reason":  utils.PathSearch("fail_reason", projectStatus, nil),
		}
		results = append(results, projectStatusMap)
	}

	return results
}
