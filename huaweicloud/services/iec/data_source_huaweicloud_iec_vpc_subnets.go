package iec

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/subnets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API IEC GET /v1/subnets
func DataSourceVpcSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcSubnetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcSubnetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iecClient, err := cfg.IECV1Client(region)
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	vpcID := d.Get("vpc_id").(string)
	listOpts := subnets.ListOpts{
		VpcID:  vpcID,
		SiteID: d.Get("site_id").(string),
	}

	allSubnets, err := subnets.List(iecClient, listOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to retrieve subnets: %s", err)
	}

	total := len(allSubnets.Subnets)
	if total == 0 {
		return diag.Errorf("no matching subnet found for vpc with id %s", vpcID)
	}

	iecSubnets := make([]map[string]interface{}, total)
	allIDs := make([]string, total)
	for i, item := range allSubnets.Subnets {
		val := map[string]interface{}{
			"id":         item.ID,
			"name":       item.Name,
			"cidr":       item.Cidr,
			"gateway_ip": item.GatewayIP,
			"site_id":    item.SiteID,
			"site_info":  item.SiteInfo,
			"dns_list":   item.DNSList,
			"status":     item.Status,
		}
		iecSubnets[i] = val
		allIDs[i] = item.ID
	}

	// set id
	d.SetId(hashcode.Strings(allIDs))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("subnets", iecSubnets),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
