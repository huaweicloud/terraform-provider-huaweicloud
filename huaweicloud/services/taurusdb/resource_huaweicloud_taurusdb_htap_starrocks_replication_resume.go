package taurusdb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var starrocksReplicationResumeNoneUpdatableParams = []string{
	"instance_id", "task_name",
}

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/resume
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication
func ResourceTaurusDBHtapStarrocksReplicationResume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksReplicationResumeCreate,
		ReadContext:   resourceTaurusDBHtapStarrocksReplicationResumeRead,
		UpdateContext: resourceTaurusDBHtapStarrocksReplicationResumeUpdate,
		DeleteContext: resourceTaurusDBHtapStarrocksReplicationResumeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(starrocksReplicationResumeNoneUpdatableParams),

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
			"task_name": {
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

func resourceTaurusDBHtapStarrocksReplicationResumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	taskName := d.Get("task_name").(string)

	err = resumeHtapReplicationTask(ctx, client, instanceId, taskName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error resuming HTAP insatnce (%s) replication task (%s): %s", instanceId, taskName, err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(id.String())

	return nil
}

func resourceTaurusDBHtapStarrocksReplicationResumeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksReplicationResumeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksReplicationResumeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting replication resume resource is not supported. The replication resume resource is only removed" +
		" from the state, the StarRocks instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
