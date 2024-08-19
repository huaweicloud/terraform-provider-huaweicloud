package dataarts

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

// @API DataArtsStudio GET /v1/{project_id}/service/apps
func DataSourceDataServiceApps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceAppsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the applications are located.`,
			},

			// Parameters in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the applications belong.`,
			},
			"dlm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of DLM engine.`,
			},

			// Query argument
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the applications to be fuzzy queried.`,
			},

			// Attribute
			"apps": {
				Type:        schema.TypeList,
				Elem:        dataserviceAppSchema(),
				Computed:    true,
				Description: `All applications that match the filter parameters.`,
			},
		},
	}
}

func dataserviceAppSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the application, in UUID format.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the application.`,
			},
			"app_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key of the application.`,
			},
			"app_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The secret of the application.`,
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
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application creator.`,
			},
			"update_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application updater.`,
			},
			"app_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the application.`,
			},
		},
	}
	return &sc
}

func buildDataServiceAppsQueryParams(d *schema.ResourceData) string {
	res := ""
	if appName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, appName)
	}
	return res
}

func queryDataServiceApps(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/service/apps?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildDataServiceAppsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm-type":     d.Get("dlm_type").(string),
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		apps := utils.PathSearch("apps", respBody, make([]interface{}, 0)).([]interface{})
		if len(apps) < 1 {
			break
		}
		result = append(result, apps...)
		offset += len(apps)
	}

	return result, nil
}

func flattenDataServiceApps(apps []interface{}) []interface{} {
	result := make([]interface{}, 0, len(apps))

	for _, app := range apps {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", app, nil),
			"name":        utils.PathSearch("name", app, nil),
			"description": utils.PathSearch("description", app, nil),
			"app_key":     utils.PathSearch("app_key", app, nil),
			"app_secret":  utils.PathSearch("app_secret", app, nil),
			"created_at":  utils.FormatTimeStampRFC3339(int64(utils.PathSearch("register_time", app, float64(0)).(float64))/1000, false),
			"updated_at":  utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", app, float64(0)).(float64))/1000, false),
			"create_user": utils.PathSearch("create_user", app, nil),
			"update_user": utils.PathSearch("update_user", app, nil),
			"app_type":    utils.PathSearch("app_type", app, nil),
		})
	}

	return result
}

func dataSourceDataServiceAppsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	appList, err := queryDataServiceApps(client, d)
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
		d.Set("apps", flattenDataServiceApps(appList)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
