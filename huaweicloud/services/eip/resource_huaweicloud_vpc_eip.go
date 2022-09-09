package eip

import (
	"time"

	"github.com/hashicorp/go-multierror"
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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceVpcEIPV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcEIPV1Create,
		Read:   resourceVpcEIPV1Read,
		Update: resourceVpcEIPV1Update,
		Delete: resourceVpcEIPV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
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
				Optional: true,
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
							Default:  "5_bgp",
							ForceNew: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"ip_version": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntInSlice([]int{4, 6}),
						},
						"port_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Deprecated",
						},
					},
				},
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
						},
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ExactlyOneOf: []string{"bandwidth.0.id", "bandwidth.0.name"},
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"charge_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
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

			// charge info: charging_mode, period_unit, period, auto_renew, auto_pay
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit([]string{"publicip.0.ip_address"}),
			"period":        common.SchemaPeriod([]string{"publicip.0.ip_address"}),
			"auto_renew":    common.SchemaAutoRenew([]string{"publicip.0.ip_address"}),
			"auto_pay":      common.SchemaAutoPay([]string{"publicip.0.ip_address"}),
		},
	}
}

func validatePrePaidBandWidth(bandwidth eips.BandwidthOpts) error {
	if bandwidth.Id != "" || bandwidth.Name == "" || bandwidth.ShareType == "WHOLE" {
		return fmtp.Errorf("shared bandwidth is not supported in prePaid charging mode")
	}
	if bandwidth.ChargeMode == "traffic" {
		return fmtp.Errorf("EIP can only be billed by bandwidth in prePaid charging mode")
	}

	return nil
}

func validatePrePaidSupportedRegion(region string) error {
	var valid bool
	// reference to: https://support.huaweicloud.com/api-eip/eip_api_0006.html#section4
	var supportedRegion = []string{
		"cn-north-4", "cn-east-3", "cn-south-1", "cn-southwest-2",
		"ap-southeast-2", "ap-southeast-3",
	}

	for _, r := range supportedRegion {
		if r == region {
			valid = true
			break
		}
	}
	if !valid {
		return fmtp.Errorf("prepaid charging mode of eip is not supported in %s region", region)
	}

	return nil
}

func resourceVpcEIPV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	networkingClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating networking v1 client: %s", err)
	}
	// networkingV2Client is used to create EIP in prePaid charging mode and tags
	networkingV2Client, err := config.NetworkingV2Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating networking v2 client: %s", err)
	}

	createOpts := eips.ApplyOpts{
		IP:        resourcePublicIP(d),
		Bandwidth: resourceBandWidth(d),
	}

	epsID := config.GetEnterpriseProjectID(d)
	if epsID != "" {
		createOpts.EnterpriseProjectID = epsID
	}

	var prePaid bool
	if d.Get("charging_mode").(string) == "prePaid" {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return err
		}
		if err := validatePrePaidBandWidth(createOpts.Bandwidth); err != nil {
			return err
		}
		if err := validatePrePaidSupportedRegion(region); err != nil {
			return err
		}

		prePaid = true
		chargeInfo := &sdk_structs.ChargeInfo{
			ChargeMode:  d.Get("charging_mode").(string),
			PeriodType:  d.Get("period_unit").(string),
			PeriodNum:   d.Get("period").(int),
			IsAutoRenew: d.Get("auto_renew").(string),
			IsAutoPay:   common.GetAutoPay(d),
		}
		createOpts.ExtendParam = chargeInfo
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	var eIP eips.PublicIp
	if prePaid {
		eIP, err = eips.Apply(networkingV2Client, createOpts).Extract()
	} else {
		eIP, err = eips.Apply(networkingClient, createOpts).Extract()
	}

	if err != nil {
		return fmtp.Errorf("Error allocating EIP: %s", err)
	}

	if eIP.ID == "" {
		return fmtp.Errorf("can not get the resource ID")
	}
	d.SetId(eIP.ID)

	// wait for order success
	if eIP.OrderID != "" {
		bssClient, err := config.BssV2Client(config.GetRegion(d))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud bss V2 client: %s", err)
		}
		if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), eIP.OrderID); err != nil {
			return err
		}
	}

	logp.Printf("[DEBUG] Waiting for EIP %s to become available.", eIP.ID)
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForEIPActive(networkingClient, eIP.ID, timeout)
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for EIP (%s) to become ready: %s",
			eIP.ID, err)
	}

	if v, ok := d.GetOk("publicip.0.port_id"); ok {
		portID := v.(string)
		err = bindPort(networkingClient, eIP.ID, portID, timeout)
		if err != nil {
			return fmtp.Errorf("Error binding EIP %s to port %s: %s", eIP.ID, portID, err)
		}
	}

	// create tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(networkingV2Client, "publicips", eIP.ID, taglist).ExtractErr(); tagErr != nil {
			return fmtp.Errorf("Error setting tags of EIP %s: %s", eIP.ID, tagErr)
		}
	}

	return resourceVpcEIPV1Read(d, meta)
}

func resourceVpcEIPV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	networkingClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating networking client: %s", err)
	}

	eIP, err := eips.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "eIP")
	}
	bandWidth, err := bandwidths.Get(networkingClient, eIP.BandwidthID).Extract()
	if err != nil {
		return fmtp.Errorf("Error fetching bandwidth: %s", err)
	}

	// build public ip
	publicIP := []map[string]interface{}{
		{
			"type":       eIP.Type,
			"ip_version": eIP.IpVersion,
			"ip_address": eIP.PublicAddress,
			"port_id":    eIP.PortID,
		},
	}

	// build bandwidth
	bW := []map[string]interface{}{
		{
			"name":        bandWidth.Name,
			"size":        eIP.BandwidthSize,
			"id":          eIP.BandwidthID,
			"share_type":  eIP.BandwidthShareType,
			"charge_mode": bandWidth.ChargeMode,
		},
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", eIP.Alias),
		d.Set("address", eIP.PublicAddress),
		d.Set("ipv6_address", eIP.PublicIpv6Address),
		d.Set("private_ip", eIP.PrivateAddress),
		d.Set("port_id", eIP.PortID),
		d.Set("enterprise_project_id", eIP.EnterpriseProjectID),
		d.Set("status", NormalizeEIPStatus(eIP.Status)),
		d.Set("publicip", publicIP),
		d.Set("bandwidth", bW),
	)

	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	// save tags
	if vpcV2Client, err := config.NetworkingV2Client(region); err == nil {
		if resourceTags, err := tags.Get(vpcV2Client, "publicips", d.Id()).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			if err := d.Set("tags", tagmap); err != nil {
				return fmtp.Errorf("Error saving tags for EIP (%s): %s", d.Id(), err)
			}
		} else {
			logp.Printf("[WARN] Error fetching tags for EIP (%s): %s", d.Id(), err)
		}
	} else {
		return fmtp.Errorf("Error creating vpc client: %s", err)
	}

	return nil
}

func resourceVpcEIPV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating networking client: %s", err)
	}

	// Update bandwidth change
	if d.HasChange("bandwidth") {
		newBWList := d.Get("bandwidth").([]interface{})
		newMap := newBWList[0].(map[string]interface{})

		// Check if eip exists
		eIP, err := eips.Get(networkingClient, d.Id()).Extract()
		if err != nil {
			return common.CheckDeleted(d, err, "eIP")
		}

		if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
			networkingV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
			if err != nil {
				return fmtp.Errorf("error creating networking v2 client: %s", err)
			}
			bandwidth := bandwidthsv2.Bandwidth{
				Name: newMap["name"].(string),
				Size: newMap["size"].(int),
			}
			extendParam := &bandwidthsv2.ExtendParam{
				IsAutoPay: common.GetAutoPay(d),
			}
			updateOpts := bandwidthsv2.UpdateOpts{
				Bandwidth:   bandwidth,
				ExtendParam: extendParam,
			}
			logp.Printf("[DEBUG] Bandwidth Update Options: %#v", updateOpts)

			order, err := bandwidthsv2.Update(networkingV2Client, eIP.BandwidthID, updateOpts)
			if err != nil {
				return fmtp.Errorf("error updating bandwidth: %s", err)
			}

			orderData, ok := order.(bandwidthsv2.PrePaid)
			if ok {
				logp.Printf("[DEBUG] The orderData is: %#v", orderData)
				// wait for order success
				bssClient, err := config.BssV2Client(config.GetRegion(d))
				if err != nil {
					return fmtp.Errorf("error creating BSS v2 client: %s", err)
				}
				if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), orderData.OrderID); err != nil {
					return err
				}
			}
		} else {
			var updateOpts bandwidths.UpdateOpts
			updateOpts.Size = newMap["size"].(int)
			updateOpts.Name = newMap["name"].(string)
			logp.Printf("[DEBUG] Bandwidth Update Options: %#v", updateOpts)

			_, err = bandwidths.Update(networkingClient, eIP.BandwidthID, updateOpts).Extract()
			if err != nil {
				return fmtp.Errorf("error updating bandwidth: %s", err)
			}
		}
	}

	newIPList := d.Get("publicip").([]interface{})
	newMap := newIPList[0].(map[string]interface{})

	// Update publicip port_id
	if d.HasChange("publicip.0.port_id") {
		timeout := d.Timeout(schema.TimeoutUpdate)
		old, new := d.GetChange("publicip.0.port_id")
		oldPort := old.(string)
		newPort := new.(string)

		if oldPort != "" {
			err = unbindPort(networkingClient, d.Id(), oldPort, timeout)
			if err != nil {
				logp.Printf("[WARN] Error trying to unbind EIP %s :%s", d.Id(), err)
			}
		}

		if newPort != "" {
			err = bindPort(networkingClient, d.Id(), newPort, timeout)
			if err != nil {
				return fmtp.Errorf("Error binding EIP %s to port %s: %s", d.Id(), newPort, err)
			}
		}
	}

	// Update publicip name and ip_version
	// API limitation: port_id and ip_version cannot be updated at the same time
	if d.HasChanges("name", "publicip.0.ip_version") {
		newName := d.Get("name").(string)
		updateOpts := eips.UpdateOpts{
			Alias:     &newName,
			IPVersion: newMap["ip_version"].(int),
		}

		logp.Printf("[DEBUG] PublicIP Update Options: %#v", updateOpts)
		_, err = eips.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating publicip: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		vpcV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud vpc client: %s", err)
		}

		tagErr := utils.UpdateResourceTags(vpcV2Client, d, "publicips", d.Id())
		if tagErr != nil {
			return fmtp.Errorf("Error updating tags of VPC %s: %s", d.Id(), tagErr)
		}
	}

	return resourceVpcEIPV1Read(d, meta)
}

func resourceVpcEIPV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating VPC client: %s", err)
	}

	id := d.Id()

	// check whether the eip exists before delete it
	// because resource could not be found cannot be deleteed
	_, err = eips.Get(networkingClient, id).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Error retrieving HuaweiCloud EIP")
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	if v, ok := d.GetOk("publicip.0.port_id"); ok {
		portID := v.(string)
		err = unbindPort(networkingClient, id, portID, timeout)
		if err != nil {
			logp.Printf("[WARN] Error trying to unbind eip %s :%s", id, err)
		}
	}

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, config, []string{id}); err != nil {
			return fmtp.Errorf("Error unsubscribe publicip: %s", err)
		}
	} else {
		if err := eips.Delete(networkingClient, id).ExtractErr(); err != nil {
			return fmtp.Errorf("Error deleting publicip: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForEIPDelete(networkingClient, d.Id()),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting EIP: %s", err)
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

func NormalizeEIPStatus(status string) string {
	var ret string = status

	// "DOWN" means the eip is active but unbound
	if status == "DOWN" {
		ret = "UNBOUND"
	} else if status == "ACTIVE" {
		ret = "BOUND"
	}

	return ret
}

func getEIPStatus(networkingClient *golangsdk.ServiceClient, eId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		e, err := eips.Get(networkingClient, eId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil, "PENDING_CREATE", nil
			}

			return nil, "", err
		}

		logp.Printf("[DEBUG] EIP: %+v", e)
		if e.Status == "DOWN" || e.Status == "ACTIVE" {
			return e, "ACTIVE", nil
		}

		return e, "", nil
	}
}

func waitForEIPActive(networkingClient *golangsdk.ServiceClient, eipID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    getEIPStatus(networkingClient, eipID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForEIPDelete(networkingClient *golangsdk.ServiceClient, eId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete EIP %s.\n", eId)

		e, err := eips.Get(networkingClient, eId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted EIP %s", eId)
				return e, "DELETED", nil
			}
			return nil, "ERROR", err
		}

		logp.Printf("[DEBUG] EIP %s still active.\n", eId)
		return e, "ACTIVE", nil
	}
}
