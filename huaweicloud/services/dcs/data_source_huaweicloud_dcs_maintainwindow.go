package dcs

import (
	"context"
	"strconv"

	"github.com/chnsz/golangsdk/openstack/dcs/v2/maintainwindows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceDcsMaintainWindow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsMaintainWindowRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"seq": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsMaintainWindowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dcsV2Client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating dcs key client: %s", err)
	}

	v, err := maintainwindows.Get(dcsV2Client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	maintainWindows := v.MaintainWindows
	filteredMVs := make([]maintainwindows.MaintainWindow, 0, len(maintainWindows))
	for _, mv := range maintainWindows {
		seq := d.Get("seq").(int)
		if seq != 0 && mv.ID != seq {
			continue
		}

		begin := d.Get("begin").(string)
		if begin != "" && mv.Begin != begin {
			continue
		}
		end := d.Get("end").(string)
		if end != "" && mv.End != end {
			continue
		}

		df, ok := d.GetOk("default")
		if ok && mv.Default != df.(bool) {
			continue
		}
		filteredMVs = append(filteredMVs, mv)
	}
	if len(filteredMVs) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	mw := filteredMVs[0]
	d.SetId(strconv.Itoa(mw.ID))
	d.Set("begin", mw.Begin)
	d.Set("end", mw.End)
	d.Set("default", mw.Default)
	logp.Printf("[DEBUG] Dcs MaintainWindow : %+v", mw)

	return nil
}
