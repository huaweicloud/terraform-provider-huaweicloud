package workspace

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

// @API Workspace GET /v1/{project_id}/app-warehouse/apps
func DataSourceAppWarehouseApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppWarehouseApplicationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the warehouse applications are located.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the application.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the application.`,
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The category of the application.`,
			},
			"verify_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The verification status of the application.`,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The record ID of the application.`,
						},
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the application.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the application.`,
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The category of the application.`,
						},
						"os_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The operating system type of the application.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the application.`,
						},
						"version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version name of the application.`,
						},
						"file_store_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage path of the application file.`,
						},
						"app_file_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The size of the application file.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the application.`,
						},
						"verify_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The verification status of the application.`,
						},
						"icon": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The base64 encoded application icon.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the application, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the application, in RFC3339 format.`,
						},
					},
				},
				Description: `All applications that match the filter parameters.`,
			},
		},
	}
}

func dataSourceAppWarehouseApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	applications, err := listAppWarehouseApplications(client, d)
	if err != nil {
		return diag.Errorf("error getting Workspace APP warehouse applications: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("applications", flattenAppWarehouseApplications(applications)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listAppWarehouseApplications(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/app-warehouse/apps"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	listPath += buildAppWarehouseApplicationsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		applications := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, applications...)
		if len(applications) < limit {
			break
		}

		offset += len(applications)
	}

	return result, nil
}

func buildAppWarehouseApplicationsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("app_id"); ok {
		res = fmt.Sprintf("%s&app_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&app_name=%v", res, v)
	}
	if v, ok := d.GetOk("category"); ok {
		res = fmt.Sprintf("%s&app_category=%v", res, v)
	}

	if v, ok := d.GetOk("verify_status"); ok {
		res = fmt.Sprintf("%s&verify_status=%v", res, v)
	}

	return res
}

func flattenAppWarehouseApplications(applications []interface{}) []map[string]interface{} {
	if len(applications) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(applications))
	for _, application := range applications {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", application, nil),
			"app_id":          utils.PathSearch("app_id", application, nil),
			"name":            utils.PathSearch("app_name", application, nil),
			"category":        utils.PathSearch("app_category", application, nil),
			"os_type":         utils.PathSearch("os_type", application, nil),
			"version":         utils.PathSearch("version_id", application, nil),
			"version_name":    utils.PathSearch("version_name", application, nil),
			"file_store_path": utils.PathSearch("appfile_store_path", application, nil),
			"app_file_size":   utils.PathSearch("app_file_size", application, nil),
			"description":     utils.PathSearch("app_description", application, nil),
			"verify_status":   utils.PathSearch("verify_status", application, nil),
			"icon":            utils.PathSearch("app_icon", application, nil),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
				application, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("modify_time",
				application, "").(string))/1000, false),
		})
	}

	return result
}
