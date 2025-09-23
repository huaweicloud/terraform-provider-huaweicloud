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
	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"
	"github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	PublicIPv6Type             = "5_dualStack"
	DefaultBandWidthChargeMode = "bandwidth"
	DefaultBandWidthSize       = 5
)

// @API EIP POST /v2.0/{project_id}/bandwidths/{bandwidth_id}/insert
// @API EIP POST /v2.0/{project_id}/bandwidths/{bandwidth_id}/remove
// @API EIP GET /v1/{project_id}/bandwidths/{id}
// @API EIP GET /v1/{project_id}/publicips/{id}
// @API VPC GET /v1/{project_id}/ports
// @API VPC GET /v1/{project_id}/ports/{portId}
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
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"port_id", "eip_id", "fixed_ip"},
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
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"network_id"},
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_ip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
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
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	bwID := d.Get("bandwidth_id").(string)
	eipID := d.Get("eip_id").(string)

	if eipID != "" {
		err = insertEIPsToBandwidth(d, bwClient, eipClient, bwID, eipID)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceID := fmt.Sprintf("%s/%s", bwID, eipID)
		d.SetId(resourceID)
	} else {
		portID := d.Get("port_id").(string)
		if portID == "" {
			networkID := d.Get("network_id").(string)
			fixedIP := d.Get("fixed_ip").(string)
			portID, err = getPortbyFixedIP(vpcClient, networkID, fixedIP)
			if err != nil {
				return diag.Errorf("unable to get port ID of %s: %s", fixedIP, err)
			}
		}

		if err = insertPortToBandwidth(bwClient, vpcClient, bwID, portID); err != nil {
			return diag.FromErr(err)
		}

		resourceID := fmt.Sprintf("%s/%s", bwID, portID)
		d.SetId(resourceID)

		if err := d.Set("port_id", portID); err != nil {
			return diag.Errorf("error setting port_id: %s", err)
		}
	}

	return resourceBandWidthAssociateRead(ctx, d, meta)
}

func resourceBandWidthAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	eipClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating bandwidth v1 client: %s", err)
	}
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	bwID := d.Get("bandwidth_id").(string)
	eipID := d.Get("eip_id").(string)
	portID := d.Get("port_id").(string)

	b, err := bandwidthsv1.Get(eipClient, bwID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "bandwidth associate")
	}

	associatedID := eipID
	if eipID == "" {
		associatedID = portID
	}
	var associatedItem *bandwidthsv1.PublicIpinfo
	allAssociatedItems := b.PublicipInfo
	for i := range allAssociatedItems {
		if allAssociatedItems[i].PublicipId == associatedID {
			associatedItem = &allAssociatedItems[i]
			break
		}
	}

	if associatedItem == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("bandwidth_id", b.ID),
		d.Set("bandwidth_name", b.Name),
		d.Set("public_ip_type", associatedItem.PublicipType),
		d.Set("ip_version", associatedItem.IPVersion),
		d.Set("public_ip", associatedItem.PublicipAddress),
		d.Set("public_ipv6", associatedItem.Publicipv6Address),
	)

	if eipID == "" {
		associatedPort, err := ports.Get(vpcClient, portID)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error fetching port")
		}
		mErr = multierror.Append(mErr,
			d.Set("port_id", associatedItem.PublicipId),
			d.Set("fixed_ip", associatedItem.Publicipv6Address),
			d.Set("network_id", associatedPort.NetworkId),
			d.Set("eip_id", ""),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("eip_id", associatedItem.PublicipId),
			d.Set("port_id", ""),
			d.Set("fixed_ip", ""),
			d.Set("network_id", ""),
		)
	}

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
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	if d.HasChanges("eip_id", "port_id", "fixed_ip", "network_id") {
		bwID := d.Get("bandwidth_id").(string)
		oldEipId, newEipId := d.GetChange("eip_id")
		oldPortId, newPortId := d.GetChange("port_id")
		newNetworkId := d.Get("network_id").(string)
		newFixedIp := d.Get("fixed_ip").(string)

		// remove old
		if oldEipId.(string) != "" {
			err = removeEIPsFromBandwidth(d, bwClient, eipClient, bwID, oldEipId.(string))
		} else {
			err = removePortFromBandwidth(bwClient, vpcClient, bwID, oldPortId.(string))
		}
		if err != nil {
			return diag.FromErr(err)
		}

		// insert new
		if newEipId.(string) != "" {
			err = insertEIPsToBandwidth(d, bwClient, eipClient, bwID, newEipId.(string))
			if err != nil {
				return diag.FromErr(err)
			}

			// update the resource ID
			resourceID := fmt.Sprintf("%s/%s", bwID, newEipId.(string))
			d.SetId(resourceID)
		} else {
			portID := newPortId.(string)
			if !d.HasChange("port_id") {
				portID, err = getPortbyFixedIP(vpcClient, newNetworkId, newFixedIp)
				if err != nil {
					return diag.Errorf("unable to get port ID of %s: %s", newFixedIp, err)
				}
			}

			if err = insertPortToBandwidth(bwClient, vpcClient, bwID, portID); err != nil {
				return diag.FromErr(err)
			}

			// update the resource ID and port ID
			resourceID := fmt.Sprintf("%s/%s", bwID, portID)
			d.SetId(resourceID)
			d.Set("port_id", portID)
		}
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
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	bwID := d.Get("bandwidth_id").(string)
	eipID := d.Get("eip_id").(string)
	portID := d.Get("port_id").(string)

	if eipID != "" {
		err = removeEIPsFromBandwidth(d, bwClient, eipClient, bwID, eipID)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err = removePortFromBandwidth(bwClient, vpcClient, bwID, portID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func insertPortToBandwidth(bwClient, vpcClient *golangsdk.ServiceClient, bwID, portID string) error {
	associatedPort, err := ports.Get(vpcClient, portID)
	if err != nil {
		return fmt.Errorf("error fetching port %s: %s", portID, err)
	}
	if associatedPort.Ipv6BandwidthId != "" {
		// the port is already associated to the shared bandwidth, return with success.
		if associatedPort.Ipv6BandwidthId == bwID {
			log.Printf("[DEBUG] port %s is already associated to shared bandwidth %s", portID, bwID)
			return nil
		}
		// the port is associated to another shared bandwidth, we should remove it first.
		if err := removePortFromBandwidth(bwClient, vpcClient, associatedPort.Ipv6BandwidthId, portID); err != nil {
			return err
		}
	}

	insertOpts := bandwidths.BandWidthInsertOpts{
		PublicipInfo: []bandwidths.PublicIpInfoID{
			{
				PublicIPID:   portID,
				PublicIPType: PublicIPv6Type,
			},
		},
	}

	log.Printf("[DEBUG] Insert port %s to bandwidth %s", portID, bwID)
	_, err = bandwidths.Insert(bwClient, bwID, insertOpts).Extract()
	if err != nil {
		return fmt.Errorf("error inserting %s into bandwidth %s: %s", portID, bwID, err)
	}

	return nil
}

func removePortFromBandwidth(bwClient, vpcClient *golangsdk.ServiceClient, bwID, portID string) error {
	if _, err := ports.Get(vpcClient, portID); err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] unnecessary to remove the non-existent port %s from bandwidth %s", portID, bwID)
			return nil
		}
		return fmt.Errorf("error fetching port %s: %s", portID, err)
	}

	// chargeMode and size are not required in actual, although they are required in docs
	removalChargeMode := "bandwidth"
	removalSize := 5

	removeOpts := bandwidths.BandWidthRemoveOpts{
		ChargeMode: removalChargeMode,
		Size:       &removalSize,
		PublicipInfo: []bandwidths.PublicIpInfoID{
			{
				PublicIPID:   portID,
				PublicIPType: PublicIPv6Type,
			},
		},
	}

	log.Printf("[DEBUG] Remove port %s from bandwidth %s", portID, bwID)
	err := bandwidths.Remove(bwClient, bwID, removeOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("error removing %s from bandwidth: %s", portID, err)
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

func resourceBandWidthAssociationImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	length := len(parts)
	if length != 2 && length != 3 {
		return nil, fmt.Errorf("invalid format for import ID, want '<bandwidth_id>/<eip_id>' or '<bandwidth_id>/<port_id>',"+
			" or '<bandwidth_id>/<network_id>/<fixed_ip>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("bandwidth_id", parts[0]),
	)

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC client: %s", err)
	}
	eipClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating EIP client: %s", err)
	}

	if length == 3 {
		networkID := parts[1]
		fixedIP := parts[2]
		portID, err := getPortbyFixedIP(vpcClient, networkID, fixedIP)
		if err != nil {
			return nil, fmt.Errorf("unable to get port ID of %s: %s", fixedIP, err)
		}

		resourceID := fmt.Sprintf("%s/%s", parts[0], portID)
		d.SetId(resourceID)
		mErr = multierror.Append(mErr,
			d.Set("network_id", parts[1]),
			d.Set("fixed_ip", parts[2]),
			d.Set("port_id", portID),
		)
	} else {
		if _, err := eips.Get(eipClient, parts[1]).Extract(); err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); !ok {
				return nil, fmt.Errorf("error getting eip(%s): %s", parts[1], err)
			}
			mErr = multierror.Append(mErr,
				d.Set("port_id", parts[1]),
			)
		} else {
			mErr = multierror.Append(mErr,
				d.Set("eip_id", parts[1]),
			)
		}
	}

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
