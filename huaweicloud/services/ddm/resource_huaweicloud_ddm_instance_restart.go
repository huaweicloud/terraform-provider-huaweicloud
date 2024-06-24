package ddm

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DDM POST /v1/{project_id}/instances/{instance_id}/action
// @API DDM GET /v1/{project_id}/instances/{instance_id}
func ResourceDdmInstanceRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDMInstanceRestartCreate,
		ReadContext:   resourceDDMInstanceRestartRead,
		DeleteContext: resourceDDMInstanceRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

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
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDDMInstanceRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDM Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	restartType := d.Get("type").(string)
	err = restartDdmInstance(ctx, client, instanceID, restartType, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error restarting instance: %s", err)
	}

	d.SetId(instanceID)

	return nil
}

func resourceDDMInstanceRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDDMInstanceRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restart resource is not supported. The restart resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
