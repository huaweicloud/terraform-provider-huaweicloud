package dms

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/dms/v2/maintainwindows"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DMS GET /v2/instances/maintain-windows
func DataSourceDmsMaintainWindow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsMaintainWindowRead,

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

func dataSourceDmsMaintainWindowRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dmsV2Client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client V2: %s", err)
	}

	maintainWindow, err := maintainwindows.Get(dmsV2Client)
	if err != nil {
		return diag.FromErr(err)
	}

	filter := make(map[string]interface{})
	if v, ok := d.GetOk("seq"); ok {
		filter["ID"] = v.(int)
	}
	if v, ok := d.GetOk("begin"); ok {
		filter["Begin"] = v.(string)
	}
	if v, ok := d.GetOk("end"); ok {
		filter["End"] = v.(string)
	}
	if v, ok := d.GetOk("default"); ok {
		filter["Default"] = v.(bool)
	}

	data, err := utils.FilterSliceWithZeroField(maintainWindow, filter)
	if err != nil {
		return diag.Errorf("error filtering DMS maintain window data, %s", err)
	}
	if len(data) < 1 {
		return diag.Errorf("your query returned no results. Please change your filters and try again.")
	}

	mw := data[0].(maintainwindows.MaintainWindow)
	log.Printf("[DEBUG] Dms MaintainWindow : %#v", mw)

	d.SetId(strconv.Itoa(mw.ID))
	mErr := multierror.Append(
		d.Set("seq", mw.ID),
		d.Set("begin", mw.Begin),
		d.Set("end", mw.End),
		d.Set("default", mw.Default),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting DMS maintain window attributes: %s", mErr)
	}

	return nil
}
