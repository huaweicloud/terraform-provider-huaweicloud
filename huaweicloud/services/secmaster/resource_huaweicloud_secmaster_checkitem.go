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

var nonUpdatableParamsCheckitem = []string{"workspace_id"}

// @API SecMaster POST /v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems
// @API SecMaster POST /v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems/search
// @API SecMaster PUT /v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems/{checkitem_id}
// @API SecMaster DELETE /v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems
func ResourceCheckitem() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCheckitemCreate,
		ReadContext:   resourceCheckitemRead,
		UpdateContext: resourceCheckitemUpdate,
		DeleteContext: resourceCheckitemDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCheckitemImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsCheckitem),

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"level": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"informational", "low", "medium", "high", "fatal"}, false,
				),
			},
			"cloud_server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"method": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice(
					[]int{0, 1, 3, 4, 5},
				),
			},
			"audit_procedure": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aggregation_handle_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"impact": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"workflow_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateCheckitemBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"description":  d.Get("description"),
		"level":        d.Get("level"),
		"cloud_server": d.Get("cloud_server"),
		"method":       d.Get("method"),
		"source":       d.Get("source"),
	}

	if v, ok := d.GetOk("audit_procedure"); ok {
		bodyParams["audit_procedure"] = v
	}
	if v, ok := d.GetOk("aggregation_handle_status"); ok {
		bodyParams["aggregation_handle_status"] = v
	}
	if v, ok := d.GetOk("impact"); ok {
		bodyParams["impact"] = v
	}
	if v, ok := d.GetOk("workflow_id"); ok {
		bodyParams["workflow_id"] = v
	}

	return bodyParams
}

func resourceCheckitemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		workspaceId   = d.Get("workspace_id").(string)
		name          = d.Get("name").(string)
		createHttpUrl = "v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCheckitemBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster checkitem: %s", err)
	}

	// The name is unique, use this as the resource ID.
	d.SetId(name)

	return resourceCheckitemRead(ctx, d, meta)
}

func GetCheckitemInfo(client *golangsdk.ServiceClient, workspaceId, name string) (interface{}, error) {
	searchHttpUrl := "v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems/search"
	searchPath := client.Endpoint + searchHttpUrl
	searchPath = strings.ReplaceAll(searchPath, "{project_id}", client.ProjectID)
	searchPath = strings.ReplaceAll(searchPath, "{workspace_id}", workspaceId)

	searchOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"name": name,
			// type:1 indicates a custom check item.
			"type": 1,
		},
	}

	resp, err := client.Request("POST", searchPath, &searchOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	checkitems := utils.PathSearch("checkitems", respBody, make([]interface{}, 0)).([]interface{})
	if len(checkitems) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return checkitems[0], nil
}

func resourceCheckitemRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := GetCheckitemInfo(client, workspaceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster checkitem")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("level", utils.PathSearch("level", respBody, nil)),
		d.Set("cloud_server", utils.PathSearch("cloud_server", respBody, nil)),
		d.Set("method", utils.PathSearch("method", respBody, nil)),
		d.Set("audit_procedure", utils.PathSearch("audit_procedure", respBody, nil)),
		d.Set("impact", utils.PathSearch("impact", respBody, nil)),
		d.Set("source", utils.PathSearch("source", respBody, nil)),
		d.Set("workflow_id", utils.PathSearch("workflow_id", respBody, nil)),
		d.Set("aggregation_handle_status", utils.PathSearch("aggregation_handle_status", respBody, nil)),
		d.Set("uuid", utils.PathSearch("uuid", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateCheckitemBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"uuid":         d.Get("uuid"),
		"name":         d.Get("name"),
		"description":  d.Get("description"),
		"level":        d.Get("level"),
		"cloud_server": d.Get("cloud_server"),
		"method":       d.Get("method"),
		"source":       d.Get("source"),
	}

	if v, ok := d.GetOk("audit_procedure"); ok {
		bodyParams["audit_procedure"] = v
	}
	if v, ok := d.GetOk("aggregation_handle_status"); ok {
		bodyParams["aggregation_handle_status"] = v
	}
	if v, ok := d.GetOk("impact"); ok {
		bodyParams["impact"] = v
	}
	if v, ok := d.GetOk("workflow_id"); ok {
		bodyParams["workflow_id"] = v
	}

	return bodyParams
}

func resourceCheckitemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		uuid        = d.Get("uuid").(string)
		httpUrl     = "v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems/{checkitem_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{checkitem_id}", uuid)

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateCheckitemBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster checkitem: %s", err)
	}

	return resourceCheckitemRead(ctx, d, meta)
}

func resourceCheckitemDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		uuid        = d.Get("uuid").(string)
		httpUrl     = "v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)

	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"batch_ids": []string{uuid},
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "SecMaster.00092013"),
			"error deleting SecMaster checkitem",
		)
	}

	return nil
}

func resourceCheckitemImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<name>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
