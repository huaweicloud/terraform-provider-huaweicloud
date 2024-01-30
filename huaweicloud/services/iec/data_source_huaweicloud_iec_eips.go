package iec

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/publicips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
)

// @API IEC GET /v1/publicips
func DataSourceEips() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEipsRead,

		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"site_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
					},
				},
			},
		},
	}
}

func dataSourceEipsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	eipClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	listpts := publicips.ListOpts{
		SiteID: d.Get("site_id").(string),
		PortID: d.Get("port_id").(string),
	}

	allEips, err := publicips.List(eipClient, listpts).Extract()
	if err != nil {
		return diag.Errorf("unable to extract IEC public IPs: %s", err)
	}

	total := len(allEips.PublicIPs)
	if total < 1 {
		return diag.Errorf("your query returned no results, please change your search criteria and try again")
	}

	log.Printf("[INFO] Retrieved [%d] IEC public IPs using given filter", total)
	firstEip := allEips.PublicIPs[0]
	d.SetId(firstEip.ID)

	iecEips := make([]map[string]interface{}, 0, total)
	for _, item := range allEips.PublicIPs {
		val := map[string]interface{}{
			"id":                   item.ID,
			"public_ip":            item.PublicIpAddress,
			"private_ip":           item.PrivateIpAddress,
			"port_id":              item.PortID,
			"status":               eip.NormalizeEipStatus(item.Status),
			"ip_version":           item.IPVersion,
			"bandwidth_id":         item.BandwidthID,
			"bandwidth_name":       item.BandwidthName,
			"bandwidth_size":       item.BandwidthSize,
			"bandwidth_share_type": item.BandwidthShareType,
		}

		iecEips = append(iecEips, val)
	}

	mErr := multierror.Append(nil,
		d.Set("site_info", firstEip.SiteInfo),
		d.Set("eips", iecEips),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving IEC public IPs: %s", err)
	}

	return nil
}
