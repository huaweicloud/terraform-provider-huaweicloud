package sdrs

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/drill"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SDRS DELETE /v1/{project_id}/disaster-recovery-drills/{id}
// @API SDRS GET /v1/{project_id}/disaster-recovery-drills/{id}
// @API SDRS PUT /v1/{project_id}/disaster-recovery-drills/{id}
// @API SDRS POST /v1/{project_id}/disaster-recovery-drills
func ResourceDrill() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrillCreate,
		ReadContext:   resourceDrillRead,
		UpdateContext: resourceDrillUpdate,
		DeleteContext: resourceDrillDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"drill_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDrillCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	createOpts := drill.CreateOpts{
		Name:       d.Get("name").(string),
		GroupID:    d.Get("group_id").(string),
		DrillVpcID: d.Get("drill_vpc_id").(string),
	}

	n, err := drill.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating SDRS DR drill: %s", err)
	}

	createTimeoutSec := int(d.Timeout(schema.TimeoutCreate).Seconds())
	if err = drill.WaitForJobSuccess(client, createTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}

	drillID, err := drill.GetJobEntity(client, n.JobID, "disaster_recovery_drill_id")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(drillID.(string))
	return resourceDrillRead(ctx, d, meta)
}

func resourceDrillRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	n, err := drill.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error.code", "SDRS.1902"),
			"error retrieving SDRS DR drill",
		)
	}

	mErr := multierror.Append(
		nil,
		d.Set("name", n.Name),
		d.Set("group_id", n.GroupID),
		d.Set("drill_vpc_id", n.DrillVpcID),
		d.Set("status", n.Status),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDrillUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	if d.HasChange("name") {
		updateOpts := drill.UpdateOpts{
			Name: d.Get("name").(string),
		}
		_, err = drill.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating SDRS DR drill name, %s", err)
		}
	}
	return resourceDrillRead(ctx, d, meta)
}

func resourceDrillDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	n, err := drill.Delete(client, d.Id()).ExtractJobResponse()
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error.code", "SDRS.1913"),
			"error deleting SDRS DR drill",
		)
	}

	deleteTimeoutSec := int(d.Timeout(schema.TimeoutDelete).Seconds())
	if err := drill.WaitForJobSuccess(client, deleteTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
