package huaweicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/hw_snatrules"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceNatSnatRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatSnatRuleV2Create,
		Read:   resourceNatSnatRuleV2Read,
		Update: resourceNatSnatRuleV2Update,
		Delete: resourceNatSnatRuleV2Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"floating_ip_id": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: utils.SuppressSnatFiplistDiffs,
			},
			"source_type": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 1),
				Optional:     true,
				ForceNew:     true,
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"cidr", "network_id"},
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"subnet_id", "network_id"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"floating_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// deprecated
			"network_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "use subnet_id instead",
			},
		},
	}
}

func resourceNatSnatRuleV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	natClient, err := config.NatGatewayClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	var subnetID string
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetID = v.(string)
	} else {
		subnetID = d.Get("network_id").(string)
	}

	sourceType := d.Get("source_type").(int)
	if sourceType == 1 && subnetID != "" {
		return fmtp.Errorf("source_type and subnet_id is incompatible in the Direct Connect scenario (source_type=1)")
	}

	createOpts := &hw_snatrules.CreateOpts{
		NatGatewayID: d.Get("nat_gateway_id").(string),
		FloatingIPID: d.Get("floating_ip_id").(string),
		Cidr:         d.Get("cidr").(string),
		Description:  d.Get("description").(string),
		NetworkID:    subnetID,
		SourceType:   sourceType,
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	snatRule, err := hw_snatrules.Create(natClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creatting Snat Rule: %s", err)
	}

	logp.Printf("[DEBUG] Waiting for HuaweiCloud Snat Rule (%s) to become available.", snatRule.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForSnatRuleActive(natClient, snatRule.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud Snat Rule: %s", err)
	}

	d.SetId(snatRule.ID)

	return resourceNatSnatRuleV2Read(d, meta)
}

func resourceNatSnatRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	natClient, err := config.NatGatewayClient(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	snatRule, err := hw_snatrules.Get(natClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Snat Rule")
	}

	d.Set("region", region)
	d.Set("nat_gateway_id", snatRule.NatGatewayID)
	d.Set("floating_ip_id", snatRule.FloatingIPID)
	d.Set("floating_ip_address", snatRule.FloatingIPAddress)
	d.Set("source_type", snatRule.SourceType)
	d.Set("subnet_id", snatRule.NetworkID)
	d.Set("cidr", snatRule.Cidr)
	d.Set("status", snatRule.Status)
	d.Set("description", snatRule.Description)

	return nil
}

func resourceNatSnatRuleV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	natClient, err := config.NatGatewayClient(region)
	if err != nil {
		return fmtp.Errorf("error creating nat client: %s", err)
	}

	ruleID := d.Id()
	updateOpts := &hw_snatrules.UpdateOpts{
		NatGatewayID: d.Get("nat_gateway_id").(string),
	}
	if d.HasChange("description") {
		desc := d.Get("description").(string)
		updateOpts.Description = &desc
	}
	if d.HasChange("floating_ip_id") {
		eipClient, err := config.NetworkingV1Client(region)
		if err != nil {
			return fmtp.Errorf("error creating networking client: %s", err)
		}

		eipIDs := d.Get("floating_ip_id").(string)
		eipList := strings.Split(eipIDs, ",")
		eipAddrs := make([]string, len(eipList))

		// get EIP address from ID
		for i, id := range eipList {
			eIP, err := eips.Get(eipClient, id).Extract()
			if err != nil {
				return fmtp.Errorf("error fetching EIP %s: %s", id, err)
			}
			eipAddrs[i] = eIP.PublicAddress
		}

		updateOpts.FloatingIPAddress = strings.Join(eipAddrs, ",")
	}

	logp.Printf("[DEBUG] update Options: %#v", updateOpts)
	_, err = hw_snatrules.Update(natClient, ruleID, updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("error updating SNAT rule: %s", err)
	}

	logp.Printf("[DEBUG] waiting for SNAT rule (%s) to become available", ruleID)
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Refresh:      waitForSnatRuleActive(natClient, ruleID),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        3 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("error updating SNAT rule: %s", err)
	}

	return resourceNatSnatRuleV2Read(d, meta)
}

func resourceNatSnatRuleV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	natClient, err := config.NatGatewayClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	natGatewayID := d.Get("nat_gateway_id").(string)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSnatRuleDelete(natClient, d.Id(), natGatewayID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Snat Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSnatRuleActive(client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := hw_snatrules.Get(client, nId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, "", nil
	}
}

func waitForSnatRuleDelete(client *golangsdk.ServiceClient, nId, natGatewayID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Snat Rule %s.\n", nId)

		n, err := hw_snatrules.Get(client, nId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Snat Rule %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = hw_snatrules.Delete(client, nId, natGatewayID).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Snat Rule %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Snat Rule %s still active.\n", nId)
		return n, "ACTIVE", nil
	}
}
