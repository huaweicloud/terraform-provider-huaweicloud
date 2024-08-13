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

// @API AOM GET /v2/{project_id}/aom/dashboards-folder
func DataSourceDashboardsFolders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDashboardsFoldersRead,

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
			"folders": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"folder_title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"folder_title_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dashboard_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"is_template": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
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
					},
				},
			},
		},
	}
}

func dataSourceDashboardsFoldersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	results, err := listDashboardsFolders(cfg, client, d)
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
		d.Set("folders", flattenFolders(results.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listDashboardsFolders(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v2/{project_id}/aom/dashboards-folder"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeadersForDataSource(cfg, d),
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboards folder: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening dashboards folder: %s", err)
	}

	return listRespBody, nil
}

func flattenFolders(folders []interface{}) []interface{} {
	if len(folders) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(folders))
	for _, folder := range folders {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("folder_id", folder, nil),
			"folder_title":          utils.PathSearch("folder_title", folder, nil),
			"folder_title_en":       utils.PathSearch("folder_title_en", folder, nil),
			"display":               utils.PathSearch("display", folder, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", folder, nil),
			"dashboard_ids":         utils.PathSearch("dashboard_ids", folder, nil),
			"is_template":           utils.PathSearch("is_template", folder, nil),
			"created_by":            utils.PathSearch("created_by", folder, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("created", folder, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("updated", folder, float64(0)).(float64))/1000, false),
		})
	}
	return result
}
