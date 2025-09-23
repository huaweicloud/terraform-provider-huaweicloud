package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM GET /v2/{project_id}/aom/dashboards
func DataSourceDashboards() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDashboardsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dashboard_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dashboards": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dashboard_title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dashboard_title_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"folder_title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"folder_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dashboard_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_favorite": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dashboard_tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeMap,
								Elem: &schema.Schema{Type: schema.TypeString},
							},
						},
						"display": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDashboardsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	results, err := listDashboards(cfg, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("dashboards", flattenDashboards(results.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListDashboardsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("dashboard_type"); ok {
		res = fmt.Sprintf("%s&dashboard_type=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func listDashboards(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v2/{project_id}/aom/dashboards"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListDashboardsQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeadersForDataSource(cfg, d),
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboards: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening dashboards: %s", err)
	}

	return utils.PathSearch("dashboards", listRespBody, make([]interface{}, 0)), nil
}

func flattenDashboards(dashboards []interface{}) []interface{} {
	if len(dashboards) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dashboards))
	for _, dashboard := range dashboards {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("dashboard_id", dashboard, nil),
			"dashboard_title":       utils.PathSearch("dashboard_title", dashboard, nil),
			"dashboard_title_en":    utils.PathSearch("dashboard_title_en", dashboard, nil),
			"folder_title":          utils.PathSearch("folder_name", dashboard, nil),
			"folder_id":             utils.PathSearch("folder_id", dashboard, nil),
			"version":               utils.PathSearch("version", dashboard, nil),
			"dashboard_type":        utils.PathSearch("dashboard_type", dashboard, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", dashboard, nil),
			"is_favorite":           utils.PathSearch("is_favorite", dashboard, nil),
			"dashboard_tags":        utils.PathSearch("dashboard_tags", dashboard, nil),
			"display":               utils.PathSearch("display", dashboard, nil),
			"created_by":            utils.PathSearch("created_by", dashboard, nil),
			"updated_by":            utils.PathSearch("updated_by", dashboard, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("created", dashboard, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("updated", dashboard, float64(0)).(float64))/1000, false),
		})
	}
	return result
}
