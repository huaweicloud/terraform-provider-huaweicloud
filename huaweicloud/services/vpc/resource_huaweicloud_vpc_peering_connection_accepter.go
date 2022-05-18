package vpc

import (
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/peerings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceVpcPeeringConnectionAccepterV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCPeeringAccepterV2Create,
		Read:   resourceVpcPeeringAccepterRead,
		Update: resourceVPCPeeringAccepterUpdate,
		Delete: resourceVPCPeeringAccepterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceVPCPeeringAccepterV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	peeringClient, err := config.NetworkingV2Client(config.GetRegion(d))

	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud Peering client: %s", err)
	}

	id := d.Get("vpc_peering_connection_id").(string)

	n, err := peerings.Get(peeringClient, id).Extract()
	if err != nil {
		return fmtp.Errorf("Error retrieving Huaweicloud Vpc Peering Connection: %s", err)
	}

	if n.Status != "PENDING_ACCEPTANCE" {
		return fmtp.Errorf("VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.")
	}

	var expectedStatus string

	if _, ok := d.GetOk("accept"); ok {

		expectedStatus = "ACTIVE"
		_, err := peerings.Accept(peeringClient, id).ExtractResult()

		if err != nil {
			return fmtp.Errorf("Unable to accept VPC Peering Connection: {{err}}", err)
		}

	} else {
		expectedStatus = "REJECTED"

		_, err := peerings.Reject(peeringClient, id).ExtractResult()

		if err != nil {
			return fmtp.Errorf("Unable to reject VPC Peering Connection: {{err}}", err)
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

	_, err = stateConf.WaitForState()
	d.SetId(n.ID)
	logp.Printf("[INFO] VPC Peering Connection status: %s", expectedStatus)

	return resourceVpcPeeringAccepterRead(d, meta)

}

func resourceVpcPeeringAccepterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	peeringclient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud peering client: %s", err)
	}

	n, err := peerings.Get(peeringclient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving Huaweicloud Vpc Peering Connection: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("status", n.Status)
	d.Set("vpc_id", n.RequestVpcInfo.VpcId)
	d.Set("peer_vpc_id", n.AcceptVpcInfo.VpcId)
	d.Set("peer_tenant_id", n.AcceptVpcInfo.TenantId)
	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceVPCPeeringAccepterUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("accept") {
		return fmtp.Errorf("VPC peering action not permitted: Can not accept/reject peering request not in pending_acceptance state.'")
	}

	return resourceVpcPeeringAccepterRead(d, meta)
}

func resourceVPCPeeringAccepterDelete(d *schema.ResourceData, meta interface{}) error {
	logp.Printf("[WARN] Will not delete VPC peering connection. Terraform will remove this resource from the state file, however resources may remain.")
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
