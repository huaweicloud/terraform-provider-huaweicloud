package deprecated

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/dms/v2/availablezones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceDmsAZ() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceDmsAZRead,
		DeprecationMessage: "Deprecated. Please use \"huaweicloud_availability_zones\" instead.",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceDmsAZRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DMS key client V2: %s", err)
	}

	v, err := availablezones.Get(dmsV2Client)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Dms az : %+v", v)
	var filteredAZs []availablezones.AvailableZone
	if v.RegionID == config.GetRegion(d) {
		AZs := v.AvailableZones
		for _, newAZ := range AZs {
			if newAZ.ResourceAvailability != "true" || newAZ.SoldOut {
				continue
			}

			name := d.Get("name").(string)
			if name != "" && newAZ.Name != name {
				continue
			}

			code := d.Get("code").(string)
			if code != "" && newAZ.Code != code {
				continue
			}

			port := d.Get("port").(string)
			if port != "" && newAZ.Port != port {
				continue
			}
			filteredAZs = append(filteredAZs, newAZ)
		}
	}

	if len(filteredAZs) < 1 {
		return fmtp.DiagErrorf("Not found any available zones")
	}

	az := filteredAZs[0]
	logp.Printf("[DEBUG] Dms az : %+v", az)

	d.SetId(az.ID)
	mErr := multierror.Append(
		d.Set("code", az.Code),
		d.Set("name", az.Name),
		d.Set("port", az.Port),
		d.Set("ipv6_enable", az.Ipv6Enable),
	)
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("Error setting DMS AZ attributes: %s", mErr)
	}
	return nil
}
