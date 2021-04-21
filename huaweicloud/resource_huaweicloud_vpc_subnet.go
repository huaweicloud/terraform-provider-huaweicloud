package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// refer to: https://support.huaweicloud.com/intl/en-us/dns_faq/dns_faq_002.html
var defaultDNSList = map[string][]string{
	"cn-north-1":     {"100.125.1.250", "100.125.21.250"},  // Beijing-1
	"cn-north-4":     {"100.125.1.250", "100.125.129.250"}, // Beijing-4
	"cn-east-2":      {"100.125.17.29", "100.125.135.29"},  // Shanghai-2
	"cn-east-3":      {"100.125.1.250", "100.125.64.250"},  // Shanghai-1
	"cn-south-1":     {"100.125.1.250", "100.125.136.29"},  // Guangzhou
	"cn-southwest-2": {"100.125.1.250", "100.125.129.250"}, // Guiyang-1
	"ap-southeast-1": {"100.125.1.250", "100.125.3.250"},   // Hong Kong
	"ap-southeast-2": {"100.125.1.250", "182.50.80.4"},     // Bangkok
	"ap-southeast-3": {"100.125.1.250", "100.125.128.250"}, // Singapore
	"af-south-1":     {"100.125.1.250", "8.8.8.8"},         // Johannesburg
	"sa-brazil-1":    {"100.125.1.22", "8.8.8.8"},          // Sao Paulo-1
	"na-mexico-1":    {"100.125.1.22", "8.8.8.8"},          // Mexico City-1
	"la-south-2":     {"100.125.1.250", "8.8.8.8"},         // Santiago
	"sa-chile-1":     {"100.125.1.250", "100.125.0.250"},   // LA-Santiago2
}

func resourceSubnetDNSListV1(d *schema.ResourceData, region string) []string {
	rawDNSN := d.Get("dns_list").([]interface{})
	dnsn := make([]string, len(rawDNSN))
	for i, raw := range rawDNSN {
		dnsn[i] = raw.(string)
	}

	// get the default DNS if it was not specified in schema
	_, hasPrimaryDNS := d.GetOk("primary_dns")
	if region != "" && len(rawDNSN) == 0 && !hasPrimaryDNS {
		dnsn = defaultDNSList[region]
	}
	return dnsn
}

func ResourceVpcSubnetV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcSubnetV1Create,
		Read:   resourceVpcSubnetV1Read,
		Update: resourceVpcSubnetV1Update,
		Delete: resourceVpcSubnetV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateString64WithChinese,
			},
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.ValidateCIDR,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.ValidateIP,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dhcp_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"primary_dns": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateIP,
				Computed:     true,
			},
			"secondary_dns": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateIP,
				RequiredWith: []string{"primary_dns"},
				Computed:     true,
			},
			"dns_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: utils.ValidateIP,
				},
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceVpcSubnetV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	subnetClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud networking client: %s", err)
	}

	enable := d.Get("ipv6_enable").(bool)
	createOpts := subnets.CreateOpts{
		Name:             d.Get("name").(string),
		CIDR:             d.Get("cidr").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		GatewayIP:        d.Get("gateway_ip").(string),
		EnableIPv6:       &enable,
		EnableDHCP:       d.Get("dhcp_enable").(bool),
		VPC_ID:           d.Get("vpc_id").(string),
		PRIMARY_DNS:      d.Get("primary_dns").(string),
		SECONDARY_DNS:    d.Get("secondary_dns").(string),
		DnsList:          resourceSubnetDNSListV1(d, region),
	}
	log.Printf("[DEBUG] Create VPC subnet options: %#v", createOpts)

	n, err := subnets.Create(subnetClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud VPC subnet: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[INFO] Vpc Subnet ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"UNKNOWN"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcSubnetActive(subnetClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for Subnet (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		vpcSubnetV2Client, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud VpcSubnet client: %s", err)
		}
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcSubnetV2Client, "subnets", n.ID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of VpcSubnet %q: %s", n.ID, tagErr)
		}
	}

	return resourceVpcSubnetV1Read(d, config)

}

func resourceVpcSubnetV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	subnetClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud networking client: %s", err)
	}

	n, err := subnets.Get(subnetClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Huaweicloud Subnets: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("cidr", n.CIDR)
	d.Set("dns_list", n.DnsList)
	d.Set("gateway_ip", n.GatewayIP)
	d.Set("ipv6_enable", n.EnableIPv6)
	d.Set("dhcp_enable", n.EnableDHCP)
	d.Set("primary_dns", n.PRIMARY_DNS)
	d.Set("secondary_dns", n.SECONDARY_DNS)
	d.Set("availability_zone", n.AvailabilityZone)
	d.Set("vpc_id", n.VPC_ID)
	d.Set("subnet_id", n.SubnetId)
	d.Set("ipv6_subnet_id", n.IPv6SubnetId)
	d.Set("ipv6_cidr", n.IPv6CIDR)
	d.Set("ipv6_gateway", n.IPv6Gateway)
	d.Set("region", GetRegion(d, config))

	// save VpcSubnet tags
	if vpcSubnetV2Client, err := config.NetworkingV2Client(GetRegion(d, config)); err == nil {
		if resourceTags, err := tags.Get(vpcSubnetV2Client, "subnets", d.Id()).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			if err := d.Set("tags", tagmap); err != nil {
				return fmt.Errorf("Error saving tags to state for Subnet (%s): %s", d.Id(), err)
			}
		} else {
			log.Printf("[WARN] Error fetching tags of Subnet (%s): %s", d.Id(), err)
		}
	} else {
		return fmt.Errorf("Error creating VpcSubnet client: %s", err)
	}

	return nil
}

func resourceVpcSubnetV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	subnetClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud networking client: %s", err)
	}

	var updateOpts subnets.UpdateOpts

	//as name is mandatory while updating subnet
	updateOpts.Name = d.Get("name").(string)

	if d.HasChange("ipv6_enable") {
		if d.Get("ipv6_enable").(bool) {
			enable := d.Get("ipv6_enable").(bool)
			updateOpts.EnableIPv6 = &enable
		} else {
			return fmt.Errorf("Parameter cannot be disabled after IPv6 enable")
		}
	}
	if d.HasChange("primary_dns") {
		updateOpts.PRIMARY_DNS = d.Get("primary_dns").(string)
	}
	if d.HasChange("secondary_dns") {
		updateOpts.SECONDARY_DNS = d.Get("secondary_dns").(string)
	}
	if d.HasChange("dns_list") {
		updateOpts.DnsList = resourceSubnetDNSListV1(d, "")
	}
	if d.HasChange("dhcp_enable") {
		updateOpts.EnableDHCP = d.Get("dhcp_enable").(bool)

	} else if d.Get("dhcp_enable").(bool) { //maintaining dhcp to be true if it was true earlier as default update option for dhcp bool is always going to be false in golangsdk
		updateOpts.EnableDHCP = true
	}

	vpc_id := d.Get("vpc_id").(string)

	_, err = subnets.Update(subnetClient, vpc_id, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating Huaweicloud VPC Subnet: %s", err)
	}

	//update tags
	if d.HasChange("tags") {
		vpcSubnetV2Client, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud VpcSubnet client: %s", err)
		}

		tagErr := utils.UpdateResourceTags(vpcSubnetV2Client, d, "subnets", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of VPC subnet %s: %s", d.Id(), tagErr)
		}
	}

	return resourceVpcSubnetV1Read(d, meta)
}

func resourceVpcSubnetV1Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*config.Config)
	subnetClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud networking client: %s", err)
	}
	vpc_id := d.Get("vpc_id").(string)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcSubnetDelete(subnetClient, vpc_id, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting Huaweicloud Subnet: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcSubnetActive(subnetClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := subnets.Get(subnetClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		//If subnet status is other than Active, send error
		if n.Status == "DOWN" || n.Status == "ERROR" {
			return nil, "", fmt.Errorf("Subnet status: '%s'", n.Status)
		}

		return n, "UNKNOWN", nil
	}
}

func waitForVpcSubnetDelete(subnetClient *golangsdk.ServiceClient, vpcId string, subnetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := subnets.Get(subnetClient, subnetId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Huaweicloud subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault500); ok {
				log.Printf("[DEBUG] Got 500 error when delting HuaweiCloud subnet %s, it should be stream control on API server, try again later", subnetId)
				return r, "ACTIVE", nil
			}
			return r, "ACTIVE", err
		}
		err = subnets.Delete(subnetClient, vpcId, subnetId).ExtractErr()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Huaweicloud subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				log.Printf("[INFO] Successfully deleted Huaweicloud subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault500); ok {
				log.Printf("[DEBUG] Got 500 error when delting HuaweiCloud subnet %s, it should be stream control on API server, try again later", subnetId)
				return r, "ACTIVE", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
