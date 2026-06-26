package secmaster

import (
	"context"
	"errors"
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

var securityReportNonUpdatableParams = []string{
	"workspace_id", "report_name", "report_period", "report_range", "report_range.*.start", "report_range.*.end",
	"language", "layout_id", "binding_wizard",
}

// Due to various issues with the API, this resource currently only supports some mandatory parameters, and other
// optional parameters are temporarily not supported.

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/sa/reports
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/sa/reports
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/sa/reports/{report_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/sa/reports/{report_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/sa/reports/{report_id}
func ResourceSecurityReport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityReportCreate,
		ReadContext:   resourceSecurityReportRead,
		UpdateContext: resourceSecurityReportUpdate,
		DeleteContext: resourceSecurityReportDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityReportImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(securityReportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located.",
			},
			// Query API no return.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workspace ID.",
			},
			"report_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the report name.",
			},
			"report_period": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the report period.",
			},
			// The value returned by the API query is inconsistent with the value created, so it is not set in the read.
			"report_range": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the start time of the data range.",
						},
						"end": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the end time of the data range.",
						},
					},
				},
				Description: "Specifies the data range.",
			},
			// Query API no return.
			"language": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the language.",
			},
			"layout_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the layout ID.",
			},
			"binding_wizard": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the report page content.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the report status.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCreateSecurityReportBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"report_name":    d.Get("report_name"),
		"report_period":  d.Get("report_period"),
		"report_range":   buildSecurityReportRangeBodyParams(d.Get("report_range").([]interface{})),
		"language":       d.Get("language"),
		"layout_id":      d.Get("layout_id"),
		"binding_wizard": d.Get("binding_wizard"),
	}

	return bodyParams
}

func buildSecurityReportRangeBodyParams(rawParams []interface{}) map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	raw, ok := rawParams[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"start": raw["start"],
		"end":   raw["end"],
	}
}

func resourceSecurityReportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/sa/reports"
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildCreateSecurityReportBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster security report: %s", err)
	}

	// After creation, the `status` of the resource defaults to **enable**.
	reportId, err := querySecurityReportIdByName(client, workspaceId, d.Get("report_name").(string),
		d.Get("report_period").(string), "enable")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(reportId)

	if v, ok := d.GetOk("status"); ok {
		if err := updateSecurityReportStatus(client, workspaceId, reportId, v.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSecurityReportRead(ctx, d, meta)
}

func querySecurityReportIdByName(client *golangsdk.ServiceClient, workspaceId, reportName, reportPeriod, status string) (string, error) {
	listPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/sa/reports"
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", workspaceId)
	listPath = fmt.Sprintf("%s?report_period=%s&status=%s", listPath, reportPeriod, status)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return "", fmt.Errorf("error querying SecMaster security reports: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return "", err
	}

	listArray, ok := listRespBody.([]interface{})
	if !ok {
		return "", errors.New("unable to convert the security report list response")
	}

	for _, v := range listArray {
		if utils.PathSearch("report_name", v, "").(string) == reportName {
			return utils.PathSearch("id", v, "").(string), nil
		}
	}

	return "", errors.New("unable to find the security report ID from the API response")
}

func resourceSecurityReportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/sa/reports/{report_id}"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{report_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster security report: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// Even if the resource ID does not exist, the query API will return `200`, so it is necessary to determine whether
	// the `id` has a value.
	if utils.PathSearch("id", respBody, "").(string) == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("report_name", utils.PathSearch("report_name", respBody, nil)),
		d.Set("report_period", utils.PathSearch("report_period", respBody, nil)),
		d.Set("layout_id", utils.PathSearch("layout_id", respBody, nil)),
		d.Set("binding_wizard", utils.PathSearch("binding_wizard", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateSecurityReportStatus(client *golangsdk.ServiceClient, workspaceId, reportId, status string) error {
	requestPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/sa/reports/{report_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{report_id}", reportId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody: map[string]interface{}{
			"status": status,
		},
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating SecMaster security report status: %s", err)
	}

	return nil
}

func resourceSecurityReportUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	if d.HasChange("status") {
		if err := updateSecurityReportStatus(client, workspaceId, d.Id(), d.Get("status").(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSecurityReportRead(ctx, d, meta)
}

func resourceSecurityReportDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/sa/reports/{report_id}"
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{report_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster security report: %s", err)
	}

	return nil
}

func resourceSecurityReportImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want '<workspace_id>/<id>', but got '%s'", importId)
	}

	d.SetId(importIdParts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", importIdParts[0])
}
