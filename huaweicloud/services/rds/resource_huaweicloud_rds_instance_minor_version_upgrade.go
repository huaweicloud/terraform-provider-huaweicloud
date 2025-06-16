package rds

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

var dbInstanceMinorVersionUpgradeNonUpdatableParams = []string{"instance_id", "is_delayed"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/db-upgrade
func ResourceRdsInstanceMinorVersionUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsInstanceMinorVersionUpgradeCreate,
		ReadContext:   resourceRdsInstanceMinorVersionUpgradeRead,
		UpdateContext: resourceRdsInstanceMinorVersionUpgradeUpdate,
		DeleteContext: resourceRdsInstanceMinorVersionUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(dbInstanceMinorVersionUpgradeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_delayed": {
				Type:     schema.TypeBool,
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

func resourceRdsInstanceMinorVersionUpgradeCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	const (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db-upgrade"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	url := client.Endpoint + httpUrl
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)
	url = strings.ReplaceAll(url, "{instance_id}", instanceID)

	body := map[string]interface{}{
		"is_delayed": d.Get("is_delayed"),
	}

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opts.JSONBody = utils.RemoveNil(body)

	_, err = client.Request("POST", url, &opts)
	if err != nil {
		return diag.Errorf("error upgrading kernel minor version for instance (%s): %s", instanceID, err)
	}

	d.SetId(instanceID)

	return nil
}

func resourceRdsInstanceMinorVersionUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsInstanceMinorVersionUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsInstanceMinorVersionUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting kernel upgrade resource is not supported. This resource is only removed from the state.",
		},
	}
}
