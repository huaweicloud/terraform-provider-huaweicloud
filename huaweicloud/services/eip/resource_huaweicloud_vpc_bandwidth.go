package eip

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	bandwidthsv1 "github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
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

	size := d.Get("size").(int)
	createOpts := bandwidths.CreateOpts{
		Name:       d.Get("name").(string),
		ChargeMode: d.Get("charge_mode").(string),
		Size:       &size,
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID != "" {
		createOpts.EnterpriseProjectId = epsID
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
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

	return resourceVpcBandWidthV2Read(ctx, d, meta)
}

func resourceVpcBandWidthV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	networkingClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	bwID := d.Id()
	if d.HasChanges("name", "size") {
		updateOpts := bandwidths.UpdateOpts{
			Bandwidth: bandwidths.Bandwidth{
				Name: d.Get("name").(string),
				Size: d.Get("size").(int),
			},
		}
		_, err := bandwidths.Update(networkingClient, bwID, updateOpts)
		if err != nil {
			return diag.Errorf("error updating bandwidth (%s): %s", bwID, err)
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
		d.Set("status", b.Status),
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
	err = bandwidths.Delete(networkingClient, bwID).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting bandwidth: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NORMAL"},
		Target:     []string{"DELETED"},
		Refresh:    waitForBandwidth(bwClient, bwID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
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
