package sdrs

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sdrs/v1/replications"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SDRS POST /v1/{project_id}/replications
// @API SDRS DELETE /v1/{project_id}/replications/{id}
// @API SDRS GET /v1/{project_id}/replications/{id}
// @API SDRS PUT /v1/{project_id}/replications/{id}
// @API SDRS GET /v1/{project_id}/jobs/{job_id}
func ResourceReplicationPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReplicationPairCreate,
		ReadContext:   resourceReplicationPairRead,
		UpdateContext: resourceReplicationPairUpdate,
		DeleteContext: resourceReplicationPairDelete,
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
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delete_target_volume": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"replication_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fault_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceReplicationPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	createOpts := replications.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		GroupID:     d.Get("group_id").(string),
		VolumeID:    d.Get("volume_id").(string),
	}

	n, err := replications.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating SDRS replication pair: %s", err)
	}

	createTimeoutSec := int(d.Timeout(schema.TimeoutCreate).Seconds())
	if err = replications.WaitForJobSuccess(client, createTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}

	replicationPairID, err := replications.GetJobEntity(client, n.JobID, "replication_pair_id")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(replicationPairID.(string))
	return resourceReplicationPairRead(ctx, d, meta)
}

func resourceReplicationPairRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	n, err := replications.Get(client, d.Id()).Extract()
	if err != nil {
		if errCode, ok := err.(golangsdk.ErrDefault400); ok {
			if resp, pErr := common.ParseErrorMsg(errCode.Body); pErr == nil && resp.ErrorCode == "SDRS.1608" {
				// `SDRS.1608` means replication pair not found
				return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
					"error retrieving SDRS replication pair")
			}
		}
		return diag.FromErr(err)
	}

	volumes := strings.Split(n.VolumeIDs, ",")
	if len(volumes) != 2 {
		return diag.Errorf("error retrieving volumes of replication pair: Invalid format. "+
			"except retrieving 2 volumes, but got %d", len(volumes))
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("group_id", n.GroupID),
		d.Set("volume_id", volumes[0]),
		d.Set("description", n.Description),
		d.Set("replication_model", n.ReplicaModel),
		d.Set("fault_level", n.FaultLevel),
		d.Set("status", n.Status),
		d.Set("target_volume_id", volumes[1]),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceReplicationPairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	if d.HasChange("name") {
		updateOpts := replications.UpdateOpts{
			Name: d.Get("name").(string),
		}
		_, err = replications.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating SDRS replication pair, %s", err)
		}
	}
	return resourceReplicationPairRead(ctx, d, meta)
}

func resourceReplicationPairDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	deleteOpts := replications.DeleteOpts{
		GroupID:      d.Get("group_id").(string),
		DeleteVolume: d.Get("delete_target_volume").(bool),
	}
	n, err := replications.Delete(client, d.Id(), deleteOpts).ExtractJobResponse()
	if err != nil {
		if errCode, ok := err.(golangsdk.ErrDefault400); ok {
			if resp, pErr := common.ParseErrorMsg(errCode.Body); pErr == nil && resp.ErrorCode == "SDRS.1608" {
				// `SDRS.1608` means replication pair not found
				return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
					"error deleting SDRS replication pair")
			}
		}
		return diag.FromErr(err)
	}

	deleteTimeoutSec := int(d.Timeout(schema.TimeoutDelete).Seconds())
	if err := replications.WaitForJobSuccess(client, deleteTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
