package coc

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var groupSyncNonUpdatableParams = []string{"group_id", "cloud_service_name", "type"}

// @API COC POST /v1/groups/{id}/sync
func ResourceGroupSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupSyncCreate,
		ReadContext:   resourceGroupSyncRead,
		UpdateContext: resourceGroupSyncUpdate,
		DeleteContext: resourceGroupSyncDelete,

		CustomizeDiff: config.FlexibleForceNew(groupSyncNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
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

func buildGroupSyncCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"provider": utils.ValueIgnoreEmpty(d.Get("cloud_service_name")),
		"type":     utils.ValueIgnoreEmpty(d.Get("type")),
	}

	return bodyParams
}

func resourceGroupSyncCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/groups/{id}/sync"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	groupID := d.Get("group_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{id}", groupID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildGroupSyncCreateOpts(d)),
	}

	createSyncResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error syncing the COC group (%s): %s", groupID, err)
	}

	createSyncRespBody, err := utils.FlattenResponse(createSyncResp)
	if err != nil {
		return diag.FromErr(err)
	}

	syncID := utils.PathSearch("data", createSyncRespBody, "").(string)
	if syncID == "" {
		return diag.Errorf("unable to find the syncing COC group sync ID from the API response")
	}

	d.SetId(syncID)

	return nil
}

func resourceGroupSyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGroupSyncUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGroupSyncDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting group sync operation resource is not supported. The group sync operation resource is only" +
		" removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
