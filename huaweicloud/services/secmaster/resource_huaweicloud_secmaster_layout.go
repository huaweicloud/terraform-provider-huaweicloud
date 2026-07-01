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

var layoutNonUpdatableParams = []string{"workspace_id", "used_by", "region_id",
	"domain_id", "thumbnail", "layout_type", "binding_id", "binding_code"}

// 1. Parameters `is_built_in` and `is_template` do not take effect during creation and update, so they are temporarily
// not supported here.
// 2. The `create_time` and `creator_name` are parameters in the API documentation, but it does not actually
// take effect. Here, it is modified to an attribute field.

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/layouts
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/layouts/{layout_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/layouts/{layout_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/layouts
func ResourceLayout() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLayoutCreate,
		UpdateContext: resourceLayoutUpdate,
		ReadContext:   resourceLayoutRead,
		DeleteContext: resourceLayoutDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLayoutImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(layoutNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workspace ID.",
			},
			// Can be updated.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the layout name.",
			},
			"used_by": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the business type that uses the layout.",
			},
			// Can be updated.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the description.",
			},
			// Can be updated.
			"cloud_pack_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the cloud pack ID.",
			},
			// Can be updated.
			"cloud_pack_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the cloud pack name.",
			},
			// Can be updated.
			"cloud_pack_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the cloud pack version.",
			},
			// Can be updated.
			"layout_json": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  "Specifies the layout information in JSON format.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region ID.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the domain ID.",
			},
			"thumbnail": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the template thumbnail.",
			},
			"layout_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the layout type.",
			},
			"binding_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the data class ID or workflow ID.",
			},
			"binding_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the data class business code.",
			},
			// Can be updated.
			// The query API may not return a value, so do not add the `Computed` attribute.
			"fields_sum": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the total number of fields.",
			},
			// Can be updated.
			// The query API may not return a value, so do not add the `Computed` attribute.
			"wizards_sum": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the total number of pages.",
			},
			// Can be updated.
			"sections_sum": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the total number of system sections.",
			},
			// Can be updated.
			"tabs_sum": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the total number of custom tabs.",
			},
			// Can be updated.
			"boa_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the BOA version.",
			},
			"is_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specifies whether to directly delete the layout.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},
			"creator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator name.",
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator ID.",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The parent layout ID.",
			},
			"en_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The English description.",
			},
			"en_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The English name.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project ID.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time.",
			},
			"layout_cfg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The layout configuration used to bind icons on the frontend.",
			},
			"binding_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The data class name or workflow name.",
			},
			"modules_sum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of system modules.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SecMaster version.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildCreateLayoutBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":               d.Get("name"),
		"used_by":            d.Get("used_by"),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
		"cloud_pack_id":      utils.ValueIgnoreEmpty(d.Get("cloud_pack_id")),
		"cloud_pack_name":    utils.ValueIgnoreEmpty(d.Get("cloud_pack_name")),
		"cloud_pack_version": utils.ValueIgnoreEmpty(d.Get("cloud_pack_version")),
		"layout_json":        utils.StringToJson(d.Get("layout_json").(string)),
		"region_id":          utils.ValueIgnoreEmpty(d.Get("region_id")),
		"domain_id":          utils.ValueIgnoreEmpty(d.Get("domain_id")),
		"thumbnail":          utils.ValueIgnoreEmpty(d.Get("thumbnail")),
		"layout_type":        utils.ValueIgnoreEmpty(d.Get("layout_type")),
		"binding_id":         utils.ValueIgnoreEmpty(d.Get("binding_id")),
		"binding_code":       utils.ValueIgnoreEmpty(d.Get("binding_code")),
		"fields_sum":         utils.ValueIgnoreEmpty(d.Get("fields_sum")),
		"wizards_sum":        utils.ValueIgnoreEmpty(d.Get("wizards_sum")),
		"sections_sum":       utils.ValueIgnoreEmpty(d.Get("sections_sum")),
		"tabs_sum":           utils.ValueIgnoreEmpty(d.Get("tabs_sum")),
		"boa_version":        utils.ValueIgnoreEmpty(d.Get("boa_version")),
	}

	return bodyParams
}

func resourceLayoutCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts"
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
		JSONBody:         utils.RemoveNil(buildCreateLayoutBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster layout: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster layout: ID is not found in API response")
	}

	d.SetId(id)

	return resourceLayoutRead(ctx, d, meta)
}

func resourceLayoutRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/{layout_id}"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{layout_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			"SecMaster.20041303"), "error retrieving SecMaster layout")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	layout := utils.PathSearch("data", respBody, nil)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", layout, nil)),
		d.Set("used_by", utils.PathSearch("used_by", layout, nil)),
		d.Set("description", utils.PathSearch("description", layout, nil)),
		d.Set("cloud_pack_id", utils.PathSearch("cloud_pack_id", layout, nil)),
		d.Set("cloud_pack_name", utils.PathSearch("cloud_pack_name", layout, nil)),
		d.Set("cloud_pack_version", utils.PathSearch("cloud_pack_version", layout, nil)),
		d.Set("layout_json", utils.JsonToString(
			utils.PathSearch("layout_json", layout, nil))),
		d.Set("region_id", utils.PathSearch("region_id", layout, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", layout, nil)),
		d.Set("thumbnail", utils.PathSearch("thumbnail", layout, nil)),
		d.Set("layout_type", utils.PathSearch("layout_type", layout, nil)),
		d.Set("binding_id", utils.PathSearch("binding_id", layout, nil)),
		d.Set("binding_code", utils.PathSearch("binding_code", layout, nil)),
		d.Set("sections_sum", utils.PathSearch("sections_sum", layout, nil)),
		d.Set("tabs_sum", utils.PathSearch("tabs_sum", layout, nil)),
		d.Set("boa_version", utils.PathSearch("boa_version", layout, nil)),
		d.Set("create_time", utils.PathSearch("create_time", layout, nil)),
		d.Set("creator_name", utils.PathSearch("creator_name", layout, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", layout, nil)),
		d.Set("parent_id", utils.PathSearch("parent_id", layout, nil)),
		d.Set("en_description", utils.PathSearch("en_description", layout, nil)),
		d.Set("en_name", utils.PathSearch("en_name", layout, nil)),
		d.Set("project_id", utils.PathSearch("project_id", layout, nil)),
		d.Set("update_time", utils.PathSearch("update_time", layout, nil)),
		d.Set("layout_cfg", utils.PathSearch("layout_cfg", layout, nil)),
		d.Set("binding_name", utils.PathSearch("binding_name", layout, nil)),
		d.Set("modules_sum", utils.PathSearch("modules_sum", layout, nil)),
		d.Set("version", utils.PathSearch("version", layout, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateLayoutBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":               d.Get("name"),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
		"cloud_pack_id":      utils.ValueIgnoreEmpty(d.Get("cloud_pack_id")),
		"cloud_pack_name":    utils.ValueIgnoreEmpty(d.Get("cloud_pack_name")),
		"cloud_pack_version": utils.ValueIgnoreEmpty(d.Get("cloud_pack_version")),
		"layout_json":        utils.StringToJson(d.Get("layout_json").(string)),
		"fields_sum":         utils.ValueIgnoreEmpty(d.Get("fields_sum")),
		"wizards_sum":        utils.ValueIgnoreEmpty(d.Get("wizards_sum")),
		"sections_sum":       utils.ValueIgnoreEmpty(d.Get("sections_sum")),
		"tabs_sum":           utils.ValueIgnoreEmpty(d.Get("tabs_sum")),
		"boa_version":        utils.ValueIgnoreEmpty(d.Get("boa_version")),
	}

	return bodyParams
}

func resourceLayoutUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/{layout_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{layout_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildUpdateLayoutBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster layout: %s", err)
	}

	return resourceLayoutRead(ctx, d, meta)
}

func resourceLayoutDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts"
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
		JSONBody: map[string]interface{}{
			"ids":       []string{d.Id()},
			"is_delete": d.Get("is_delete").(bool),
		},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster layout: %s", err)
	}

	return nil
}

func resourceLayoutImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want '<workspace_id>/<id>', but got '%s'", importId)
	}

	d.SetId(importIdParts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", importIdParts[0])
}
