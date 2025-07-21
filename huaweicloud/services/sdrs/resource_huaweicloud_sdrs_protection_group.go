package sdrs

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectiongroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SDRS POST /v1/{project_id}/server-groups/{id}/action
// @API SDRS DELETE /v1/{project_id}/server-groups/{id}
// @API SDRS GET /v1/{project_id}/server-groups/{id}
// @API SDRS PUT /v1/{project_id}/server-groups/{id}
// @API SDRS POST /v1/{project_id}/server-groups
func ResourceProtectionGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtectionGroupCreate,
		ReadContext:   resourceProtectionGroupRead,
		UpdateContext: resourceProtectionGroupUpdate,
		DeleteContext: resourceProtectionGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of a protection group.`,
			},
			"source_availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the production site AZ of a protection group.`,
			},
			"target_availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the disaster recovery site AZ of a protection group.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of an active-active domain.`,
			},
			"source_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the VPC for the production site.`,
			},
			"dr_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: `Specifies the deployment model. The default value is **migration**, indicating migration
within a VPC.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the description of a protection group.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether enable the protection group start protecting`,
			},
		},
	}
}

func resourceProtectionGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	createOpts := protectiongroups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		SourceAZ:    d.Get("source_availability_zone").(string),
		TargetAZ:    d.Get("target_availability_zone").(string),
		DomainID:    d.Get("domain_id").(string),
		SourceVpcID: d.Get("source_vpc_id").(string),
		DrType:      d.Get("dr_type").(string),
	}
	n, err := protectiongroups.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating SDRS protection group: %s", err)
	}

	createTimeoutSec := int(d.Timeout(schema.TimeoutCreate).Seconds())
	if err = protectiongroups.WaitForJobSuccess(client, createTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}

	groupID, err := protectiongroups.GetJobEntity(client, n.JobID, "server_group_id")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(groupID.(string))

	// enable the protection group start protecting.
	// It can only be set to true when there's replication pairs within the protection group.
	if d.Get("enable").(bool) {
		if err := enableAndWaitForProtectionGroupSuccess(d, client); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceProtectionGroupRead(ctx, d, meta)
}

func resourceProtectionGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	n, err := protectiongroups.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error.code", "SDRS.1013"),
			"error retrieving SDRS protection group")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("description", n.Description),
		d.Set("source_availability_zone", n.SourceAZ),
		d.Set("target_availability_zone", n.TargetAZ),
		d.Set("domain_id", n.DomainID),
		d.Set("source_vpc_id", n.SourceVpcID),
		d.Set("dr_type", n.DrType),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceProtectionGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	if d.HasChange("name") {
		updateOpts := protectiongroups.UpdateOpts{
			Name: d.Get("name").(string),
		}
		if _, err = protectiongroups.Update(client, d.Id(), updateOpts).Extract(); err != nil {
			return diag.Errorf("error updating SDRS protection group name, %s", err)
		}
	}

	if d.HasChange("enable") {
		var enableErr error
		if d.Get("enable").(bool) {
			enableErr = enableAndWaitForProtectionGroupSuccess(d, client)
		} else {
			enableErr = disableAndWaitForProtectionGroupSuccess(d, client)
		}
		if enableErr != nil {
			return diag.FromErr(enableErr)
		}
	}
	return resourceProtectionGroupRead(ctx, d, meta)
}

func resourceProtectionGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}
	// use for WaitForJobSuccess
	client.Endpoint = client.ResourceBase

	n, err := protectiongroups.Delete(client, d.Id()).ExtractJobResponse()
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error.code", "SDRS.1013"),
			"error deleting SDRS protection group")
	}

	deleteTimeoutSec := int(d.Timeout(schema.TimeoutDelete).Seconds())
	if err = protectiongroups.WaitForJobSuccess(client, deleteTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func enableAndWaitForProtectionGroupSuccess(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	n, err := protectiongroups.Enable(client, d.Id()).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("error start SDRS protection group: %s", err)
	}

	updateTimeoutSec := int(d.Timeout(schema.TimeoutUpdate).Seconds())
	return protectiongroups.WaitForJobSuccess(client, updateTimeoutSec, n.JobID)
}

func disableAndWaitForProtectionGroupSuccess(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	n, err := protectiongroups.Disable(client, d.Id()).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("error stop SDRS protection group: %s", err)
	}

	updateTimeoutSec := int(d.Timeout(schema.TimeoutUpdate).Seconds())
	return protectiongroups.WaitForJobSuccess(client, updateTimeoutSec, n.JobID)
}
