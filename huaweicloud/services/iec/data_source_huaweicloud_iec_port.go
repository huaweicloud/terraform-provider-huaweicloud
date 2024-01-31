package iec

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/iec/v1/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC GET /v1/ports
func DataSourcePort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePortRead,

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
			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"mac_address": {
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
				Computed: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourcePortRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	listOpts := ports.ListOpts{
		ID:         d.Get("id").(string),
		NetworkID:  d.Get("subnet_id").(string),
		MacAddress: d.Get("mac_address").(string),
	}

	var ipFilter bool
	if v, ipFilter := d.GetOk("fixed_ip"); ipFilter {
		listOpts.FixedIPs = []string{fmt.Sprintf("ip_address=%s", v)}
	}

	allPorts, err := ports.List(iecClient, listOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to retrieve IEC port: %s", err)
	}

	total := len(allPorts.Ports)
	if total < 1 {
		return diag.Errorf("your query returned no results. Please change your search criteria and try again")
	}
	if total > 1 {
		return diag.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	port := allPorts.Ports[0]
	log.Printf("[DEBUG] Retrieved IEC port %s: %+v", port.ID, port)
	d.SetId(port.ID)
	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("mac_address", port.MacAddress),
		d.Set("subnet_id", port.NetworkID),
		d.Set("status", port.Status),
		d.Set("site_id", port.SiteID),
		d.Set("security_groups", port.SecurityGroups),
	)

	if !ipFilter && len(port.FixedIPs) > 0 {
		fixedIP := port.FixedIPs[0].IpAddress
		mErr = multierror.Append(mErr, d.Set("fixed_ip", fixedIP))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}
