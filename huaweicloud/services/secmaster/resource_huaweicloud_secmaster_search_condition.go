package secmaster

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

var searchConditionNonUpdatableParams = []string{"workspace_id", "dataspace_id", "pipe_id"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/search/conditions
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/search/conditions/{condition_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/siem/search/conditions/{condition_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/siem/search/conditions/{condition_id}
func ResourceSearchCondition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSearchConditionCreate,
		ReadContext:   resourceSearchConditionRead,
		UpdateContext: resourceSearchConditionUpdate,
		DeleteContext: resourceSearchConditionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSearchConditionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(searchConditionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The query detail API actually does not return this parameter.
			"dataspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipe_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
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

func buildSearchConditionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dataspace_id":   d.Get("dataspace_id"),
		"pipe_id":        d.Get("pipe_id"),
		"condition_name": d.Get("condition_name"),
		"query":          d.Get("query"),
	}

	return bodyParams
}

func resourceSearchConditionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/search/conditions"
		workspaceId   = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildSearchConditionBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating search condition: %s", err)
	}

	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	conditionId := utils.PathSearch("condition_id", respBody, "").(string)
	if conditionId == "" {
		return diag.Errorf("unable to find the condition ID from the API response")
	}

	d.SetId(conditionId)

	return resourceSearchConditionRead(ctx, d, meta)
}

func GetSearchConditionInfo(client *golangsdk.ServiceClient, workspaceId, conditionId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/siem/search/conditions/{condition_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{condition_id}", conditionId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceSearchConditionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	searchCondition, err := GetSearchConditionInfo(client, workspaceId, d.Id())
	if err != nil {
		// When the search condition does not exist, the response HTTP status code of the query API is `404`
		return common.CheckDeletedDiag(d, err, "error retrieving search condition")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("pipe_id", utils.PathSearch("pipe_id", searchCondition, nil)),
		d.Set("condition_name", utils.PathSearch("condition_name", searchCondition, nil)),
		d.Set("query", utils.PathSearch("query", searchCondition, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateSearchConditionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"condition_name": d.Get("condition_name"),
		"query":          d.Get("query"),
	}

	return bodyParams
}

func resourceSearchConditionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/search/conditions/{condition_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{condition_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateSearchConditionBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating search condition: %s", err)
	}

	return resourceSearchConditionRead(ctx, d, meta)
}

func resourceSearchConditionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/search/conditions/{condition_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{condition_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the search condition does not exist, the response HTTP status code of the deletion API is `404`.
		return common.CheckDeletedDiag(d, err, "error deleting search condition")
	}

	return nil
}

func resourceSearchConditionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
