package eip

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	bandwidthsv1 "github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	DefaultBandWidthChargeMode = "bandwidth"
	DefaultBandWidthSize       = 5
)

// @API EIP POST /v1/{project_id}/bandwidths/{ID}/insert
// @API EIP POST /v1/{project_id}/bandwidths/{ID}/remove
// @API EIP GET /v1/{project_id}/bandwidths/{id}
// @API EIP GET /v1/{project_id}/publicips/{id}
func ResourceBandWidthAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBandWidthAssociateCreate,
		ReadContext:   resourceBandWidthAssociateRead,
		UpdateContext: resourceBandWidthAssociateUpdate,
		DeleteContext: resourceBandWidthAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceBandWidthAssociationImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth_charge_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The charge mode after removal bandwidth.`,
			},
			"bandwidth_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The size (Mbits/s) after removal bandwidth.`,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBandWidthAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bwClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating bandwidth client: %s", err)
	}
	eipClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EIP client: %s", err)
	}

	bwID := d.Get("bandwidth_id").(string)
	eipID := d.Get("eip_id").(string)

	err = insertEIPsToBandwidth(d, bwClient, eipClient, bwID, eipID)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceID := fmt.Sprintf("%s/%s", bwID, eipID)
	d.SetId(resourceID)
	return resourceBandWidthAssociateRead(ctx, d, meta)
}

func resourceBandWidthAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	eipClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating bandwidth v1 client: %s", err)
	}

	bwID := d.Get("bandwidth_id").(string)
	eipID := d.Get("eip_id").(string)

	b, err := bandwidthsv1.Get(eipClient, bwID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "bandwidth associate")
	}

	var eipItem *bandwidthsv1.PublicIpinfo
	allEIPs := b.PublicipInfo
	for i := range allEIPs {
		if allEIPs[i].PublicipId == eipID {
			eipItem = &allEIPs[i]
			break
		}
	}

	if eipItem == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("bandwidth_id", b.ID),
		d.Set("bandwidth_name", b.Name),
		d.Set("eip_id", eipItem.PublicipId),
		d.Set("public_ip", eipItem.PublicipAddress),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBandWidthAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bwClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating bandwidth client: %s", err)
	}
	eipClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EIP client: %s", err)
	}

	if d.HasChange("eip_id") {
		bwID := d.Get("bandwidth_id").(string)
		oldVal, newVal := d.GetChange("eip_id")

		err = removeEIPsFromBandwidth(d, bwClient, eipClient, bwID, oldVal.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = insertEIPsToBandwidth(d, bwClient, eipClient, bwID, newVal.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		// update the resource ID
		resourceID := fmt.Sprintf("%s/%s", bwID, newVal.(string))
		d.SetId(resourceID)
	}

	return resourceBandWidthAssociateRead(ctx, d, meta)
}

func resourceBandWidthAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bwClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating bandwidth v2.0 client: %s", err)
	}
	eipClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EIP client: %s", err)
	}

	bwID := d.Get("bandwidth_id").(string)
	eipID := d.Get("eip_id").(string)

	err = removeEIPsFromBandwidth(d, bwClient, eipClient, bwID, eipID)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func insertEIPsToBandwidth(d *schema.ResourceData, bwClient, eipClient *golangsdk.ServiceClient, bwID, eipID string) error {
	publicIp, err := eips.Get(eipClient, eipID).Extract()
	if err != nil {
		return fmt.Errorf("error fetching EIP %s: %s", eipID, err)
	}

	if publicIp.BandwidthShareType == "WHOLE" {
		// the EIP is already associated to the shared bandwidth, return with success.
		if publicIp.BandwidthID == bwID {
			log.Printf("[DEBUG] EIP %s is already associated to shared bandwidth %s", eipID, bwID)
			return nil
		}

		// the EIP is associated to another shared bandwidth, we should remove it first.
		if err := removeEIPsFromBandwidth(d, bwClient, eipClient, publicIp.BandwidthID, eipID); err != nil {
			return err
		}
	}

	publicipOpts := []bandwidths.PublicIpInfoID{
		{
			PublicIPID: eipID,
		},
	}
	insertOpts := bandwidths.BandWidthInsertOpts{
		PublicipInfo: publicipOpts,
	}

	log.Printf("[DEBUG] Add EIP %s to bandwidth %s", eipID, bwID)
	if _, err := bandwidths.Insert(bwClient, bwID, insertOpts).Extract(); err != nil {
		return fmt.Errorf("error inserting EIP %s into bandwidth %s: %s", eipID, bwID, err)
	}
	return nil
}

func removeEIPsFromBandwidth(d *schema.ResourceData, bwClient, eipClient *golangsdk.ServiceClient, bwID, eipID string) error {
	if _, err := eips.Get(eipClient, eipID).Extract(); err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] unnecessary to remove the non-existent EIP %s from bandwidth %s", eipID, bwID)
			return nil
		}
		return fmt.Errorf("error fetching EIP %s: %s", eipID, err)
	}

	bwChargeMode := d.Get("bandwidth_charge_mode").(string)
	if bwChargeMode == "" {
		bwChargeMode = DefaultBandWidthChargeMode
	}

	bwSize := d.Get("bandwidth_size").(int)
	if bwSize == 0 {
		bwSize = DefaultBandWidthSize
	}

	removeOpts := bandwidths.BandWidthRemoveOpts{
		ChargeMode: bwChargeMode,
		Size:       &bwSize,
		PublicipInfo: []bandwidths.PublicIpInfoID{
			{
				PublicIPID: eipID,
			},
		},
	}

	log.Printf("[DEBUG] Remove EIP %s from bandwidth %s", eipID, bwID)
	err := bandwidths.Remove(bwClient, bwID, removeOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("error removing EIP %s from bandwidth %s: %s", eipID, bwID, err)
	}

	return nil
}

func resourceBandWidthAssociationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<bandwidth_id>/<eip_id>', but '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("bandwidth_id", parts[0]),
		d.Set("eip_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
