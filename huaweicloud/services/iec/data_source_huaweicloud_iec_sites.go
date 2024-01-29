package iec

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ieccommon "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/sites"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC GET /v1/sites
func DataSourceSites() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSitesV1Read,

		Schema: map[string]*schema.Schema{
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

func dataSourceSitesV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	listOpts := sites.ListSiteOpts{
		Area:     d.Get("area").(string),
		Province: d.Get("province").(string),
		City:     d.Get("city").(string),
	}
	pages, err := sites.List(iecClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to retrieve IEC sites: %s", err)
	}

	allSites, err := sites.ExtractSites(pages)
	if err != nil {
		return diag.Errorf("unable to extract IEC sites: %s", err)
	}
	total := len(allSites.Sites)
	if total < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	log.Printf("[INFO] Retrieved [%d] IEC sites using given filter", total)
	iecSites := make([]map[string]interface{}, 0, total)
	for i := range allSites.Sites {
		item := allSites.Sites[i]
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

	mErr := multierror.Append(nil,
		d.Set("sites", iecSites),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving IEC sites: %s", err)
	}

	site := allSites.Sites[0]
	d.SetId(site.ID)

	return nil
}

func flattenSiteLines(site *ieccommon.Site) []map[string]interface{} {
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
