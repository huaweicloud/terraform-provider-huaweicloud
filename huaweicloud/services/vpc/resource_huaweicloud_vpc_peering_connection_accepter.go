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
)

// @API VPC PUT /v2.0/vpc/peerings/{id}/accept
// @API VPC PUT /v2.0/vpc/peerings/{id}/reject
// @API VPC GET /v2.0/vpc/peerings/{id}
func ResourceVpcPeeringConnectionAccepterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCPeeringAccepterCreate,
		ReadContext:   resourceVpcPeeringAccepterRead,
		UpdateContext: resourceVPCPeeringAccepterUpdate,
		DeleteContext: resourceVPCPeeringAccepterDelete,
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
				Computed: true,
				ForceNew: true,
			},
			"vpc_peering_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accept": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVPCPeeringAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	peeringClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection client: %s", err)
	}

	id := d.Get("vpc_peering_connection_id").(string)
	n, err := peerings.Get(peeringClient, id).Extract()
	if err != nil {
		return diag.Errorf("error retrieving Vpc Peering Connection: %s", err)
	}

	if n.Status != "PENDING_ACCEPTANCE" {
		return diag.Errorf("VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.")
	}

	var expectedStatus string

	if _, ok := d.GetOk("accept"); ok {
		expectedStatus = "ACTIVE"

		_, err := peerings.Accept(peeringClient, id).ExtractResult()
		if err != nil {
			return diag.Errorf("unable to accept VPC Peering Connection: %s", err)
		}
	} else {
		expectedStatus = "REJECTED"

		_, err := peerings.Reject(peeringClient, id).ExtractResult()
		if err != nil {
			return diag.Errorf("unable to reject VPC Peering Connection: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{expectedStatus},
		Refresh:    waitForVpcPeeringConnStatus(peeringClient, n.ID, expectedStatus),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the VPC Peering Connection: %s", err)
	}

	d.SetId(n.ID)
	return resourceVpcPeeringAccepterRead(ctx, d, meta)
}

func resourceVpcPeeringAccepterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceVPCPeeringAccepterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("accept") {
		return diag.Errorf("VPC peering action not permitted: Can not accept/reject peering request not in pending_acceptance state.")
	}

	return resourceVpcPeeringAccepterRead(ctx, d, meta)
}

func resourceVPCPeeringAccepterDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("[WARN] Will not delete VPC peering connection. Terraform will remove this resource from the state file, resources may remain.")
	d.SetId("")
	return nil
}

func waitForVpcPeeringConnStatus(peeringClient *golangsdk.ServiceClient, peeringId, expectedStatus string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := peerings.Get(peeringClient, peeringId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == expectedStatus {
			return n, expectedStatus, nil
		}

		return n, "PENDING", nil
	}
}
