package dcs

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/dcs/v2/maintainwindows"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/instances/maintain-windows
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dcsV2Client, err := cfg.DcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
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
		return diag.Errorf("error while filtering data : %s", err)
	}

	if len(filteredMVs) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	mw := filteredMVs[0].(maintainwindows.MaintainWindow)
	d.SetId(strconv.Itoa(mw.ID))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("begin", mw.Begin),
		d.Set("end", mw.End),
		d.Set("default", mw.Default),
	)
	log.Printf("[DEBUG] Dcs MaintainWindow : %+v", mw)

	return diag.FromErr(mErr.ErrorOrNil())
}
