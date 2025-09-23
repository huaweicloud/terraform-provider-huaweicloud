package workspace

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

var appHdaBatchUpgradeNonUpdatableParams = []string{"server_sids"}

// @API Workspace PATCH /v1/{project_id}/app-servers/access-agent/actions/upgrade
func ResourceAppHdaBatchUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppHdaBatchUpgradeCreate,
		ReadContext:   resourceAppHdaBatchUpgradeRead,
		UpdateContext: resourceAppHdaBatchUpgradeUpdate,
		DeleteContext: resourceAppHdaBatchUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(appHdaBatchUpgradeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the servers to be upgraded are located.`,
			},

			// Required parameters.
			"server_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of server IDs to be upgraded HDA in batches.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildAppHdaBatchUpgradeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"items": d.Get("server_ids").([]interface{}),
	}
}

func resourceAppHdaBatchUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpUrl := "v1/{project_id}/app-servers/access-agent/actions/upgrade"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAppHdaBatchUpgradeBodyParams(d),
	}

	_, err = client.Request("PATCH", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error executing HDA batch upgrade: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceAppHdaBatchUpgradeRead(ctx, d, meta)
}

func resourceAppHdaBatchUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppHdaBatchUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppHdaBatchUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for batch upgrading HDA versions of servers. Deleting this
resource will not clear the corresponding upgrade request record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
