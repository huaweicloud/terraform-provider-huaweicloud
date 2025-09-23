package codeartspipeline

import (
	"context"
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

var basicPluginNonUpdatableParams = []string{
	"plugin_name", "business_type", "business_type_display_name", "runtime_attribution", "plugin_composition_type",
}

// @API CodeArtsPipeline POST /v3/{domain_id}/extension/info/add
// @API CodeArtsPipeline POST /v1/{domain_id}/agent-plugin/query-all
// @API CodeArtsPipeline POST /v3/{domain_id}/extension/info/update
// @API CodeArtsPipeline DELETE /v3/{domain_id}/extension/info/delete
func ResourceCodeArtsPipelineBasicPlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineBasicPluginCreate,
		ReadContext:   resourcePipelineBasicPluginRead,
		UpdateContext: resourcePipelineBasicPluginUpdate,
		DeleteContext: resourcePipelineBasicPluginDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(basicPluginNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"plugin_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the basic plugin name.`,
			},
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the service type.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the display name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the basic plugin description.`,
			},
			"icon_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the icon URL.`,
			},
			"business_type_display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the display name of service type.`,
			},
			"runtime_attribution": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the runtime attributes.`,
			},
			"plugin_composition_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the combination extension type.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"maintainers": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the maintenance engineer.`,
			},
			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique ID.`,
			},
		},
	}
}

func resourcePipelineBasicPluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v3/{domain_id}/extension/info/add"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{domain_id}", cfg.DomainID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineBasicPluginBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline basic plugin: %s", err)
	}

	d.SetId(d.Get("plugin_name").(string))

	return resourcePipelineBasicPluginRead(ctx, d, meta)
}

func buildCreateOrUpdatePipelineBasicPluginBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"plugin_name":                d.Get("plugin_name"),
		"business_type":              d.Get("business_type"),
		"display_name":               d.Get("display_name"),
		"description":                d.Get("description"),
		"business_type_display_name": utils.ValueIgnoreEmpty(d.Get("business_type_display_name")),
		"icon_url":                   utils.ValueIgnoreEmpty(d.Get("icon_url")),
		"runtime_attribution":        utils.ValueIgnoreEmpty(d.Get("runtime_attribution")),
		"plugin_composition_type":    utils.ValueIgnoreEmpty(d.Get("plugin_composition_type")),
	}

	return bodyParams
}

func resourcePipelineBasicPluginRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getRespBody, err := GetPipelineBasicPlugin(client, cfg.DomainID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Pipeline basic plugin")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plugin_name", utils.PathSearch("plugin_name", getRespBody, nil)),
		d.Set("business_type", utils.PathSearch("business_type", getRespBody, nil)),
		d.Set("display_name", utils.PathSearch("display_name", getRespBody, nil)),
		d.Set("business_type_display_name", utils.PathSearch("business_type_display_name", getRespBody, nil)),
		d.Set("icon_url", utils.PathSearch("icon_url", getRespBody, nil)),
		d.Set("runtime_attribution", utils.PathSearch("runtime_attribution", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("plugin_composition_type", utils.PathSearch("plugin_composition_type", getRespBody, nil)),
		d.Set("maintainers", utils.PathSearch("maintainers", getRespBody, nil)),
		d.Set("unique_id", utils.PathSearch("unique_id", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetPipelineBasicPlugin(client *golangsdk.ServiceClient, domainId, id string) (interface{}, error) {
	getHttpUrl := "v1/{domain_id}/agent-plugin/query-all?limit=1&offset=0"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"plugin_name":   id,
			"business_type": []string{"Build", "Gate", "Deploy", "Test", "Normal"},
		},
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	plugin := utils.PathSearch("data[0]", getRespBody, nil)
	if plugin == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return plugin, nil
}

func resourcePipelineBasicPluginUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v3/{domain_id}/extension/info/update"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{domain_id}", cfg.DomainID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineBasicPluginBodyParams(d)),
	}

	_, err = client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating CodeArts Pipeline basic plugin: %s", err)
	}

	return resourcePipelineBasicPluginRead(ctx, d, meta)
}

func resourcePipelineBasicPluginDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v3/{domain_id}/extension/info/delete?plugin_name={plugin_name}&type=all"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{domain_id}", cfg.DomainID)
	deletePath = strings.ReplaceAll(deletePath, "{plugin_name}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DEVPIPE.30011001"),
			"error deleting CodeArts Pipeline basic plugin")
	}

	return nil
}
