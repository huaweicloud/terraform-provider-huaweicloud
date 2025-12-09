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

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/modules
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/modules/{module_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/modules/{module_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/modules/{module_id}
func ResourceModule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModuleCreate,
		ReadContext:   resourceModuleRead,
		UpdateContext: resourceModuleUpdate,
		DeleteContext: resourceModuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceModuleImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"workspace_id",
			"project_id",
		}),

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
			// update
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// update
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// update
			"module_json": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// update
			"module_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// update
			"metric_ids": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// update
			"thumbnail": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// update
			"data_query": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// update
			"boa_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// update
			"cloud_pack_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// update
			"cloud_pack_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// update
			"cloud_pack_version": {
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
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Computed: true}),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"en_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"en_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_built_in": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateModuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":         d.Get("name"),
		"workspace_id": d.Get("workspace_id"),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"project_id":   utils.ValueIgnoreEmpty(d.Get("project_id")),
		"module_json":  utils.ValueIgnoreEmpty(d.Get("module_json")),
		"module_type":  utils.ValueIgnoreEmpty(d.Get("module_type")),
		"metric_ids":   utils.ValueIgnoreEmpty(d.Get("metric_ids")),
		"thumbnail":    utils.ValueIgnoreEmpty(d.Get("thumbnail")),
		"data_query":   utils.ValueIgnoreEmpty(d.Get("data_query")),
		"boa_version":  utils.ValueIgnoreEmpty(d.Get("boa_version")),
	}
}

func whetherHasUpdateFields(d *schema.ResourceData) bool {
	updateFields := []string{"cloud_pack_id", "cloud_pack_name", "cloud_pack_version"}
	for _, field := range updateFields {
		if _, ok := d.GetOk(field); ok {
			return true
		}
	}
	return false
}

func resourceModuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/modules"
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
		JSONBody:         utils.RemoveNil(buildCreateModuleBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster module: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster module: ID is not found in API response")
	}
	d.SetId(id)

	// Some fields can be configured only through update API.
	if whetherHasUpdateFields(d) {
		err := updateModule(client, d)
		if err != nil {
			return diag.Errorf("error updating SecMaster module: %s", err)
		}
	}

	return resourceModuleRead(ctx, d, meta)
}

func resourceModuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/modules/{module_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{module_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "SecMaster.20097001"),
			"error retrieving SecMaster module",
		)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cloud_pack_id", utils.PathSearch("data.cloud_pack_id", respBody, nil)),
		d.Set("cloud_pack_name", utils.PathSearch("data.cloud_pack_name", respBody, nil)),
		d.Set("cloud_pack_version", utils.PathSearch("data.cloud_pack_version", respBody, nil)),
		d.Set("create_time", utils.PathSearch("data.create_time", respBody, nil)),
		d.Set("creator_id", utils.PathSearch("data.creator_id", respBody, nil)),
		d.Set("description", utils.PathSearch("data.description", respBody, nil)),
		d.Set("en_description", utils.PathSearch("data.en_description", respBody, nil)),
		d.Set("module_json", utils.PathSearch("data.module_json", respBody, nil)),
		d.Set("name", utils.PathSearch("data.name", respBody, nil)),
		d.Set("en_name", utils.PathSearch("data.en_name", respBody, nil)),
		d.Set("project_id", utils.PathSearch("data.project_id", respBody, nil)),
		d.Set("workspace_id", utils.PathSearch("data.workspace_id", respBody, nil)),
		d.Set("update_time", utils.PathSearch("data.update_time", respBody, nil)),
		d.Set("thumbnail", utils.PathSearch("data.thumbnail", respBody, nil)),
		d.Set("module_type", utils.PathSearch("data.module_type", respBody, nil)),
		d.Set("tag", utils.PathSearch("data.tag", respBody, nil)),
		d.Set("is_built_in", utils.PathSearch("data.is_built_in", respBody, nil)),
		d.Set("data_query", utils.PathSearch("data.data_query", respBody, nil)),
		d.Set("boa_version", utils.PathSearch("data.boa_version", respBody, nil)),
		d.Set("version", utils.PathSearch("data.version", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateModuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":               d.Get("name"),
		"cloud_pack_id":      utils.ValueIgnoreEmpty(d.Get("cloud_pack_id")),
		"cloud_pack_name":    utils.ValueIgnoreEmpty(d.Get("cloud_pack_name")),
		"cloud_pack_version": utils.ValueIgnoreEmpty(d.Get("cloud_pack_version")),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
		"module_json":        utils.ValueIgnoreEmpty(d.Get("module_json")),
		"thumbnail":          utils.ValueIgnoreEmpty(d.Get("thumbnail")),
		"module_type":        utils.ValueIgnoreEmpty(d.Get("module_type")),
		"metric_ids":         utils.ValueIgnoreEmpty(d.Get("metric_ids")),
		"data_query":         utils.ValueIgnoreEmpty(d.Get("data_query")),
		"boa_version":        utils.ValueIgnoreEmpty(d.Get("boa_version")),
	}
}

func updateModule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/soc/modules/{module_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{module_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildUpdateModuleBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceModuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	if err := updateModule(client, d); err != nil {
		return diag.Errorf("error updating SecMaster module: %s", err)
	}

	return resourceModuleRead(ctx, d, meta)
}

func resourceModuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/modules/{module_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{module_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster module: %s", err)
	}

	return nil
}

func resourceModuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want <workspace_id>/<id>, but got %s", importId)
	}

	d.SetId(importIdParts[1])
	return []*schema.ResourceData{d}, d.Set("workspace_id", importIdParts[0])
}
