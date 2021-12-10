package huaweicloud

import (
	"context"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/iec/v1/bandwidths"
	iec_common "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func dataSourceIECBandWidths() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIECBandWidthsRead,

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

func dataSourceIECBandWidthsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud IEC client: %s", err)
	}

	listOpts := bandwidths.ListOpts{
		SiteID: d.Get("site_id").(string),
	}

	allBWs, err := bandwidths.List(iecClient, listOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to extract iec bandwidths: %s", err)
	}

	total := len(allBWs.BandWidth)
	logp.Printf("[INFO] Retrieved [%d] IEC bandwidths using given filter", total)

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

	if err := d.Set("bandwidths", iecBWs); err != nil {
		return fmtp.DiagErrorf("Error saving IEC bandwidths: %s", err)
	}
	if total > 0 {
		d.Set("site_info", allBWs.BandWidth[0].SiteInfo)
	}

	return nil
}

func getLineName(operator iec_common.Operator) string {
	if operator.Name != "" {
		return operator.Name
	} else if operator.I18nName != "" {
		return operator.I18nName
	} else if operator.Sa != "" {
		return operator.Sa
	}

	return ""
}
