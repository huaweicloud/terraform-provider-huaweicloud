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

var backupStopNonUpdatableParams = []string{"instance_id"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/backups/stop
func ResourceRdsBackupStop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsBackupStopCreate,
		ReadContext:   resourceRdsBackupStopRead,
		UpdateContext: resourceRdsBackupStopUpdate,
		DeleteContext: resourceRdsBackupStopDelete,

		CustomizeDiff: config.FlexibleForceNew(backupStopNonUpdatableParams),

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceRdsBackupStopCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/backups/stop"
		product    = "rds"
		instanceID = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	url := client.Endpoint + httpUrl
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)
	url = strings.ReplaceAll(url, "{instance_id}", instanceID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", url, &opts)
	if err != nil {
		return diag.Errorf("error stopping backup for instance (%s): %s", instanceID, err)
	}

	d.SetId(instanceID)

	return nil
}

func resourceRdsBackupStopRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsBackupStopUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsBackupStopDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS backup stop is not supported. " +
		"The backup stop resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
