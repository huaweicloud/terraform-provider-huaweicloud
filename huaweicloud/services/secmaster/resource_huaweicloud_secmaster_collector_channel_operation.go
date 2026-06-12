package secmaster

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

var collectorChannelOptNonUpdatableParams = []string{"workspace_id", "channel_id", "action"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/collector/channels/{channel_id}/operation
func ResourceCollectorChannelOperation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectorChannelOperationCreate,
		UpdateContext: resourceCollectorChannelOperationUpdate,
		ReadContext:   resourceCollectorChannelOperationRead,
		DeleteContext: resourceCollectorChannelOperationDelete,

		CustomizeDiff: config.FlexibleForceNew(collectorChannelOptNonUpdatableParams),

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
			"channel_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
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

func resourceCollectorChannelOperationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		channelId     = d.Get("channel_id").(string)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/collector/channels/{channel_id}/operation"
	)
	createClient, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := createClient.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", createClient.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", d.Get("workspace_id").(string))
	createPath = strings.ReplaceAll(createPath, "{channel_id}", channelId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createOpt.JSONBody = map[string]interface{}{
		"action": d.Get("action"),
	}

	_, err = createClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating collector channel operation: %s", err)
	}

	d.SetId(channelId)

	return nil
}

func resourceCollectorChannelOperationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCollectorChannelOperationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCollectorChannelOperationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for collector channel operation. Deleting this resource will not change
		the status of the currently collector channel operation, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
