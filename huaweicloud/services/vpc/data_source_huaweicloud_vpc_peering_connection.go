package vpc

import (
	"github.com/chnsz/golangsdk/openstack/networking/v2/peerings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceVpcPeeringConnectionV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcPeeringConnectionV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateString64WithChinese,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"peer_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcPeeringConnectionV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	peeringClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	listOpts := peerings.ListOpts{
		ID:         d.Get("id").(string),
		Name:       d.Get("name").(string),
		Status:     d.Get("status").(string),
		TenantId:   d.Get("peer_tenant_id").(string),
		VpcId:      d.Get("vpc_id").(string),
		Peer_VpcId: d.Get("peer_vpc_id").(string),
	}

	refinedPeering, err := peerings.List(peeringClient, listOpts)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve vpc peering connections: %s", err)
	}

	if len(refinedPeering) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedPeering) > 1 {
		return fmtp.Errorf("Multiple VPC peering connections matched." +
			" Use additional constraints to reduce matches to a single VPC peering connection")
	}

	Peering := refinedPeering[0]

	logp.Printf("[INFO] Retrieved Vpc peering Connections using given filter %s: %+v", Peering.ID, Peering)
	d.SetId(Peering.ID)

	d.Set("id", Peering.ID)
	d.Set("name", Peering.Name)
	d.Set("status", Peering.Status)
	d.Set("vpc_id", Peering.RequestVpcInfo.VpcId)
	d.Set("peer_vpc_id", Peering.AcceptVpcInfo.VpcId)
	d.Set("peer_tenant_id", Peering.AcceptVpcInfo.TenantId)
	d.Set("region", config.GetRegion(d))

	return nil
}
