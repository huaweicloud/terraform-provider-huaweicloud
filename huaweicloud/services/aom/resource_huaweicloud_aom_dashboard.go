package aom

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v2/{project_id}/aom/dashboards
// @API AOM DELETE /v2/{project_id}/aom/dashboards/{dashboard_id}
// @API AOM GET /v2/{project_id}/aom/dashboards/{dashboard_id}
// @API AOM GET /v2/{project_id}/aom/dashboards
func ResourceDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDashboardCreate,
		ReadContext:   resourceDashboardRead,
		UpdateContext: resourceDashboardUpdate,
		DeleteContext: resourceDashboardDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dashboard_title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"folder_title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dashboard_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"is_favorite": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"dashboard_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			"charts": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDashboardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/aom/dashboards"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(cfg, d),
		JSONBody:         utils.RemoveNil(buildCreateDashboardBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM dashboard: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening creating dashboard response: %s", err)
	}

	id := utils.PathSearch("dashboard_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating AOM dashboard: can not found dashboard_id in return")
	}

	d.SetId(id)

	return resourceDashboardRead(ctx, d, meta)
}

func buildCreateDashboardBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dashboard_title": d.Get("dashboard_title"),
		"folder_name":     d.Get("folder_title"),
		"dashboard_type":  d.Get("dashboard_type"),
		"is_favorite":     utils.ValueIgnoreEmpty(d.Get("is_favorite")),
		"dashboard_tags":  utils.ValueIgnoreEmpty(d.Get("dashboard_tags")),
	}

	if v, ok := d.GetOk("charts"); ok {
		bodyParams["charts"] = parseCharts(v.(string))
	}

	return bodyParams
}

func parseCharts(v string) interface{} {
	if v == "" {
		return nil
	}
	// charts may be list or map
	var data interface{}
	err := json.Unmarshal([]byte(v), &data)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse JSON: %s", err)
		return v
	}

	return data
}

func encodeCharts(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	rst, err := json.Marshal(v)
	if err != nil {
		log.Printf("[DEBUG] Unable to encode charts: %s", err)
		return nil
	}

	return string(rst)
}

func resourceDashboardRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	// use list to check delete for get will retrieve a non exist dashboard
	err = filterDashboard(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving dashboard")
	}

	// use get to retrieve info for some details are not in list retrun
	dashboard, err := getDashboard(cfg, client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving dashboard")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("dashboard_title", utils.PathSearch("dashboard_title", dashboard, nil)),
		d.Set("folder_title", utils.PathSearch("folder_name", dashboard, nil)),
		d.Set("dashboard_type", utils.PathSearch("dashboard_type", dashboard, nil)),
		d.Set("is_favorite", utils.PathSearch("is_favorite", dashboard, nil)),
		d.Set("dashboard_tags", utils.PathSearch("dashboard_tags", dashboard, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", dashboard, nil)),
		d.Set("charts", encodeCharts(utils.PathSearch("charts", dashboard, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDashboard(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getHttpUrl := "v2/{project_id}/aom/dashboards/{dashboard_id}?version=v1"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{dashboard_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(cfg, d),
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboard: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening dashboard: %s", err)
	}

	return getRespBody, nil
}

func filterDashboard(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	listHttpUrl := "v2/{project_id}/aom/dashboards"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":          "application/json",
			"Enterprise-Project-Id": "all_granted_eps",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return fmt.Errorf("error retrieving dashboards: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return fmt.Errorf("error flattening dashboards: %s", err)
	}

	jsonPath := fmt.Sprintf("dashboards[?dashboard_id=='%s']|[0]", d.Id())
	dashboard := utils.PathSearch(jsonPath, listRespBody, nil)
	if dashboard == nil {
		return golangsdk.ErrDefault404{}
	}

	return nil
}

func resourceDashboardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	changes := []string{
		"dashboard_title",
		"folder_name",
		"dashboard_type",
		"is_favorite",
		"dashboard_tags",
		"charts",
	}

	if d.HasChanges(changes...) {
		updateHttpUrl := "v2/{project_id}/aom/dashboards"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      buildHeaders(cfg, d),
			JSONBody:         buildUpdateDashboardBodyParams(d),
		}

		_, err = client.Request("POST", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating dashboard: %s", err)
		}
	}

	return resourceDashboardRead(ctx, d, meta)
}

func buildUpdateDashboardBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"is_create_action": false,
		"version":          "v1",
		"dashboard_id":     d.Id(),
		"dashboard_title":  d.Get("dashboard_title"),
		"folder_name":      d.Get("folder_title"),
		"dashboard_type":   d.Get("dashboard_type"),
		"is_favorite":      utils.ValueIgnoreEmpty(d.Get("is_favorite")),
		"dashboard_tags":   utils.ValueIgnoreEmpty(d.Get("dashboard_tags")),
		"charts":           parseCharts(d.Get("charts").(string)),
	}

	return bodyParams
}

func resourceDashboardDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	// Delete will return 200 even deleting a non exist dashboard
	err = filterDashboard(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting dashboard")
	}

	deleteHttpUrl := "v2/{project_id}/aom/dashboards/{dashboard_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{dashboard_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(cfg, d),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting dashboard")
	}

	return nil
}
