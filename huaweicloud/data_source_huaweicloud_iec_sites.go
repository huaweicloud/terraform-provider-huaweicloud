package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	iec_common "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/sites"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func dataSourceIecSites() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIecSitesV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"area": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"province": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"city": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sites": {
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
						"area": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"province": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lines": {
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
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
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

func dataSourceIecSitesV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	listOpts := sites.ListSiteOpts{
		Area:     d.Get("area").(string),
		Province: d.Get("province").(string),
		City:     d.Get("city").(string),
	}
	pages, err := sites.List(iecClient, listOpts).AllPages()
	if err != nil {
		return fmtp.Errorf("Unable to retrieve iec sites: %s", err)
	}

	allSites, err := sites.ExtractSites(pages)
	if err != nil {
		return fmtp.Errorf("Unable to extract iec sites: %s", err)
	}
	total := len(allSites.Sites)
	if total < 1 {
		return fmtp.Errorf("Your query returned no results of huaweicloud_iec_sites. " +
			"Please change your search criteria and try again.")
	}

	logp.Printf("[INFO] Retrieved [%d] IEC sites using given filter", total)
	iecSites := make([]map[string]interface{}, 0, total)
	for _, item := range allSites.Sites {
		val := map[string]interface{}{
			"id":       item.ID,
			"name":     item.Name,
			"area":     item.Area,
			"province": item.Province,
			"city":     item.City,
			"status":   item.Status,
			"lines":    flattenSiteLines(&item),
		}
		iecSites = append(iecSites, val)
	}
	if err := d.Set("sites", iecSites); err != nil {
		return fmtp.Errorf("Error saving IEC sites: %s", err)
	}

	site := allSites.Sites[0]
	d.SetId(site.ID)

	return nil
}

func flattenSiteLines(site *iec_common.Site) []map[string]interface{} {
	siteLines := make([]map[string]interface{}, len(site.EipPools))
	for i, item := range site.EipPools {
		siteLines[i] = map[string]interface{}{
			"id":         item.PoolID,
			"name":       item.DisplayName,
			"operator":   item.OperatorID.Name,
			"ip_version": item.IPVersion,
		}
	}

	return siteLines
}
