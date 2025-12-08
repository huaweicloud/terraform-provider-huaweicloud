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

// The details API has issues; we are using the list API here.

// @API SecMaster POST /v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields
// @API SecMaster GET /v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields
// @API SecMaster PUT /v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields/{field_id}
// @API SecMaster DELETE /v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields
func ResourceLayoutField() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLayoutFieldCreate,
		ReadContext:   resourceLayoutFieldRead,
		UpdateContext: resourceLayoutFieldUpdate,
		DeleteContext: resourceLayoutFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLayoutFieldImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{"workspace_id"}),

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
			"field_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// This field `layout_id` is not exist in API response.
			"layout_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"wizard_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aopworkflow_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aopworkflow_version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"playbook_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"playbook_version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default_value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// This field `display_type` is not exist in API response.
			"display_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"field_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"extra_json": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"field_tooltip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"json_schema": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"searchable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"visible": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"maintainable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"editable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"creatable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"boa_version": {
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
			// Computed attributes
			"cloud_pack_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_pack_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataclass_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_pack_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"en_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"en_default_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"en_field_tooltip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_built_in": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateLayoutFieldBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                   d.Get("name"),
		"field_key":              d.Get("field_key"),
		"description":            utils.ValueIgnoreEmpty(d.Get("description")),
		"layout_id":              utils.ValueIgnoreEmpty(d.Get("layout_id")),
		"wizard_id":              utils.ValueIgnoreEmpty(d.Get("wizard_id")),
		"aopworkflow_id":         utils.ValueIgnoreEmpty(d.Get("aopworkflow_id")),
		"aopworkflow_version_id": utils.ValueIgnoreEmpty(d.Get("aopworkflow_version_id")),
		"playbook_id":            utils.ValueIgnoreEmpty(d.Get("playbook_id")),
		"playbook_version_id":    utils.ValueIgnoreEmpty(d.Get("playbook_version_id")),
		"default_value":          utils.ValueIgnoreEmpty(d.Get("default_value")),
		"display_type":           utils.ValueIgnoreEmpty(d.Get("display_type")),
		"field_type":             utils.ValueIgnoreEmpty(d.Get("field_type")),
		"extra_json":             utils.ValueIgnoreEmpty(d.Get("extra_json")),
		"field_tooltip":          utils.ValueIgnoreEmpty(d.Get("field_tooltip")),
		"json_schema":            utils.ValueIgnoreEmpty(d.Get("json_schema")),
		"readonly":               d.Get("readonly"),
		"required":               d.Get("required"),
		"searchable":             d.Get("searchable"),
		"visible":                d.Get("visible"),
		"maintainable":           d.Get("maintainable"),
		"editable":               d.Get("editable"),
		"creatable":              d.Get("creatable"),
		"boa_version":            utils.ValueIgnoreEmpty(d.Get("boa_version")),
	}

	return bodyParams
}

func resourceLayoutFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildCreateLayoutFieldBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster layout field: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster layout field: ID is not found in API response")
	}
	d.SetId(id)

	return resourceLayoutFieldRead(ctx, d, meta)
}

func resourceLayoutFieldRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster layout field: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	targetField := utils.PathSearch(fmt.Sprintf("[?id == '%s']|[0]", d.Id()), respBody, nil)
	if targetField == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cloud_pack_id", utils.PathSearch("cloud_pack_id", targetField, nil)),
		d.Set("cloud_pack_name", utils.PathSearch("cloud_pack_name", targetField, nil)),
		d.Set("dataclass_id", utils.PathSearch("dataclass_id", targetField, nil)),
		d.Set("cloud_pack_version", utils.PathSearch("cloud_pack_version", targetField, nil)),
		d.Set("field_key", utils.PathSearch("field_key", targetField, nil)),
		d.Set("name", utils.PathSearch("name", targetField, nil)),
		d.Set("description", utils.PathSearch("description", targetField, nil)),
		d.Set("en_description", utils.PathSearch("en_description", targetField, nil)),
		d.Set("default_value", utils.PathSearch("default_value", targetField, nil)),
		d.Set("en_default_value", utils.PathSearch("en_default_value", targetField, nil)),
		d.Set("field_type", utils.PathSearch("field_type", targetField, nil)),
		d.Set("extra_json", utils.PathSearch("extra_json", targetField, nil)),
		d.Set("field_tooltip", utils.PathSearch("field_tooltip", targetField, nil)),
		d.Set("en_field_tooltip", utils.PathSearch("en_field_tooltip", targetField, nil)),
		d.Set("json_schema", utils.PathSearch("json_schema", targetField, nil)),
		d.Set("is_built_in", utils.PathSearch("is_built_in", targetField, nil)),
		d.Set("readonly", utils.PathSearch("read_only", targetField, nil)),
		d.Set("required", utils.PathSearch("required", targetField, nil)),
		d.Set("searchable", utils.PathSearch("searchable", targetField, nil)),
		d.Set("visible", utils.PathSearch("visible", targetField, nil)),
		d.Set("maintainable", utils.PathSearch("maintainable", targetField, nil)),
		d.Set("editable", utils.PathSearch("editable", targetField, nil)),
		d.Set("creatable", utils.PathSearch("creatable", targetField, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", targetField, nil)),
		d.Set("creator_name", utils.PathSearch("creator_name", targetField, nil)),
		d.Set("modifier_id", utils.PathSearch("modifier_id", targetField, nil)),
		d.Set("modifier_name", utils.PathSearch("modifier_name", targetField, nil)),
		d.Set("create_time", utils.PathSearch("create_time", targetField, nil)),
		d.Set("update_time", utils.PathSearch("update_time", targetField, nil)),
		d.Set("wizard_id", utils.PathSearch("wizard_id", targetField, nil)),
		d.Set("aopworkflow_id", utils.PathSearch("aopworkflow_id", targetField, nil)),
		d.Set("aopworkflow_version_id", utils.PathSearch("aopworkflow_version_id", targetField, nil)),
		d.Set("playbook_id", utils.PathSearch("playbook_id", targetField, nil)),
		d.Set("playbook_version_id", utils.PathSearch("playbook_version_id", targetField, nil)),
		d.Set("boa_version", utils.PathSearch("boa_version", targetField, nil)),
		d.Set("version", utils.PathSearch("version", targetField, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateLayoutFieldBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                   d.Get("name"),
		"field_key":              d.Get("field_key"),
		"description":            utils.ValueIgnoreEmpty(d.Get("description")),
		"layout_id":              utils.ValueIgnoreEmpty(d.Get("layout_id")),
		"wizard_id":              utils.ValueIgnoreEmpty(d.Get("wizard_id")),
		"aopworkflow_id":         utils.ValueIgnoreEmpty(d.Get("aopworkflow_id")),
		"aopworkflow_version_id": utils.ValueIgnoreEmpty(d.Get("aopworkflow_version_id")),
		"playbook_id":            utils.ValueIgnoreEmpty(d.Get("playbook_id")),
		"playbook_version_id":    utils.ValueIgnoreEmpty(d.Get("playbook_version_id")),
		"default_value":          utils.ValueIgnoreEmpty(d.Get("default_value")),
		"display_type":           utils.ValueIgnoreEmpty(d.Get("display_type")),
		"field_type":             utils.ValueIgnoreEmpty(d.Get("field_type")),
		"extra_json":             utils.ValueIgnoreEmpty(d.Get("extra_json")),
		"field_tooltip":          utils.ValueIgnoreEmpty(d.Get("field_tooltip")),
		"json_schema":            utils.ValueIgnoreEmpty(d.Get("json_schema")),
		"readonly":               d.Get("readonly"),
		"required":               d.Get("required"),
		"searchable":             d.Get("searchable"),
		"visible":                d.Get("visible"),
		"maintainable":           d.Get("maintainable"),
		"editable":               d.Get("editable"),
		"creatable":              d.Get("creatable"),
		"boa_version":            utils.ValueIgnoreEmpty(d.Get("boa_version")),
	}

	return bodyParams
}

func resourceLayoutFieldUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields/{field_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{field_id}", d.Id())

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildUpdateLayoutFieldBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster layout field: %s", err)
	}

	return resourceLayoutFieldRead(ctx, d, meta)
}

func resourceLayoutFieldDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         []string{d.Id()},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster layout field: %s", err)
	}

	return nil
}

func resourceLayoutFieldImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want <workspace_id>/<id>, but got %s", importId)
	}

	d.SetId(importIdParts[1])
	return []*schema.ResourceData{d}, d.Set("workspace_id", importIdParts[0])
}
