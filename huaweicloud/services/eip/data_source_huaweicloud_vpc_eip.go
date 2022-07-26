package eip

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceVpcEip() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcEipRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
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

func dataSourceVpcEipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating VPC client: %s", err)
	}

	var listOpts eips.ListOpts
	if portId, ok := d.GetOk("port_id"); ok {
		listOpts.PortId = []string{portId.(string)}
	}

	if publicIp, ok := d.GetOk("public_ip"); ok {
		listOpts.PublicIp = []string{publicIp.(string)}
	}

	listOpts.EnterpriseProjectId = config.DataGetEnterpriseProjectID(d)

	pages, err := eips.List(vpcClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve eips: %s ", err)
	}

	if len(allEips) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allEips) > 1 {
		return fmtp.DiagErrorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Eip := allEips[0]

	d.SetId(Eip.ID)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("status", NormalizeEIPStatus(Eip.Status)),
		d.Set("public_ip", Eip.PublicAddress),
		d.Set("ipv6_address", Eip.PublicIpv6Address),
		d.Set("ip_version", Eip.IpVersion),
		d.Set("port_id", Eip.PortID),
		d.Set("type", Eip.Type),
		d.Set("private_ip", Eip.PrivateAddress),
		d.Set("bandwidth_id", Eip.BandwidthID),
		d.Set("bandwidth_size", Eip.BandwidthSize),
		d.Set("bandwidth_share_type", Eip.BandwidthShareType),
		d.Set("enterprise_project_id", Eip.EnterpriseProjectID),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting eip fields: %s", err)
	}

	return nil
}
