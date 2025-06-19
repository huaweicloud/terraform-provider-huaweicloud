package cbr

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableCheckpointSyncParams = []string{
	"vault_id",
	"auto_trigger",
}

// This resource uses the API for synchronizing hybrid cloud checkpoints.
// Due to the lack of test scenarios, this code is not tested and is not documented externally.
// Documentation is only stored in docs/incubating.

// @API CBR POST /v3/{project_id}/checkpoints/sync
// @API CBR GET /v3/{project_id}/operation-logs/{operation_log_id}
func ResourceCheckpointSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCheckpointSyncCreate,
		ReadContext:   resourceCheckpointSyncRead,
		UpdateContext: resourceCheckpointSyncUpdate,
		DeleteContext: resourceCheckpointSyncDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableCheckpointSyncParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level
region will be used.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the hybrid cloud vault to sync checkpoints to.`,
			},
			"auto_trigger": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether this checkpoint sync is automatically triggered.`,
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

func buildCheckpointSyncCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"sync": map[string]interface{}{
			"vault_id":     d.Get("vault_id"),
			"auto_trigger": d.Get("auto_trigger"),
		},
	}
}

func resourceCheckpointSyncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/checkpoints/sync"
		product = "cbr"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCheckpointSyncCreateBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CBR checkpoint sync: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	operationLogID := utils.PathSearch("sync.operation_log_id", respBody, "").(string)
	if operationLogID == "" {
		return diag.Errorf("error creating CBR checkpoint sync: Operation Log ID is not found in API response")
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID for CBR checkpoint sync: %s", err)
	}
	d.SetId(id)

	if err := waitingForSyncTaskSuccess(ctx, client, d.Timeout(schema.TimeoutCreate), operationLogID); err != nil {
		return diag.Errorf("error waiting for CBR checkpoint sync to complete: %s", err)
	}

	return resourceCheckpointSyncRead(ctx, d, meta)
}

func resourceCheckpointSyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceCheckpointSyncUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceCheckpointSyncDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to synchronize a CBR checkpoint. Deleting this 
resource will not change the current checkpoint synchronization result, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
