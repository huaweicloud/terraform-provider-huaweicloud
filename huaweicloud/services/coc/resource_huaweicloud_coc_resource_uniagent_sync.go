package coc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var resourceUniAgentSyncNonUpdatableParams = []string{"resource_infos", "resource_infos.*.region_id",
	"resource_infos.*.resource_id", "vendor"}

// @API COC POST /v1/resources/uniagent/sync
func ResourceResourceUniAgentSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceUniAgentSyncCreate,
		ReadContext:   resourceResourceUniAgentSyncRead,
		UpdateContext: resourceResourceUniAgentSyncUpdate,
		DeleteContext: resourceResourceUniAgentSyncDelete,

		CustomizeDiff: config.FlexibleForceNew(resourceUniAgentSyncNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"resource_infos": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"vendor": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildResourceUniAgentSyncCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resource_infos": buildResourceUniAgentSyncResourceInfosCreateOpts(d.Get("resource_infos")),
		"vendor":         utils.ValueIgnoreEmpty(d.Get("vendor")),
	}

	return bodyParams
}

func buildResourceUniAgentSyncResourceInfosCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"region_id":   raw["region_id"],
				"resource_id": raw["resource_id"],
			}
		}
		return params
	}

	return nil
}

func resourceResourceUniAgentSyncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/resources/uniagent/sync"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResourceUniAgentSyncCreateOpts(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error syncing the COC uniagent resource: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	syncId := utils.PathSearch("data", createRespBody, "").(string)
	if syncId == "" {
		return diag.Errorf("unable to find the syncing uniagent resource ID from the API response")
	}
	d.SetId(syncId)

	return resourceResourceUniAgentSyncRead(ctx, d, meta)
}

func resourceResourceUniAgentSyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourceUniAgentSyncUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourceUniAgentSyncDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting resource uniagent sync resource is not supported. The resource uniagent sync resource is" +
		" only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
