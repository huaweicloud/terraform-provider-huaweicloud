package cbr

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableCheckpointCopyParams = []string{
	"vault_id",
	"destination_project_id",
	"destination_region",
	"destination_vault_id",
	"enable_acceleration",
	"auto_trigger",
}

// @API CBR POST /v3/{project_id}/checkpoints/replicate
func ResourceCheckpointCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCheckpointCopyCreate,
		ReadContext:   resourceCheckpointCopyRead,
		UpdateContext: resourceCheckpointCopyUpdate,
		DeleteContext: resourceCheckpointCopyDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableCheckpointCopyParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the source vault.`,
			},
			"destination_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the destination project to which the backup is to be copied.`,
			},
			"destination_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the destination region to which the backup is to be copied.`,
			},
			"destination_vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the destination vault to which the backup is to be copied.`,
			},
			"auto_trigger": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to automatically trigger the replication.`,
			},
			"enable_acceleration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable acceleration for cross-region replication.`,
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

func buildCheckpointCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"replicate": map[string]interface{}{
			"auto_trigger":           d.Get("auto_trigger"),
			"destination_project_id": d.Get("destination_project_id"),
			"destination_region":     d.Get("destination_region"),
			"destination_vault_id":   d.Get("destination_vault_id"),
			"enable_acceleration":    d.Get("enable_acceleration"),
			"vault_id":               d.Get("vault_id"),
		},
	}
}

func resourceCheckpointCopyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/checkpoints/replicate"
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
		JSONBody:         buildCheckpointCopyBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error copying CBR checkpoint: %s", err)
	}

	// Generate UUID as resource ID
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func resourceCheckpointCopyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This is a one-time operation resource, no need to do anything on read
	return nil
}

func resourceCheckpointCopyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This is a one-time operation resource, update is not supported
	return nil
}

func resourceCheckpointCopyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to copy a CBR checkpoint. Deleting this 
resource will not change the current copy checkpoint result, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
