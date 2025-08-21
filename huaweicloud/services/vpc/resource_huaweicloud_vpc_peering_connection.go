package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/peerings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC DELETE /v2.0/vpc/peerings/{id}
// @API VPC GET /v2.0/vpc/peerings/{id}
// @API VPC PUT /v2.0/vpc/peerings/{id}
// @API VPC POST /v2.0/vpc/peerings
func ResourceVpcPeeringConnectionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCPeeringCreate,
		ReadContext:   resourceVPCPeeringRead,
		UpdateContext: resourceVPCPeeringUpdate,
		DeleteContext: resourceVPCPeeringDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateString64WithChinese,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVPCPeeringCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	peeringClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection client: %s", err)
	}

	requestvpcinfo := peerings.VpcInfo{
		VpcId: d.Get("vpc_id").(string),
	}

	acceptvpcinfo := peerings.VpcInfo{
		VpcId:    d.Get("peer_vpc_id").(string),
		TenantId: d.Get("peer_tenant_id").(string),
	}

	createOpts := peerings.CreateOpts{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		RequestVpcInfo: requestvpcinfo,
		AcceptVpcInfo:  acceptvpcinfo,
	}

	n, err := peerings.Create(peeringClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] Waiting for VPC Peering Connection(%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"PENDING_ACCEPTANCE", "ACTIVE"},
		Refresh:    waitForVpcPeeringActive(peeringClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection: %s", err)
	}

	return resourceVPCPeeringRead(ctx, d, meta)
}

func resourceVPCPeeringRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	peeringClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection client: %s", err)
	}

	n, err := peerings.Get(peeringClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPC Peering Connection")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("status", n.Status),
		d.Set("description", n.Description),
		d.Set("vpc_id", n.RequestVpcInfo.VpcId),
		d.Set("peer_vpc_id", n.AcceptVpcInfo.VpcId),
		d.Set("peer_tenant_id", n.AcceptVpcInfo.TenantId),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting VPC Peering Connection fields: %s", err)
	}

	return nil
}

func resourceVPCPeeringUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	peeringClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection client: %s", err)
	}

	updateOpts := peerings.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: utils.String(d.Get("description").(string)),
	}

	_, err = peerings.Update(peeringClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating VPC Peering Connection: %s", err)
	}

	return resourceVPCPeeringRead(ctx, d, meta)
}

func resourceVPCPeeringDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	peeringClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcPeeringDelete(peeringClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting VPC Peering Connection: %s", err)
	}

	return nil
}

func waitForVpcPeeringActive(peeringClient *golangsdk.ServiceClient, peeringId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := peerings.Get(peeringClient, peeringId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "PENDING_ACCEPTANCE" || n.Status == "ACTIVE" {
			return n, n.Status, nil
		}

		return n, "CREATING", nil
	}
}

func waitForVpcPeeringDelete(peeringClient *golangsdk.ServiceClient, peeringId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := peerings.Get(peeringClient, peeringId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted vpc peering connection %s", peeringId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = peerings.Delete(peeringClient, peeringId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted vpc peering connection %s", peeringId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault409); ok {
				return r, "ACTIVE", nil
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
