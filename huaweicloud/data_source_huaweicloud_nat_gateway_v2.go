package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/natgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceNatGatewayV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNatGatewayV2Read,

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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// deprecated
			"router_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use vpc_id instead",
			},
			"internal_network_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use subnet_id instead",
			},
		},
	}
}

func dataSourceNatGatewayV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	natClient, err := config.NatGatewayClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	var vpcID, subnetID string
	if v1, ok := d.GetOk("vpc_id"); ok {
		vpcID = v1.(string)
	} else {
		vpcID = d.Get("router_id").(string)
	}
	if v2, ok := d.GetOk("subnet_id"); ok {
		subnetID = v2.(string)
	} else {
		subnetID = d.Get("internal_network_id").(string)
	}

	listOpts := natgateways.ListOpts{
		ID:                  d.Get("id").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Spec:                d.Get("spec").(string),
		RouterID:            vpcID,
		InternalNetworkID:   subnetID,
		Status:              d.Get("status").(string),
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
	}

	pages, err := natgateways.List(natClient, listOpts).AllPages()
	if err != nil {
		return err
	}

	allNatGateways, err := natgateways.ExtractNatGateways(pages)

	if err != nil {
		return fmtp.Errorf("Unable to retrieve natgateways: %s", err)
	}

	if len(allNatGateways) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allNatGateways) > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	natgateway := allNatGateways[0]

	logp.Printf("[DEBUG] Retrieved Natgateway %s: %+v", natgateway.ID, natgateway)

	d.SetId(natgateway.ID)
	d.Set("name", natgateway.Name)
	d.Set("description", natgateway.Description)
	d.Set("spec", natgateway.Spec)
	d.Set("vpc_id", natgateway.RouterID)
	d.Set("subnet_id", natgateway.InternalNetworkID)
	d.Set("status", natgateway.Status)
	d.Set("enterprise_project_id", natgateway.EnterpriseProjectID)

	return nil
}
