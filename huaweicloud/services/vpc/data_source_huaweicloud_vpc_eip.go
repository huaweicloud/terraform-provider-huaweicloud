package vpc

import (
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceVpcEip() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcEipRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth_share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcEipRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating networking client: %s", err)
	}

	listOpts := &eips.ListOpts{
		PortId:   []string{d.Get("port_id").(string)},
		PublicIp: []string{d.Get("public_ip").(string)},
	}

	epsID := config.GetEnterpriseProjectID(d)
	if epsID != "" {
		listOpts.EnterpriseProjectId = epsID
	}

	pages, err := eips.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve eips: %s ", err)
	}

	if len(allEips) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allEips) > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Eip := allEips[0]

	d.SetId(Eip.ID)
	d.Set("region", config.GetRegion(d))
	d.Set("status", NormalizeEIPStatus(Eip.Status))
	d.Set("public_ip", Eip.PublicAddress)
	d.Set("port_id", Eip.PortID)
	d.Set("type", Eip.Type)
	d.Set("private_ip", Eip.PrivateAddress)
	d.Set("bandwidth_id", Eip.BandwidthID)
	d.Set("bandwidth_size", Eip.BandwidthSize)
	d.Set("bandwidth_share_type", Eip.BandwidthShareType)
	d.Set("enterprise_project_id", Eip.EnterpriseProjectID)

	return nil
}
