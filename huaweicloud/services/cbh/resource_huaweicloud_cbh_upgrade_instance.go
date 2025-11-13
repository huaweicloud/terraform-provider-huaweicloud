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

// @API CBH POST /v2/{project_id}/cbs/instance/upgrade
func ResourceUpgradeInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUpgradeInstanceCreate,
		UpdateContext: resourceUpgradeInstanceUpdate,
		ReadContext:   resourceUpgradeInstanceRead,
		DeleteContext: resourceUpgradeInstanceDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"server_id",
			"upgrade_time",
			"cancel",
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
			"upgrade_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cancel": {
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

func buildCancelValue(d *schema.ResourceData) interface{} {
	cancelString := d.Get("cancel").(string)
	if cancelString == "true" {
		return true
	}

	if cancelString == "false" {
		return false
	}

	return nil
}

func buildUpgradeInstanceRequestOpt(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"server_id":    d.Get("server_id"),
		"upgrade_time": d.Get("upgrade_time"),
		"cancel":       buildCancelValue(d),
	}
}

func resourceUpgradeInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v2/{project_id}/cbs/instance/upgrade"
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
		JSONBody:         utils.RemoveNil(buildUpgradeInstanceRequestOpt(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error upgrading CBH instance: %s", err)
	}

	d.SetId(serverId)

	return resourceUpgradeInstanceRead(ctx, d, meta)
}

func resourceUpgradeInstanceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceUpgradeInstanceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceUpgradeInstanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to upgrade CBH instance.
Deleting this resource will not recover the upgrade CBH instance, but will only remove the resource information from
the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
