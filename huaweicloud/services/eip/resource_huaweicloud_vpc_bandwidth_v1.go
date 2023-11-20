package eip

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ResourceVpcBandWidthV1
// Add resource bandwidth with update function calls v1 API to support provider which only published v1 API to update the bandwidth.
func ResourceVpcBandWidthV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcBandWidthV2Create,
		ReadContext:   resourceVpcBandWidthV2Read,
		UpdateContext: resourceVpcBandWidthV1Update,
		DeleteContext: resourceVpcBandWidthV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicips": publicIPListComputedSchema(),
		},
	}
}

func resourceVpcBandWidthV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	networkingClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	bssClient, err := cfg.BssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS V2 client: %s", err)
	}

	bwID := d.Id()
	if d.HasChanges("name", "size") {
		updateOpts := bandwidths.UpdateOpts{
			Bandwidth: bandwidths.Bandwidth{
				Name: d.Get("name").(string),
				Size: d.Get("size").(int),
			},
		}

		// ExtendParam is valid and mandatory when changing size field in pre-paid billing mode
		if d.HasChange("size") && d.Get("charging_mode").(string) == "prePaid" {
			updateOpts.ExtendParam = &bandwidths.ExtendParam{
				IsAutoPay: "true",
			}
		}

		log.Printf("[DEBUG] bandwidth update options: %#v", updateOpts)
		resp, err := bandwidths.Update(networkingClient, bwID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating bandwidth (%s): %s", bwID, err)
		}

		if resp.OrderID != "" {
			if err := common.WaitOrderComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("auto_renew") {
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), bwID); err != nil {
			return diag.Errorf("error updating the auto-renew of the bandwidth (%s): %s", bwID, err)
		}
	}

	return resourceVpcBandWidthV2Read(ctx, d, meta)
}
