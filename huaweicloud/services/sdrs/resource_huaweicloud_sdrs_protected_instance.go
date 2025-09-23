package sdrs

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectedinstances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SDRS POST /v1/{project_id}/protected-instances/{id}/tags/action
// @API SDRS GET /v1/{project_id}/protected-instances/{id}
// @API SDRS PUT /v1/{project_id}/protected-instances/{id}
// @API SDRS DELETE /v1/{project_id}/protected-instances/{id}
// @API SDRS POST /v1/{project_id}/protected-instances
// @API SDRS GET /v1/{project_id}/protected-instances/{protected_instance_id}/tags
// @API SDRS POST /v1/{project_id}/protected-instances/{protected_instance_id}/tags/action
// @API SDRS GET /v1/{project_id}/jobs/{job_id}
func ResourceProtectedInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtectedInstanceCreate,
		ReadContext:   resourceProtectedInstanceRead,
		UpdateContext: resourceProtectedInstanceUpdate,
		DeleteContext: resourceProtectedInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"primary_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"primary_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delete_target_server": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_target_eip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"target_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceProtectedInstanceCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	createOpts := protectedinstances.CreateOpts{
		GroupID:     d.Get("group_id").(string),
		ServerID:    d.Get("server_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ClusterID:   d.Get("cluster_id").(string),
		SubnetID:    d.Get("primary_subnet_id").(string),
		IpAddress:   d.Get("primary_ip_address").(string),
	}
	n, err := protectedinstances.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating SDRS protected instance: %s", err)
	}

	createTimeoutSec := int(d.Timeout(schema.TimeoutCreate).Seconds())
	if err = protectedinstances.WaitForJobSuccess(client, createTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}

	instanceID, err := protectedinstances.GetJobEntity(client, n.JobID, "protected_instance_id")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(instanceID.(string))

	// add tags
	if err := utils.CreateResourceTags(client, d, "protected-instances", d.Id()); err != nil {
		return diag.Errorf("error setting tags of SDRS protected instance %s: %s", d.Id(), err)
	}

	return resourceProtectedInstanceRead(ctx, d, meta)
}

func resourceProtectedInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	n, err := protectedinstances.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error.code", "SDRS.1320"),
			"error retrieving SDRS protected instance",
		)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("group_id", n.GroupID),
		d.Set("server_id", n.SourceServer),
		d.Set("description", n.Description),
		d.Set("target_server", n.TargetServer),
		d.Set("tags", d.Get("tags")),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}

	if err := utils.SetResourceTagsToState(d, client, "protected-instances", d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceProtectedInstanceUpdate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	if d.HasChange("name") {
		updateOpts := protectedinstances.UpdateOpts{
			Name: d.Get("name").(string),
		}
		_, err = protectedinstances.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating SDRS protected instance name, %s", err)
		}
	}

	if d.HasChange("tags") {
		if err := utils.UpdateResourceTags(client, d, "protected-instances", d.Id()); err != nil {
			return diag.Errorf("error updating tags of SDRS protected instance %s: %s", d.Id(), err)
		}
	}
	return resourceProtectedInstanceRead(ctx, d, meta)
}

func resourceProtectedInstanceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	deleteOpts := protectedinstances.DeleteOpts{
		DeleteTargetServer: d.Get("delete_target_server").(bool),
		DeleteTargetEip:    d.Get("delete_target_eip").(bool),
	}
	n, err := protectedinstances.Delete(client, d.Id(), deleteOpts).ExtractJobResponse()
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error.code", "SDRS.1301"),
			"error deleting SDRS protected instance",
		)
	}

	deleteTimeoutSec := int(d.Timeout(schema.TimeoutDelete).Seconds())
	if err := protectedinstances.WaitForJobSuccess(client, deleteTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
