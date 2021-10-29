package dcs

import (
	"context"
	"strconv"

	"github.com/chnsz/golangsdk/openstack/dcs/v2/maintainwindows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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

func dataSourceDcsMaintainWindowRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dcsV2Client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating DCS key client: %s", err)
	}

	v, err := maintainwindows.Get(dcsV2Client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	filteredMVs, err := utils.FilterSliceWithField(v.MaintainWindows, map[string]interface{}{
		"ID":      d.Get("seq").(int),
		"Begin":   d.Get("begin").(string),
		"End":     d.Get("end").(string),
		"Default": d.Get("default").(bool),
	})
	if err != nil {
		return fmtp.DiagErrorf("Error while filtering data : %s", err)
	}

	if len(filteredMVs) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	mw := filteredMVs[0].(maintainwindows.MaintainWindow)
	d.SetId(strconv.Itoa(mw.ID))
	d.Set("begin", mw.Begin)
	d.Set("end", mw.End)
	d.Set("default", mw.Default)
	logp.Printf("[DEBUG] Dcs MaintainWindow : %+v", mw)

	return nil
}
