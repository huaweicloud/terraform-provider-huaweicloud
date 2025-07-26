package ces

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParams = []string{"enterprise_project_id", "dashboard_id"}

// @API CES POST /v2/{project_id}/dashboards
// @API CES GET  /v2/{project_id}/dashboards
// @API CES PUT  /v2/{project_id}/dashboards/{dashboard_id}
// @API CES POST /v2/{project_id}/dashboards/batch-delete
func ResourceDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDashboardCreate,
		UpdateContext: resourceDashboardUpdate,
		ReadContext:   resourceDashboardRead,
		DeleteContext: resourceDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the dashboard name.`,
			},
			"row_widget_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the monitoring view display mode.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID of the dashboard.`,
			},
			"dashboard_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the copied dashboard ID.`,
			},
			"is_favorite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the dashboard is favorite.`,
			},
			"extend_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the metric aggregation method.`,
						},
						"period": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the metric aggregation period.`,
						},
						"display_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the display time.`,
						},
						"refresh_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the refresh time.`,
						},
						"from": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the start time.`,
						},
						"to": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the end time.`,
						},
						"screen_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the monitoring screen background color.`,
						},
						"enable_screen_auto_play": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether the monitoring screen switches automatically.`,
						},
						"time_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the automatic switching time interval of the monitoring screen.`,
						},
						"enable_legend": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether to enable the legend.`,
						},
						"full_screen_widget_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the number of large screen display views.`,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"creator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator name of the dashboard.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the dashboard.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the namespace.`,
			},
			"sub_product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the sub-product ID.`,
			},
			"dashboard_template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the monitoring disk template ID.`,
			},
			"widgets_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of views under the board.`,
			},
		},
	}
}

func resourceDashboardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/dashboards"
		product = "ces"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opt.JSONBody = utils.RemoveNil(buildCreateDashboardBodyParams(d, cfg))
	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return diag.Errorf("error creating CES dashboard: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("dashboard_id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CES dashboard: ID is not found in API response")
	}
	d.SetId(id)

	rowWidgetNumNeedUpdate := false
	dashboard, err := getDashboard(client, d)
	if err != nil {
		return diag.Errorf("error retrieving CES dashboard: %s", err)
	}
	specifiedRowWidgetNum := d.Get("row_widget_num")
	rowWidgetNum := utils.PathSearch("row_widget_num", dashboard, float64(0)).(float64)
	rowWidgetNumNeedUpdate = specifiedRowWidgetNum != int(rowWidgetNum)

	_, isFavoriteOk := d.GetOk("is_favorite")
	_, extendInfoOk := d.GetOk("extend_info")
	if rowWidgetNumNeedUpdate || isFavoriteOk || extendInfoOk {
		err = updateDashboard(client, d)
		if err != nil {
			return diag.Errorf("error updating CES dashboard: %s", err)
		}
	}
	return resourceDashboardRead(ctx, d, meta)
}

func buildCreateDashboardBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"dashboard_name": d.Get("name"),
		"row_widget_num": d.Get("row_widget_num"),
		"dashboard_id":   utils.ValueIgnoreEmpty(d.Get("dashboard_id")),
		"enterprise_id":  utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
}

func resourceDashboardRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	dashboard, err := getDashboard(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES dashboard")
	}

	createdAt := utils.PathSearch("create_time", dashboard, float64(0)).(float64) / 1000

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("dashboard_name", dashboard, nil)),
		d.Set("row_widget_num", utils.PathSearch("row_widget_num", dashboard, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_id", dashboard, nil)),
		d.Set("is_favorite", utils.PathSearch("is_favorite", dashboard, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createdAt), false)),
		d.Set("creator_name", utils.PathSearch("creator_name", dashboard, nil)),
		d.Set("namespace", utils.PathSearch("namespace", dashboard, nil)),
		d.Set("sub_product", utils.PathSearch("sub_product", dashboard, nil)),
		d.Set("dashboard_template_id", utils.PathSearch("dashboard_template_id", dashboard, nil)),
		d.Set("widgets_num", utils.PathSearch("widgets_num", dashboard, nil)),
		d.Set("extend_info", flattenCesDashboardDetail(
			utils.PathSearch("extend_info", dashboard, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCesDashboardDetail(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"filter":                  utils.PathSearch("filter", param, nil),
			"period":                  utils.PathSearch("period", param, nil),
			"display_time":            utils.PathSearch("display_time", param, nil),
			"refresh_time":            utils.PathSearch("refresh_time", param, nil),
			"from":                    utils.PathSearch("from", param, nil),
			"to":                      utils.PathSearch("to", param, nil),
			"screen_color":            utils.PathSearch("screen_color", param, nil),
			"enable_screen_auto_play": utils.PathSearch("enable_screen_auto_play", param, nil),
			"time_interval":           utils.PathSearch("time_interval", param, nil),
			"enable_legend":           utils.PathSearch("enable_legend", param, nil),
			"full_screen_widget_num":  utils.PathSearch("full_screen_widget_num", param, nil),
		},
	}

	return rst
}

func getDashboard(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/dashboards?dashboard_id={id}"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{id}", d.Id())

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", path, &opt)

	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	dashboard := utils.PathSearch("dashboards|[0]", respBody, nil)
	return dashboard, nil
}

func resourceDashboardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	err = updateDashboard(client, d)
	if err != nil {
		return diag.Errorf("error updating CES dashboard: %s", err)
	}

	return resourceDashboardRead(ctx, d, meta)
}

func updateDashboard(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/dashboards/{dashboard_id}"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{dashboard_id}", d.Id())

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	opt.JSONBody = utils.RemoveNil(buildUpdateDashboardBodyParams(d))
	_, err := client.Request("PUT", path, &opt)
	if err != nil {
		return fmt.Errorf("error updating CES dashboard: %s", err)
	}

	return nil
}

func buildUpdateDashboardBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"dashboard_name": d.Get("name"),
		"row_widget_num": d.Get("row_widget_num"),
		"is_favorite":    d.Get("is_favorite"),
		"extend_info":    buildUpdateDashboardExtendInfoBodyParams(d.Get("extend_info")),
	}
}

func buildUpdateDashboardExtendInfoBodyParams(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"filter":                  utils.ValueIgnoreEmpty(raw["filter"]),
			"period":                  utils.ValueIgnoreEmpty(raw["period"]),
			"display_time":            utils.ValueIgnoreEmpty(raw["display_time"]),
			"refresh_time":            utils.ValueIgnoreEmpty(raw["refresh_time"]),
			"from":                    utils.ValueIgnoreEmpty(raw["from"]),
			"to":                      utils.ValueIgnoreEmpty(raw["to"]),
			"screen_color":            utils.ValueIgnoreEmpty(raw["screen_color"]),
			"enable_screen_auto_play": utils.ValueIgnoreEmpty(raw["enable_screen_auto_play"]),
			"time_interval":           utils.ValueIgnoreEmpty(raw["time_interval"]),
			"enable_legend":           utils.ValueIgnoreEmpty(raw["enable_legend"]),
			"full_screen_widget_num":  utils.ValueIgnoreEmpty(raw["full_screen_widget_num"]),
		}

		return param
	}

	return nil
}

func resourceDashboardDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/dashboards/batch-delete"
		product = "ces"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opt.JSONBody = map[string]interface{}{
		"dashboard_ids": []string{d.Id()},
	}
	_, err = client.Request("POST", path, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CES dashboard")
	}

	return nil
}
