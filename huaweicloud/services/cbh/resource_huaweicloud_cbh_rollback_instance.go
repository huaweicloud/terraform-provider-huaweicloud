package cbh

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

// @API CBH POST /v2/{project_id}/cbs/instance/rollback
func ResourceRollbackInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRollbackInstanceCreate,
		UpdateContext: resourceRollbackInstanceUpdate,
		ReadContext:   resourceRollbackInstanceRead,
		DeleteContext: resourceRollbackInstanceDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"server_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
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

func resourceRollbackInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v2/{project_id}/cbs/instance/rollback"
		product  = "cbh"
		serverId = d.Get("server_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"server_id": serverId,
		},
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error rolling back CBH instance: %s", err)
	}

	d.SetId(serverId)

	return resourceRollbackInstanceRead(ctx, d, meta)
}

func resourceRollbackInstanceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceRollbackInstanceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceRollbackInstanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to rollback CBH instance.
Deleting this resource will not recover the rollback CBH instance, but will only remove the resource information from
the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
