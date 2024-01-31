package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v2/peerings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC GET /v2.0/vpc/peerings
func DataSourceVpcPeeringConnectionV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcPeeringConnectionV2Read,

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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcPeeringConnectionV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	peeringClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering Connection client: %s", err)
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
		return diag.Errorf("unable to retrieve VPC Peering Connections: %s", err)
	}

	if len(refinedPeering) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedPeering) > 1 {
		return diag.Errorf("multiple VPC Peering Connections matched." +
			" Use additional constraints to reduce matches to a single VPC Peering Connection")
	}

	item := refinedPeering[0]

	log.Printf("[INFO] Retrieved Vpc peering Connections using given filter %s: %+v", item.ID, item)
	d.SetId(item.ID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", item.Name),
		d.Set("status", item.Status),
		d.Set("description", item.Description),
		d.Set("vpc_id", item.RequestVpcInfo.VpcId),
		d.Set("peer_vpc_id", item.AcceptVpcInfo.VpcId),
		d.Set("peer_tenant_id", item.AcceptVpcInfo.TenantId),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
