package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/bss/v2/orders"
	sdk_structs "github.com/huaweicloud/golangsdk/openstack/common/structs"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/eips"

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
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"publicip": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
			"tags": tagsSchema(),
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// charge info: charging_mode, period_unit, period, auto_renew
			"charging_mode": schemeChargingMode(nil),
			"period_unit":   schemaPeriodUnit([]string{"publicip.0.ip_address"}),
			"period":        schemaPeriod([]string{"publicip.0.ip_address"}),
			"auto_renew":    schemaAutoRenew([]string{"publicip.0.ip_address"}),
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
	region := GetRegion(d, config)

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

	epsID := GetEnterpriseProjectID(d, config)
	if epsID != "" {
		createOpts.EnterpriseProjectID = epsID
	}

	var prePaid bool
	if d.Get("charging_mode").(string) == "prePaid" {
		if err := validatePrePaidChargeInfo(d); err != nil {
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
			IsAutoPay:   "true",
			IsAutoRenew: d.Get("auto_renew").(string),
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
		bssClient, err := config.BssV2Client(GetRegion(d, config))
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

	err = bindToPort(d, eIP.ID, networkingClient, timeout)
	if err != nil {
		return fmtp.Errorf("Error binding eip:%s to port: %s", eIP.ID, err)
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
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating networking client: %s", err)
	}

	eIP, err := eips.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "eIP")
	}
	bandWidth, err := bandwidths.Get(networkingClient, eIP.BandwidthID).Extract()
	if err != nil {
		return fmtp.Errorf("Error fetching bandwidth: %s", err)
	}

	// Set public ip
	publicIP := []map[string]string{
		{
			"type":       eIP.Type,
			"ip_address": eIP.PublicAddress,
			"port_id":    eIP.PortID,
		},
	}
	d.Set("publicip", publicIP)

	// Set bandwidth
	bW := []map[string]interface{}{
		{
			"name":        bandWidth.Name,
			"size":        eIP.BandwidthSize,
			"id":          eIP.BandwidthID,
			"share_type":  eIP.BandwidthShareType,
			"charge_mode": bandWidth.ChargeMode,
		},
	}
	d.Set("bandwidth", bW)
	d.Set("address", eIP.PublicAddress)
	d.Set("enterprise_project_id", eIP.EnterpriseProjectID)
	d.Set("region", GetRegion(d, config))
	d.Set("status", normalizeEIPStatus(eIP.Status))

	// save tags
	if vpcV2Client, err := config.NetworkingV2Client(GetRegion(d, config)); err == nil {
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
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating networking client: %s", err)
	}

	// Update bandwidth change
	if d.HasChange("bandwidth") {
		var updateOpts bandwidths.UpdateOpts

		newBWList := d.Get("bandwidth").([]interface{})
		newMap := newBWList[0].(map[string]interface{})
		updateOpts.Size = newMap["size"].(int)
		updateOpts.Name = newMap["name"].(string)

		logp.Printf("[DEBUG] Bandwidth Update Options: %#v", updateOpts)

		eIP, err := eips.Get(networkingClient, d.Id()).Extract()
		if err != nil {
			return CheckDeleted(d, err, "eIP")
		}
		_, err = bandwidths.Update(networkingClient, eIP.BandwidthID, updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating bandwidth: %s", err)
		}

	}

	// Update publicip change
	if d.HasChange("publicip") {
		var updateOpts eips.UpdateOpts

		newIPList := d.Get("publicip").([]interface{})
		newMap := newIPList[0].(map[string]interface{})
		updateOpts.PortID = newMap["port_id"].(string)

		logp.Printf("[DEBUG] PublicIP Update Options: %#v", updateOpts)
		_, err = eips.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating publicip: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		vpcV2Client, err := config.NetworkingV2Client(GetRegion(d, config))
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
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating VPC client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	err = unbindToPort(d, d.Id(), networkingClient, timeout)
	if err != nil {
		logp.Printf("[WARN] Error trying to unbind eip %s :%s", d.Id(), err)
	}

	id := d.Id()
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err := UnsubscribePrePaidResource(d, config, []string{id}); err != nil {
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
		Type:    rawMap["type"].(string),
		Address: rawMap["ip_address"].(string),
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

func bindToPort(d *schema.ResourceData, eipID string, networkingClient *golangsdk.ServiceClient, timeout time.Duration) error {
	publicIPRaw := d.Get("publicip").([]interface{})
	rawMap := publicIPRaw[0].(map[string]interface{})
	portID, ok := rawMap["port_id"]
	if !ok || portID == "" {
		return nil
	}

	pd := portID.(string)
	logp.Printf("[DEBUG] Bind eip:%s to port: %s", eipID, pd)

	updateOpts := eips.UpdateOpts{PortID: pd}
	_, err := eips.Update(networkingClient, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}
	return waitForEIPActive(networkingClient, eipID, timeout)
}

func unbindToPort(d *schema.ResourceData, eipID string, networkingClient *golangsdk.ServiceClient, timeout time.Duration) error {
	publicIPRaw := d.Get("publicip").([]interface{})
	rawMap := publicIPRaw[0].(map[string]interface{})
	portID, ok := rawMap["port_id"]
	if !ok || portID == "" {
		return nil
	}

	pd := portID.(string)
	logp.Printf("[DEBUG] Unbind eip:%s to port: %s", eipID, pd)

	updateOpts := eips.UpdateOpts{PortID: ""}
	_, err := eips.Update(networkingClient, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}
	return waitForEIPActive(networkingClient, eipID, timeout)
}

func normalizeEIPStatus(status string) string {
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
