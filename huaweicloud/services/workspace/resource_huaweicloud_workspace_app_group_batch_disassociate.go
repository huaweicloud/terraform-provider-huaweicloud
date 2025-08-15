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

var appServerGroupBatchDisassociateNonUpdateParams = []string{"server_group_id"}

// @API Workspace POST /v1/{project_id}/app-groups/actions/disassociate-app-group
func ResourceAppServerGroupBatchDisassociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppServerGroupBatchDisassociateCreate,
		ReadContext:   resourceAppServerGroupBatchDisassociateRead,
		UpdateContext: resourceAppServerGroupBatchDisassociateUpdate,
		DeleteContext: resourceAppServerGroupBatchDisassociateDelete,

		CustomizeDiff: config.FlexibleForceNew(appServerGroupBatchDisassociateNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the server group is located.`,
			},

			// Required parameter(s).
			"server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the server group to disassociate all application groups.`,
			},

			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceAppServerGroupBatchDisassociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		httpUrl       = "v1/{project_id}/app-groups/actions/disassociate-app-group"
		serverGroupId = d.Get("server_group_id").(string)
	)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createPath += "?server_group_id=" + serverGroupId

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error disassociating application groups from server group (%s): %s", serverGroupId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	return resourceAppServerGroupBatchDisassociateRead(ctx, d, meta)
}

func resourceAppServerGroupBatchDisassociateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppServerGroupBatchDisassociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppServerGroupBatchDisassociateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to disassociate all application groups from a
specified server group. Deleting this resource will not clear the corresponding request record, but will only remove
the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
