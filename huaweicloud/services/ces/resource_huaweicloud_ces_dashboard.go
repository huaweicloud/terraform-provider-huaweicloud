package ces

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
			"creator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator name of the dashboard.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the dashboard.`,
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
	if rowWidgetNumNeedUpdate || isFavoriteOk {
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
	)

	return diag.FromErr(mErr.ErrorOrNil())
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
	opt.JSONBody = buildUpdateDashboardBodyParams(d)
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
	}
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
