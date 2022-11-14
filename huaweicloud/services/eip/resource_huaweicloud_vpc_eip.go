package eip

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	sdk_structs "github.com/chnsz/golangsdk/openstack/common/structs"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	bandwidthsv2 "github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type BgpType string         // The BGP type of the public IP.
type IpVersion int          // The Version of the EIP protocol.
type BandwidthType string   // The bandwidth type bound by EIP.
type ChargeMode string      // The charging mode of the bandwidth.
type EipStatus string       // The current status of the EIP.
type NormalizeStatus string // The Normalized status value.

const (
	BgpTypeDynamic BgpType = "5_bgp"  // Dynamic BGP
	BgpTypeStatic  BgpType = "5_sbgp" // Static BGP

	IpVersionV4 IpVersion = 4 // IPv4
	IpVersionV6 IpVersion = 6 // IPv6

	BandwidthTypeDedicated BandwidthType = "PER"   // Dedicated bandwidth
	BandwidthTypeShared    BandwidthType = "WHOLE" // Shared bandwidth

	ChargeModeTraffic   ChargeMode = "traffic"   // Billing based on traffic
	ChargeModeBandwidth ChargeMode = "bandwidth" // Billing based on bandwidth

	EipStatusDown   EipStatus = "DOWN"
	EipStatusActive EipStatus = "ACTIVE"

	NormalizeStatusBound   NormalizeStatus = "BOUND"
	NormalizeStatusUnbound NormalizeStatus = "UNBOUND"
)

// In order to be compatible with other providers, keep the exposed function name (ResourceVpcEIPV1) unchanged.
func ResourceVpcEIPV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcEipCreate,
		ReadContext:   resourceVpcEipRead,
		UpdateContext: resourceVpcEipUpdate,
		DeleteContext: resourceVpcEipDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to create the EIP resource.`,
			},
			"publicip": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(BgpTypeDynamic),
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(BgpTypeDynamic), string(BgpTypeStatic),
							}, false),
							Description: `The EIP type.`,
						},
						"ip_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validation.IsIPv4Address,
							Description:  `The EIP address to be assigned.`,
						},
						"ip_version": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.IntInSlice([]int{
								int(IpVersionV4), int(IpVersionV6),
							}),
							Description: `The IP version.`,
						},
						"port_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Deprecated",
						},
					},
				},
				Description: `The EIP configuration.`,
			},
			"bandwidth": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(BandwidthTypeDedicated), string(BandwidthTypeShared),
							}, false),
							Description: `Whether the bandwidth is dedicated or shared.`,
						},
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ExactlyOneOf: []string{"bandwidth.0.name"},
							Description:  `The shared bandwidth ID.`,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5\\w-.]*$"),
									"The name can only contain letters, digits, underscores (_), hyphens (-), and periods (.)."),
								validation.StringLenBetween(1, 64),
							),
							Description: `The bandwidth name.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The bandwidth size.`,
						},
						"charge_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(ChargeModeTraffic), string(ChargeModeBandwidth),
							}, false),
							Description: `Whether the bandwidth is billed by traffic or by bandwidth size.`,
						},
					},
				},
				Description: `The bandwidth configuration.`,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5\\w-.]*$"),
						"The name can only contain letters, digits, underscores (_), hyphens (-), and periods (.)."),
					validation.StringLenBetween(1, 64),
				),
				Description: `The name of the EIP.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The enterprise project ID to which the EIP belongs.`,
			},
			"tags": common.TagsSchema(),

			// Charge info: charging_mode, period_unit, period, auto_renew, auto_pay
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit([]string{"publicip.0.ip_address"}),
			"period":        common.SchemaPeriod([]string{"publicip.0.ip_address"}),
			"auto_renew":    common.SchemaAutoRenewUpdatable([]string{"publicip.0.ip_address"}),
			"auto_pay":      common.SchemaAutoPay([]string{"publicip.0.ip_address"}),

			// Attributes
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func validatePrePaidBandWidth(bandwidth eips.BandwidthOpts) error {
	if bandwidth.Id != "" || bandwidth.Name == "" || bandwidth.ShareType == string(BandwidthTypeShared) {
		return fmt.Errorf("shared bandwidth is not supported in prePaid charging mode")
	}
	if bandwidth.ChargeMode == string(ChargeModeTraffic) {
		return fmt.Errorf("the EIP can only be billed by bandwidth in prePaid charging mode")
	}

	return nil
}

func buildVpcEipCreateOpts(config *config.Config, d *schema.ResourceData, isPrePaid bool) (eips.ApplyOpts, error) {
	bandwidth := resourceBandWidth(d)
	result := eips.ApplyOpts{
		IP:                  resourcePublicIP(d),
		Bandwidth:           bandwidth,
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
	}

	if isPrePaid {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return result, err
		}
		if err := validatePrePaidBandWidth(bandwidth); err != nil {
			return result, err
		}

		result.ExtendParam = &sdk_structs.ChargeInfo{
			ChargeMode:  d.Get("charging_mode").(string),
			PeriodType:  d.Get("period_unit").(string),
			PeriodNum:   d.Get("period").(int),
			IsAutoRenew: d.Get("auto_renew").(string),
			IsAutoPay:   common.GetAutoPay(d),
		}
	}
	return result, nil
}

func createPrePaidEip(ctx context.Context, config *config.Config, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	timeout := d.Timeout(schema.TimeoutCreate)
	createOpts, err := buildVpcEipCreateOpts(config, d, true)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	resp, err := eips.Apply(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error allocating prePaid EIP: %s", err)
	}

	// Waiting for EIP creation completed.
	bssClient, err := config.BssV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}
	if err := common.WaitOrderComplete(ctx, bssClient, resp.OrderID, timeout); err != nil {
		return err
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, resp.OrderID, timeout)
	if err != nil {
		return err
	}

	d.SetId(resourceId)
	return nil
}

func createPostPaidEip(ctx context.Context, config *config.Config, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	createOpts, err := buildVpcEipCreateOpts(config, d, false)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	resp, err := eips.Apply(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error allocating EIP: %s", err)
	}
	d.SetId(resp.ID)

	log.Printf("[DEBUG] Waiting for EIP (%s) to become available", resp.ID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      eipStatusRefreshFunc(client, resp.ID, []string{"DOWN", "ACTIVE"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for EIP (%s) to become ready: %s", resp.ID, err)
	}
	return nil
}

func resourceVpcEipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	vpcV1Client, err := config.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v1 client: %s", err)
	}
	vpcV2Client, err := config.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v2.0 client: %s", err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		err := createPrePaidEip(ctx, config, vpcV2Client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err := createPostPaidEip(ctx, config, vpcV1Client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("publicip.0.port_id"); ok {
		err = updateEipPortId(vpcV1Client, d)
		if err != nil {
			return diag.Errorf("error binding EIP (%s) to port %s: %s", d.Id(), v.(string), err)
		}
	}

	// create tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcV2Client, "publicips", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of EIP (%s): %s", d.Id(), tagErr)
		}
	}

	return resourceVpcEipRead(ctx, d, meta)
}

// NormalizeEipStatus is a method to change an incomprehensible status into an easy-to-understand status.
func NormalizeEipStatus(status string) string {
	// The 'DOWN' status means the EIP is active but not bound.
	if status == string(EipStatusDown) {
		return string(NormalizeStatusUnbound)
	}
	if status == string(EipStatusActive) {
		return string(NormalizeStatusBound)
	}

	// Other running statuses.
	return status
}

func eipStatusRefreshFunc(networkingClient *golangsdk.ServiceClient, eipId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := eips.Get(networkingClient, eipId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				if len(targets) < 1 {
					return resp, "COMPLETED", nil
				}
				return resp, "PENDING", nil
			}

			return nil, "", err
		}
		log.Printf("[DEBUG] The details of the EIP (%s) is: %+v", eipId, resp)

		if utils.StrSliceContains([]string{"BIND_ERROR", "ERROR"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpected status '%s'", resp.Status)
		}
		if utils.StrSliceContains(targets, resp.Status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func flattenEipPublicIpDetails(publicIp eips.PublicIp) []map[string]interface{} {
	if reflect.DeepEqual(publicIp, eips.PublicIp{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":       publicIp.Type,
			"ip_version": publicIp.IpVersion,
			"ip_address": publicIp.PublicAddress,
			"port_id":    publicIp.PortID,
		},
	}
}

func flattenEipBandwidthDetails(publicIp eips.PublicIp, bandWidth bandwidths.BandWidth) []map[string]interface{} {
	if reflect.DeepEqual(publicIp, eips.PublicIp{}) || reflect.DeepEqual(bandWidth, bandwidths.BandWidth{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":        bandWidth.Name,
			"size":        publicIp.BandwidthSize,
			"id":          publicIp.BandwidthID,
			"share_type":  publicIp.BandwidthShareType,
			"charge_mode": bandWidth.ChargeMode,
		},
	}
}

func resourceVpcEipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	networkingClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v1 client: %s", err)
	}

	resourceId := d.Id()
	publicIp, err := eips.Get(networkingClient, resourceId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "EIP")
	}
	bandWidth, err := bandwidths.Get(networkingClient, publicIp.BandwidthID).Extract()
	if err != nil {
		return diag.Errorf("error fetching bandwidth: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", publicIp.Alias),
		d.Set("address", publicIp.PublicAddress),
		d.Set("ipv6_address", publicIp.PublicIpv6Address),
		d.Set("private_ip", publicIp.PrivateAddress),
		d.Set("port_id", publicIp.PortID),
		d.Set("enterprise_project_id", publicIp.EnterpriseProjectID),
		d.Set("status", NormalizeEipStatus(publicIp.Status)),
		d.Set("publicip", flattenEipPublicIpDetails(publicIp)),
		d.Set("bandwidth", flattenEipBandwidthDetails(publicIp, bandWidth)),
	)

	// Save tags via EIP API.
	if vpcV2Client, err := config.NetworkingV2Client(region); err == nil {
		if resourceTags, err := tags.Get(vpcV2Client, "publicips", resourceId).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			if err := d.Set("tags", tagmap); err != nil {
				mErr = multierror.Append(mErr, fmt.Errorf("error saving tags for EIP (%s): %s", resourceId, err))
			}
		} else {
			log.Printf("[WARN] Error fetching tags for EIP (%s): %s", resourceId, err)
		}
	} else {
		mErr = multierror.Append(mErr, fmt.Errorf("error creating VPC v2.0 client: %s", err))
	}

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}
	return nil
}

func updateEipConfig(vpcV1Client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var updateOpts = eips.UpdateOpts{}

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		updateOpts.Alias = &newName
	}
	if d.HasChange("publicip.0.ip_version") {
		updateOpts.IPVersion = d.Get("publicip.0.ip_version").(int)
	}

	if updateOpts != (eips.UpdateOpts{}) {
		log.Printf("[DEBUG] PublicIP Update Options: %#v", updateOpts)
		_, err := eips.Update(vpcV1Client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating public IP: %s", err)
		}
	}
	return nil
}

func updateEipPortId(vpcV1Client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	resourceId := d.Id()
	timeout := d.Timeout(schema.TimeoutUpdate)
	old, new := d.GetChange("publicip.0.port_id")
	oldPort := old.(string)
	newPort := new.(string)

	if oldPort != "" {
		err := unbindPort(vpcV1Client, resourceId, oldPort, timeout)
		if err != nil {
			log.Printf("[WARN] Error trying to unbind EIP (%s): %s", resourceId, err)
		}
	}
	if newPort != "" {
		err := bindPort(vpcV1Client, resourceId, newPort, timeout)
		if err != nil {
			return fmt.Errorf("error binding EIP (%s) to port (%s): %s", resourceId, newPort, err)
		}
	}
	return nil
}

func updateEipBandwidth(vpcV1Client *golangsdk.ServiceClient, config *config.Config, d *schema.ResourceData) error {
	old, new := d.GetChange("bandwidth")
	oldRaw := old.([]interface{})
	newRaw := new.([]interface{})
	// Bandwidth blocks are required and must be present.
	oldMap := oldRaw[0].(map[string]interface{})
	newMap := newRaw[0].(map[string]interface{})

	bandwidthId := oldMap["id"].(string)
	if d.Get("charging_mode").(string) == "prePaid" {
		vpcV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating VPC v2.0 client: %s", err)
		}
		updateOpts := bandwidthsv2.UpdateOpts{
			Bandwidth: bandwidthsv2.Bandwidth{
				Name: newMap["name"].(string),
				Size: newMap["size"].(int),
			},
			ExtendParam: &bandwidthsv2.ExtendParam{
				IsAutoPay: common.GetAutoPay(d),
			},
		}
		log.Printf("[DEBUG] Bandwidth Update Options: %#v", updateOpts)

		order, err := bandwidthsv2.Update(vpcV2Client, bandwidthId, updateOpts)
		if err != nil {
			return fmt.Errorf("error updating bandwidth: %s", err)
		}

		if orderData, ok := order.(bandwidthsv2.PrePaid); ok {
			log.Printf("[DEBUG] The orderData is: %#v", orderData)
			// Wait for order success.
			bssClient, err := config.BssV2Client(config.GetRegion(d))
			if err != nil {
				return fmt.Errorf("error creating BSS v2 client: %s", err)
			}
			if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), orderData.OrderID); err != nil {
				return err
			}
		} else if d.HasChange("bandwidth.0.size") {
			return fmt.Errorf("unable to find order structure in API response: %#v", order)
		}
		return nil
	}

	updateOpts := bandwidths.UpdateOpts{
		Size: newMap["size"].(int),
		Name: newMap["name"].(string),
	}
	log.Printf("[DEBUG] Bandwidth Update Options: %#v", updateOpts)
	_, err := bandwidths.Update(vpcV1Client, bandwidthId, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating bandwidth: %s", err)
	}

	return nil
}

func resourceVpcEipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcV1Client, err := config.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v1 client: %s", err)
	}

	// API limitation: port_id and ip_version cannot be updated at the same time
	if d.HasChanges("name", "publicip.0.ip_version") {
		err = updateEipConfig(vpcV1Client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("publicip.0.port_id") {
		err = updateEipPortId(vpcV1Client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("bandwidth") {
		err = updateEipBandwidth(vpcV1Client, config, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		vpcV2Client, err := config.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v2 client: %s", err)
		}

		tagErr := utils.UpdateResourceTags(vpcV2Client, d, "publicips", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPC (%s): %s", d.Id(), tagErr)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := config.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the EIP (%s): %s", d.Id(), err)
		}
	}

	return resourceVpcEipRead(ctx, d, meta)
}

func resourceVpcEipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	resourceId := d.Id()

	// check whether the eip exists before delete it
	// because resource could not be found cannot be deleteed
	_, err = eips.Get(networkingClient, resourceId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EIP")
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	if v, ok := d.GetOk("publicip.0.port_id"); ok {
		portID := v.(string)
		err = unbindPort(networkingClient, resourceId, portID, timeout)
		if err != nil {
			log.Printf("[WARN] Error trying to unbind eip %s :%s", resourceId, err)
		}
	}

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, config, []string{resourceId}); err != nil {
			return diag.Errorf("error unsubscribe publicip: %s", err)
		}
	} else {
		if err := eips.Delete(networkingClient, resourceId).ExtractErr(); err != nil {
			return diag.Errorf("error deleting publicip: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"COMPLETED"},
		Refresh:    eipStatusRefreshFunc(networkingClient, resourceId, nil),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for EIP (%s) to be deleted: %s", resourceId, err)
	}

	d.SetId("")
	return nil
}

func resourcePublicIP(d *schema.ResourceData) eips.PublicIpOpts {
	publicIPRaw := d.Get("publicip").([]interface{})
	rawMap := publicIPRaw[0].(map[string]interface{})

	publicip := eips.PublicIpOpts{
		Alias:     d.Get("name").(string),
		Type:      rawMap["type"].(string),
		Address:   rawMap["ip_address"].(string),
		IPVersion: rawMap["ip_version"].(int),
	}
	return publicip
}

func resourceBandWidth(d *schema.ResourceData) eips.BandwidthOpts {
	bandwidthRaw := d.Get("bandwidth").([]interface{})
	rawMap := bandwidthRaw[0].(map[string]interface{})

	bandwidth := eips.BandwidthOpts{
		Id:         rawMap["id"].(string),
		Name:       rawMap["name"].(string),
		Size:       rawMap["size"].(int),
		ShareType:  rawMap["share_type"].(string),
		ChargeMode: rawMap["charge_mode"].(string),
	}
	return bandwidth
}
