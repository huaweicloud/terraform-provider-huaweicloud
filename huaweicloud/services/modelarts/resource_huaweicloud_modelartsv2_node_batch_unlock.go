package modelarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v2NodeBatchUnlockNonUpdatableParams = []string{
	"pool_id",
	"node_names",
}

// @API ModelArts POST /v2/{project_id}/pools/{pool_name}/nodes/batch-unlock
func ResourceV2NodeBatchUnlock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NodeBatchUnlockCreate,
		ReadContext:   resourceV2NodeBatchUnlockRead,
		UpdateContext: resourceV2NodeBatchUnlockUpdate,
		DeleteContext: resourceV2NodeBatchUnlockDelete,

		CustomizeDiff: config.FlexibleForceNew(v2NodeBatchUnlockNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource nodes are located.`,
			},

			// Required parameters.
			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource pool ID to which the resource nodes belong.`,
			},
			"node_names": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The name list of resource nodes to be unlocked.`,
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
					}),
			},
		},
	}
}

func buildNodeBatchUnlockBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"nodeNames": d.Get("node_names").([]interface{}),
		// Whether locking or unlocking, the `actions` parameter must be an array.
		// And the array must contain only one element named `delete`.
		"actions": []interface{}{"delete"},
	}
}

func resourceV2NodeBatchUnlockCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/pools/{pool_name}/nodes/batch-unlock"
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{pool_name}", d.Get("pool_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildNodeBatchUnlockBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating node batch unlock: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return resourceV2NodeBatchUnlockRead(ctx, d, meta)
}

func resourceV2NodeBatchUnlockRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchUnlockUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchUnlockDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch unlocking the ModelArts nodes. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
