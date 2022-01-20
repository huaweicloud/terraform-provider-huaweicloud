package eip

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk"
	bandwidthsv1 "github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(5, 2000),
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bandwidth", "95peak_plus",
				}, false),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

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
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	NetworkingV1Client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating networking client: %s", err)
	}

	size := d.Get("size").(int)
	createOpts := bandwidths.CreateOpts{
		Name:       d.Get("name").(string),
		ChargeMode: d.Get("charge_mode").(string),
		Size:       &size,
	}

	epsID := config.GetEnterpriseProjectID(d)
	if epsID != "" {
		createOpts.EnterpriseProjectId = epsID
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	b, err := bandwidths.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating Bandwidth: %s", err)
	}

	logp.Printf("[DEBUG] Waiting for Bandwidth (%s) to become available.", b.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"NORMAL"},
		Pending:    []string{"CREATING"},
		Refresh:    waitForBandwidth(NetworkingV1Client, b.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf(
			"Error waiting for Bandwidth (%s) to become ACTIVE for creation: %s",
			b.ID, err)
	}
	d.SetId(b.ID)

	return resourceVpcBandWidthV2Read(ctx, d, meta)
}

func resourceVpcBandWidthV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating networking client: %s", err)
	}

	if d.HasChanges("name", "size") {
		updateOpts := bandwidths.UpdateOpts{
			Bandwidth: bandwidths.Bandwidth{
				Name: d.Get("name").(string),
				Size: d.Get("size").(int),
			},
		}
		_, err := bandwidths.Update(networkingClient, d.Id(), updateOpts)
		if err != nil {
			return fmtp.DiagErrorf("Error updating HuaweiCloud BandWidth (%s): %s", d.Id(), err)
		}
	}

	return resourceVpcBandWidthV2Read(ctx, d, meta)
}

func resourceVpcBandWidthV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating networking client: %s", err)
	}

	b, err := bandwidthsv1.Get(networkingClient, d.Id()).Extract()
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
		d.Set("status", b.Status),
		d.Set("publicips", flattenPublicIPs(b)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting bandwidth fields: %s", err)
	}

	return nil
}

func resourceVpcBandWidthV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	NetworkingV1Client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating networking client: %s", err)
	}

	err = bandwidths.Delete(networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud Bandwidth: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NORMAL"},
		Target:     []string{"DELETED"},
		Refresh:    waitForBandwidth(NetworkingV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting Bandwidth: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForBandwidth(networkingClient *golangsdk.ServiceClient, Id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		b, err := bandwidthsv1.Get(networkingClient, Id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return b, "DELETED", nil
			}
			return nil, "", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Bandwidth (%s) current status: %s", b.ID, b.Status)
		return b, b.Status, nil
	}
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
