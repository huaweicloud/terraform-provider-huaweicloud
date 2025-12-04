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

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials/{asset_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials/{asset_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials/{asset_id}
func ResourceOperationConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOperationConnectionCreate,
		ReadContext:   resourceOperationConnectionRead,
		UpdateContext: resourceOperationConnectionUpdate,
		DeleteContext: resourceOperationConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOperationConnectionImportState,
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
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Field `description` can be modify to empty, so no `Computed` attribute is added here.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
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
			"update_time": {
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
			"defense_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_enterprise_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_enterprise_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateOperationConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":                 d.Get("name"),
		"component_id":         d.Get("component_id"),
		"component_version_id": d.Get("component_version_id"),
		"config":               d.Get("config"),
		"description":          utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceOperationConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials"
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
		JSONBody:         utils.RemoveNil(buildCreateOperationConnectionBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster operation connection: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("asset.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster operation connection: ID is not found in API response")
	}
	d.SetId(id)

	return resourceOperationConnectionRead(ctx, d, meta)
}

func resourceOperationConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials/{asset_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{asset_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "code", "SecMaster.20041002"),
			"error retrieving SecMaster operation connection",
		)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("project_id", utils.PathSearch("asset.project_id", respBody, nil)),
		d.Set("workspace_id", utils.PathSearch("asset.workspace_id", respBody, nil)),
		d.Set("name", utils.PathSearch("asset.name", respBody, nil)),
		d.Set("component_id", utils.PathSearch("asset.component_id", respBody, nil)),
		d.Set("component_name", utils.PathSearch("asset.component_name", respBody, nil)),
		d.Set("component_version_id", utils.PathSearch("asset.component_version_id", respBody, nil)),
		d.Set("type", utils.PathSearch("asset.type", respBody, nil)),
		d.Set("status", utils.PathSearch("asset.status", respBody, nil)),
		d.Set("config", utils.PathSearch("asset.config", respBody, nil)),
		d.Set("description", utils.PathSearch("asset.description", respBody, nil)),
		d.Set("enabled", utils.PathSearch("asset.enabled", respBody, nil)),
		d.Set("create_time", utils.PathSearch("asset.create_time", respBody, nil)),
		d.Set("creator_id", utils.PathSearch("asset.creator_id", respBody, nil)),
		d.Set("creator_name", utils.PathSearch("asset.creator_name", respBody, nil)),
		d.Set("update_time", utils.PathSearch("asset.update_time", respBody, nil)),
		d.Set("modifier_id", utils.PathSearch("asset.modifier_id", respBody, nil)),
		d.Set("modifier_name", utils.PathSearch("asset.modifier_name", respBody, nil)),
		d.Set("defense_type", utils.PathSearch("asset.defense_type", respBody, nil)),
		d.Set("target_project_id", utils.PathSearch("asset.target_project_id", respBody, nil)),
		d.Set("target_project_name", utils.PathSearch("asset.target_project_name", respBody, nil)),
		d.Set("target_enterprise_id", utils.PathSearch("asset.target_enterprise_id", respBody, nil)),
		d.Set("target_enterprise_name", utils.PathSearch("asset.target_enterprise_name", respBody, nil)),
		d.Set("region_id", utils.PathSearch("asset.region_id", respBody, nil)),
		d.Set("region_name", utils.PathSearch("asset.region_name", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateOperationConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":                 d.Get("name"),
		"component_id":         d.Get("component_id"),
		"component_version_id": d.Get("component_version_id"),
		"config":               d.Get("config"),
		"description":          d.Get("description"),
	}
}

func resourceOperationConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials/{asset_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{asset_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         buildUpdateOperationConnectionBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster operation connection: %s", err)
	}

	return resourceOperationConnectionRead(ctx, d, meta)
}

func resourceOperationConnectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials/{asset_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{asset_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster operation connection: %s", err)
	}

	return nil
}

func resourceOperationConnectionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want <workspace_id>/<id>, but got %s", importId)
	}

	d.SetId(importIdParts[1])
	return []*schema.ResourceData{d}, d.Set("workspace_id", importIdParts[0])
}
