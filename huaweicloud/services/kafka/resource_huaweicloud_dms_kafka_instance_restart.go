package kafka

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API Kafka POST /v2/{project_id}/instances/action
// @API Kafka GET /v2/{project_id}/instances/{instance_id}
func ResourceDmsKafkaInstanceRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaInstanceRestartCreate,
		ReadContext:   resourceDmsKafkaInstanceRestartRead,
		DeleteContext: resourceDmsKafkaInstanceRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the Kafka instance.`,
			},
		},
	}
}

func resourceDmsKafkaInstanceRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	if err = restartKafkaInstance(ctx, d.Timeout(schema.TimeoutCreate), client, instanceID); err != nil {
		return diag.Errorf("error restarting instance: %s", err)
	}

	d.SetId(instanceID)

	return nil
}

func resourceDmsKafkaInstanceRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsKafkaInstanceRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting resource is not supported. The resource is only removed from the state, the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
