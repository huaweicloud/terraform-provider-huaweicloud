package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/iec/v1/publicips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func dataSourceIECNetworkEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIECNetworkEipsRead,

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

func dataSourceIECNetworkEipsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	eipClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	listpts := publicips.ListOpts{
		SiteID: d.Get("site_id").(string),
		PortID: d.Get("port_id").(string),
	}

	allEips, err := publicips.List(eipClient, listpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to extract iec public ips: %s", err)
	}
	total := len(allEips.PublicIPs)
	if total < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	logp.Printf("[INFO] Retrieved [%d] IEC public IPs using given filter", total)
	firstEip := allEips.PublicIPs[0]
	d.SetId(firstEip.ID)
	d.Set("site_info", firstEip.SiteInfo)

	iecEips := make([]map[string]interface{}, 0, total)
	for _, item := range allEips.PublicIPs {
		val := map[string]interface{}{
			"id":                   item.ID,
			"public_ip":            item.PublicIpAddress,
			"private_ip":           item.PrivateIpAddress,
			"port_id":              item.PortID,
			"status":               eip.NormalizeEIPStatus(item.Status),
			"ip_version":           item.IPVersion,
			"bandwidth_id":         item.BandwidthID,
			"bandwidth_name":       item.BandwidthName,
			"bandwidth_size":       item.BandwidthSize,
			"bandwidth_share_type": item.BandwidthShareType,
		}

		iecEips = append(iecEips, val)
	}
	if err := d.Set("eips", iecEips); err != nil {
		return fmtp.Errorf("Error saving IEC public IPs: %s", err)
	}

	return nil
}
