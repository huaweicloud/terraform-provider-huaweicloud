package dms

import (
	"context"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/dms/v2/maintainwindows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

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
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DMS client V2: %s", err)
	}

	maintainWindows, err := maintainwindows.Get(dmsV2Client)
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

	data, err := utils.FilterSliceWithZeroField(maintainWindows, filter)
	if err != nil {
		return fmtp.DiagErrorf("Error filtering DMS maintain window data, %s", err)
	}
	if len(data) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. Please change your filters and try again.")
	}

	mw := data[0].(maintainwindows.MaintainWindow)
	logp.Printf("[DEBUG] Dms MaintainWindow : %#v", mw)

	d.SetId(strconv.Itoa(mw.ID))
	mErr := multierror.Append(
		d.Set("seq", mw.ID),
		d.Set("begin", mw.Begin),
		d.Set("end", mw.End),
		d.Set("default", mw.Default),
	)
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting DMS maintain window attributes: %s", mErr)
	}

	return nil
}
