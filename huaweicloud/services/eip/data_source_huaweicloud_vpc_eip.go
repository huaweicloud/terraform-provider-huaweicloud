package eip

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API EIP GET /v1/{project_id}/publicips
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
			"name": {
				Type:     schema.TypeString,
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
			"bandwidth_name": {
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcEipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	var listOpts eips.ListOpts
	if portId, ok := d.GetOk("port_id"); ok {
		listOpts.PortId = []string{portId.(string)}
	}

	if publicIp, ok := d.GetOk("public_ip"); ok {
		listOpts.PublicIp = []string{publicIp.(string)}
	}

	listOpts.EnterpriseProjectId = cfg.GetEnterpriseProjectID(d, "all_granted_eps")

	pages, err := eips.List(vpcClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve eips: %s ", err)
	}

	if len(allEips) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allEips) > 1 {
		return diag.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	eipInfo := allEips[0]

	d.SetId(eipInfo.ID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", eipInfo.Alias),
		d.Set("status", NormalizeEipStatus(eipInfo.Status)),
		d.Set("public_ip", eipInfo.PublicAddress),
		d.Set("ipv6_address", eipInfo.PublicIpv6Address),
		d.Set("ip_version", eipInfo.IpVersion),
		d.Set("port_id", eipInfo.PortID),
		d.Set("type", eipInfo.Type),
		d.Set("private_ip", eipInfo.PrivateAddress),
		d.Set("bandwidth_id", eipInfo.BandwidthID),
		d.Set("bandwidth_name", eipInfo.BandwidthName),
		d.Set("bandwidth_size", eipInfo.BandwidthSize),
		d.Set("bandwidth_share_type", eipInfo.BandwidthShareType),
		d.Set("enterprise_project_id", eipInfo.EnterpriseProjectID),
		d.Set("created_at", eipInfo.CreateTime),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting eip fields: %s", err)
	}

	return nil
}
