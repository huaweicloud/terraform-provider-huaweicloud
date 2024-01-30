package iec

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/vpcs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC GET /v1/vpcs
func DataSourceVpc() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcRead,

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

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iecClient, err := cfg.IECV1Client(region)
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:   d.Get("id").(string),
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Query VPCs using given filter: %+v", listOpts)
	allVpcs, err := vpcs.List(iecClient, listOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to retrieve vpcs: %s", err)
	}

	total := len(allVpcs.Vpcs)
	if total < 1 {
		return diag.Errorf("your query returned no results," +
			"please change your search criteria and try again.")
	}
	if total > 1 {
		return diag.Errorf("your query returned more than one result," +
			" please try a more specific search criteria")
	}

	vpcInfo := allVpcs.Vpcs[0]
	log.Printf("[DEBUG] Retrieved IEC VPC %s: %+v", vpcInfo.ID, vpcInfo)

	d.SetId(vpcInfo.ID)

	mErr := multierror.Append(nil,
		d.Set("name", vpcInfo.Name),
		d.Set("cidr", vpcInfo.Cidr),
		d.Set("mode", vpcInfo.Mode),
		d.Set("subnet_num", vpcInfo.SubnetNum),
		d.Set("region", region),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving IEC VPC: %s", err)
	}

	return nil
}
