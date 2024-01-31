package iec

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/bandwidths"
	ieccommon "github.com/chnsz/golangsdk/openstack/iec/v1/common"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API IEC GET /v1/bandwidths
func DataSourceBandWidths() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBandWidthsRead,

		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"site_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidths": {
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
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceBandWidthsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	listOpts := bandwidths.ListOpts{
		SiteID: d.Get("site_id").(string),
	}

	allBWs, err := bandwidths.List(iecClient, listOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to extract IEC bandwidths: %s", err)
	}

	total := len(allBWs.BandWidth)
	log.Printf("[INFO] Retrieved [%d] IEC bandwidths using given filter", total)

	ids := make([]string, 0, total)
	iecBWs := make([]map[string]interface{}, total)
	for i, item := range allBWs.BandWidth {
		ids = append(ids, item.ID)
		iecBWs[i] = map[string]interface{}{
			"id":          item.ID,
			"name":        item.Name,
			"size":        item.Size,
			"share_type":  item.ShareType,
			"charge_mode": item.ChargeMode,
			"status":      item.Status,
			"line":        getLineName(item.Operator),
		}
	}

	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("bandwidths", iecBWs),
	)

	if total > 0 {
		mErr = multierror.Append(mErr, d.Set("site_info", allBWs.BandWidth[0].SiteInfo))
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving IEC bandwidths: %s", err)
	}

	return nil
}

func getLineName(operator ieccommon.Operator) string {
	if operator.Name != "" {
		return operator.Name
	}

	if operator.I18nName != "" {
		return operator.I18nName
	}

	if operator.Sa != "" {
		return operator.Sa
	}

	return ""
}
