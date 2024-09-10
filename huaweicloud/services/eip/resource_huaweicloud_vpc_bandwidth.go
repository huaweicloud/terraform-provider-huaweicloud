package eip

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/structs"
	bandwidthsv1 "github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API EIP POST /v2.0/{project_id}/bandwidths/change-to-period
// @API EIP PUT /v2.0/{project_id}/bandwidths/{ID}
// @API EIP PUT /v1/{project_id}/bandwidths/{ID}
// @API EIP DELETE /v2.0/{project_id}/bandwidths/{ID}
// @API EIP GET /v1/{project_id}/bandwidths/{id}
// @API EIP POST /v2.0/{project_id}/bandwidths
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
func ResourceVpcBandWidthV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcBandWidthV2Create,
		ReadContext:   resourceVpcBandWidthV2Read,
		UpdateContext: resourceVpcBandWidthV2Update,
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
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			// charging_mode,  period_unit and period only support changing post-paid to pre-paid billing mode.
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
			"bandwidth_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicips": publicIPListComputedSchema(),
		},
	}
}

func publicIPListComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"ip_version": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"ip_address": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func resourceVpcBandWidthV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	networkingClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	bwClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating bandwidth v1 client: %s", err)
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	size := d.Get("size").(int)
	createOpts := bandwidths.CreateOpts{
		Name:              d.Get("name").(string),
		ChargeMode:        d.Get("charge_mode").(string),
		Size:              &size,
		PublicBorderGroup: d.Get("public_border_group").(string),
		BandwidthType:     d.Get("bandwidth_type").(string),
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID != "" {
		createOpts.EnterpriseProjectId = epsID
	}

	log.Printf("[DEBUG] bandwidth create options: %#v", createOpts)
	b, err := bandwidths.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating bandwidth: %s", err)
	}

	d.SetId(b.ID)

	log.Printf("[DEBUG] Waiting for bandwidth (%s) to become available.", b.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"NORMAL"},
		Pending:    []string{"CREATING"},
		Refresh:    waitForBandwidth(bwClient, b.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for bandwidth (%s) to become ACTIVE: %s",
			b.ID, err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		// we can not create a bandwidth with pre-paid directly due to the API does not support
		// call the change-to-period API as a workaround
		err := changeBandwidthToPeriod(ctx, d, networkingClient, bssClient)
		if err != nil {
			return err
		}
	}

	return resourceVpcBandWidthV2Read(ctx, d, meta)
}

func changeBandwidthToPeriod(ctx context.Context, d *schema.ResourceData, networkingClient,
	bssClient *golangsdk.ServiceClient) diag.Diagnostics {
	changeOpts := bandwidths.ChangeToPeriodOpts{
		BandwidthIDs: []string{d.Id()},
		ExtendParam: structs.ChargeInfo{
			ChargeMode:  "prePaid",
			PeriodType:  d.Get("period_unit").(string),
			PeriodNum:   d.Get("period").(int),
			IsAutoRenew: d.Get("auto_renew").(string),
			IsAutoPay:   "true",
		},
	}
	orderID, err := bandwidths.ChangeToPeriod(networkingClient, changeOpts).Extract()
	if err != nil {
		return diag.Errorf("error changing bandwidth (%s) to pre-paid billing mode: %s", d.Id(), err)
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpcBandWidthV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	networkingV1Client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	networkingClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS V2 client: %s", err)
	}

	bwID := d.Id()
	if d.HasChange("charging_mode") {
		if d.Get("charging_mode").(string) == "postPaid" {
			return diag.Errorf("error updating the charging mode of the bandwidth (%s): %s", bwID,
				"only support changing post-paid bandwidth to pre-paid")
		}
		err := changeBandwidthToPeriod(ctx, d, networkingClient, bssClient)
		if err != nil {
			return err
		}
	} else if d.HasChange("auto_renew") {
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), bwID); err != nil {
			return diag.Errorf("error updating the auto-renew of the bandwidth (%s): %s", bwID, err)
		}
	}

	if d.HasChanges("name", "size", "charge_mode") {
		if d.Get("charging_mode").(string) == "prePaid" {
			updateOpts := bandwidths.UpdateOpts{
				Bandwidth: bandwidths.Bandwidth{
					Name: d.Get("name").(string),
					Size: d.Get("size").(int),
				},
				ExtendParam: &bandwidths.ExtendParam{
					IsAutoPay: "true",
				},
			}

			log.Printf("[DEBUG] bandwidth update options: %#v", updateOpts)
			resp, err := bandwidths.Update(networkingClient, bwID, updateOpts).Extract()
			if err != nil {
				return diag.Errorf("error updating pre-paid bandwidth (%s): %s", bwID, err)
			}

			if resp.OrderID != "" {
				if err := common.WaitOrderComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutUpdate)); err != nil {
					return diag.FromErr(err)
				}
			}
		} else {
			updateOpts := bandwidthsv1.UpdateOpts{
				Name:       d.Get("name").(string),
				Size:       d.Get("size").(int),
				ChargeMode: d.Get("charge_mode").(string),
			}
			log.Printf("[DEBUG] bandwidth update options: %#v", updateOpts)
			_, err := bandwidthsv1.Update(networkingV1Client, bwID, updateOpts).Extract()
			if err != nil {
				return diag.Errorf("error updating post-paid bandwidth (%s): %s", bwID, err)
			}
		}
	}

	return resourceVpcBandWidthV2Read(ctx, d, meta)
}

func resourceVpcBandWidthV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bwClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating bandwidth v1 client: %s", err)
	}

	b, err := bandwidthsv1.Get(bwClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "bandwidth")
	}

	mErr := multierror.Append(
		d.Set("name", b.Name),
		d.Set("size", b.Size),
		d.Set("charge_mode", b.ChargeMode),
		d.Set("enterprise_project_id", b.EnterpriseProjectID),
		d.Set("share_type", b.ShareType),
		d.Set("bandwidth_type", b.BandwidthType),
		d.Set("public_border_group", b.PublicBorderGroup),
		d.Set("status", b.Status),
		d.Set("created_at", b.CreatedAt),
		d.Set("updated_at", b.UpdatedAt),
		d.Set("charging_mode", normalizeChargingMode(b.BillingInfo)),
		d.Set("publicips", flattenPublicIPs(b)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting bandwidth fields: %s", err)
	}

	return nil
}

func resourceVpcBandWidthV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	networkingClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	bwClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating bandwidth v1 client: %s", err)
	}

	bwID := d.Id()
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{bwID}); err != nil {
			return diag.Errorf("error unsubscribe bandwidth: %s", err)
		}
	} else {
		if err := bandwidths.Delete(networkingClient, bwID).ExtractErr(); err != nil {
			return diag.Errorf("error deleting bandwidth: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NORMAL"},
		Target:     []string{"DELETED"},
		Refresh:    waitForBandwidth(bwClient, bwID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for bandwidth (%s) to become DELETED: %s",
			bwID, err)
	}

	return nil
}

func waitForBandwidth(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		b, err := bandwidthsv1.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return b, "DELETED", nil
			}
			return nil, "", err
		}

		log.Printf("[DEBUG] the current status of bandwidth (%s) is %s", b.ID, b.Status)
		return b, b.Status, nil
	}
}

func normalizeChargingMode(billing string) string {
	if billing != "" {
		return "prePaid"
	}
	return "postPaid"
}

func flattenPublicIPs(band bandwidthsv1.BandWidth) []map[string]interface{} {
	allIPs := make([]map[string]interface{}, len(band.PublicipInfo))
	for i, ipInfo := range band.PublicipInfo {
		address := ipInfo.PublicipAddress
		if ipInfo.Publicipv6Address != "" {
			address = ipInfo.Publicipv6Address
		}

		allIPs[i] = map[string]interface{}{
			"id":         ipInfo.PublicipId,
			"type":       ipInfo.PublicipType,
			"ip_version": ipInfo.IPVersion,
			"ip_address": address,
		}
	}

	return allIPs
}
