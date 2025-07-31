package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dns/v2/nameservers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v1/{project_id}/subnets
// @API VPC GET /v1/{project_id}/subnets/{id}
// @API VPC PUT /v1/{project_id}/vpcs/{vpcid}/subnets/{id}
// @API VPC DELETE /v1/{project_id}/vpcs/{vpcid}/subnets/{id}
// @API VPC POST /v2.0/{project_id}/subnets/{id}/tags/action
// @API VPC GET /v2.0/{project_id}/subnets/{id}/tags
// @API DNS GET /v2/nameservers

// refer to: https://support.huaweicloud.com/intl/en-us/dns_faq/dns_faq_002.html
var privateDNSList = map[string][]string{
	"cn-north-1":     {"100.125.1.250", "100.125.21.250"},  // Beijing-1
	"cn-north-4":     {"100.125.1.250", "100.125.129.250"}, // Beijing-4
	"cn-north-9":     {"100.125.1.250", "100.125.107.250"}, // Ulanqab
	"cn-east-2":      {"100.125.17.29", "100.125.135.29"},  // Shanghai-2
	"cn-east-3":      {"100.125.1.250", "100.125.64.250"},  // Shanghai-1
	"cn-south-1":     {"100.125.1.250", "100.125.136.29"},  // Guangzhou
	"cn-south-4":     {"100.125.0.167"},                    // Guangzhou-InvitationOnly
	"cn-southwest-2": {"100.125.1.250", "100.125.129.250"}, // Guiyang-1
	"ap-southeast-1": {"100.125.1.250", "100.125.3.250"},   // Hong Kong
	"ap-southeast-2": {"100.125.1.250", "100.125.1.251"},   // Bangkok
	"ap-southeast-3": {"100.125.1.250", "100.125.128.250"}, // Singapore
	"af-south-1":     {"100.125.1.250", "100.125.1.14"},    // Johannesburg
	"tr-west-1":      {"100.125.2.250", "100.125.2.251"},   // Turkey Istanbul
	"sa-brazil-1":    {"100.125.1.22", "100.125.1.90"},     // LA-Sao Paulo-1
	"na-mexico-1":    {"100.125.1.22", "100.125.1.90"},     // LA-Mexico City-1
	"la-north-2":     {"100.125.1.250", "100.125.1.242"},   // LA-Mexico City-2
	"la-south-2":     {"100.125.1.250", "100.125.0.250"},   // LA-Santiago
	"sa-chile-1":     {"100.125.1.250", "100.125.0.250"},   // LA-Santiago2
}

// buildSubnetDNSList is used to obtain the DNS list according to the region in the following orders:
// 1. return the dns_list if it was specified;
// 2. return nil if primary_dns was specified, otherwise continue;
// 3. return the private DNS list if the region is in **privateDNSList**
// 4. try to get the private DNS list through API https://support.huaweicloud.com/intl/en-us/api-dns/dns_api_69001.html
// 5. return the public DNS if any error occurs;
func buildSubnetDNSList(d *schema.ResourceData, cfg *config.Config, region string) []string {
	if raw, ok := d.GetOk("dns_list"); ok {
		return utils.ExpandToStringList(raw.([]interface{}))
	}

	// get the DNS list only when primary_dns was not specified
	_, hasPrimaryDNS := d.GetOk("primary_dns")
	if hasPrimaryDNS {
		return nil
	}

	if dnsn, ok := privateDNSList[region]; ok {
		return dnsn
	}

	// public DNS: 8.8.8.8(google-public-dns-a.google.com) and 114.114.114.114(China)
	publicDNSList := []string{"8.8.8.8", "114.114.114.114"}

	dnsClient, err := cfg.DnsWithRegionClient(region)
	if err != nil {
		log.Printf("[WARN] cannot generate DNS client, use %v as the DNS list", publicDNSList)
		return publicDNSList
	}

	opts := nameservers.ListOpts{
		Region: region,
	}
	dsnServers, err := nameservers.List(dnsClient, &opts)
	if len(dsnServers) > 0 {
		records := dsnServers[0].Records
		dnsList := make([]string, len(records))
		for i := range records {
			dnsList[i] = records[i].Address
		}
		return dnsList
	}

	if err != nil {
		log.Printf("[WARN] failed to fetch the name servers: %s", err)
	}
	log.Printf("[WARN] use %v as the DNS list", publicDNSList)
	return publicDNSList
}

func ResourceVpcSubnetV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcSubnetCreate,
		ReadContext:   resourceVpcSubnetRead,
		UpdateContext: resourceVpcSubnetUpdate,
		DeleteContext: resourceVpcSubnetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{ // request and response parameters
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
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.ValidateCIDR,
			},
			"gateway_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.ValidateIP,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"ntp_server_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp_lease_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_ipv6_lease_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "schema: Deprecated",
			},
			"ipv4_subnet_id": {
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
			"tags": common.TagsSchema(),
		},
	}
}

func buildDhcpOpts(d *schema.ResourceData, update bool) []subnets.ExtraDhcpOpt {
	var result []subnets.ExtraDhcpOpt
	if v, ok := d.GetOk("dhcp_lease_time"); ok {
		addressVal := v.(string)
		addressTime := subnets.ExtraDhcpOpt{
			OptName:  "addresstime",
			OptValue: &addressVal,
		}
		result = append(result, addressTime)
	}

	if v, ok := d.GetOk("dhcp_ipv6_lease_time"); ok {
		ipv6AddressVal := v.(string)
		ipv6AddressTime := subnets.ExtraDhcpOpt{
			OptName:  "ipv6_addresstime",
			OptValue: &ipv6AddressVal,
		}
		result = append(result, ipv6AddressTime)
	}

	if v, ok := d.GetOk("ntp_server_address"); ok {
		ntpVal := v.(string)
		ntp := subnets.ExtraDhcpOpt{
			OptName:  "ntp",
			OptValue: &ntpVal,
		}
		result = append(result, ntp)
	} else if update {
		ntp := subnets.ExtraDhcpOpt{
			OptName: "ntp",
		}
		result = append(result, ntp)
	}

	if v, ok := d.GetOk("dhcp_domain_name"); ok {
		domainNameVal := v.(string)
		domainName := subnets.ExtraDhcpOpt{
			OptName:  "domainname",
			OptValue: &domainNameVal,
		}
		result = append(result, domainName)
	} else if update {
		domainName := subnets.ExtraDhcpOpt{
			OptName: "domainname",
		}
		result = append(result, domainName)
	}

	return result
}

func resourceVpcSubnetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	subnetClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	enable := d.Get("ipv6_enable").(bool)
	createOpts := subnets.CreateOpts{
		Name:             d.Get("name").(string),
		CIDR:             d.Get("cidr").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		GatewayIP:        d.Get("gateway_ip").(string),
		Description:      d.Get("description").(string),
		EnableIPv6:       &enable,
		EnableDHCP:       d.Get("dhcp_enable").(bool),
		VPC_ID:           d.Get("vpc_id").(string),
		PRIMARY_DNS:      d.Get("primary_dns").(string),
		SECONDARY_DNS:    d.Get("secondary_dns").(string),
		DnsList:          buildSubnetDNSList(d, config, region),
		ExtraDhcpOpts:    buildDhcpOpts(d, false),
	}
	log.Printf("[DEBUG] Create VPC subnet options: %#v", createOpts)

	n, err := subnets.Create(subnetClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC subnet: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[INFO] Vpc Subnet ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"UNKNOWN"},
		Target:       []string{"ACTIVE"},
		Refresh:      waitForVpcSubnetActive(subnetClient, n.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"Error waiting for Subnet (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		vpcSubnetV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating VpcSubnet client: %s", err)
		}
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcSubnetV2Client, "subnets", n.ID, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of VpcSubnet %q: %s", n.ID, tagErr)
		}
	}

	return resourceVpcSubnetRead(ctx, d, config)

}

// GetVpcSubnetById is a method to obtain subnet informations from special region through subnet ID.
func GetVpcSubnetById(config *config.Config, region, networkId string) (*subnets.Subnet, error) {
	subnetClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return nil, err
	}

	return subnets.Get(subnetClient, networkId).Extract()
}

func resourceVpcSubnetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	n, err := GetVpcSubnetById(config, region, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error obtain Subnet information")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("description", n.Description),
		d.Set("cidr", n.CIDR),
		d.Set("dns_list", n.DnsList),
		d.Set("gateway_ip", n.GatewayIP),
		d.Set("ipv6_enable", n.EnableIPv6),
		d.Set("dhcp_enable", n.EnableDHCP),
		d.Set("primary_dns", n.PRIMARY_DNS),
		d.Set("secondary_dns", n.SECONDARY_DNS),
		d.Set("availability_zone", n.AvailabilityZone),
		d.Set("vpc_id", n.VPC_ID),
		d.Set("subnet_id", n.SubnetId),
		d.Set("ipv4_subnet_id", n.SubnetId),
		d.Set("ipv6_subnet_id", n.IPv6SubnetId),
		d.Set("ipv6_cidr", n.IPv6CIDR),
		d.Set("ipv6_gateway", n.IPv6Gateway),
	)

	// save VpcSubnet tags
	if vpcSubnetV2Client, err := config.NetworkingV2Client(region); err == nil {
		if resourceTags, err := tags.Get(vpcSubnetV2Client, "subnets", d.Id()).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			mErr = multierror.Append(mErr, d.Set("tags", tagmap))
		} else {
			log.Printf("[WARN] Error fetching tags of Subnet (%s): %s", d.Id(), err)
		}
	} else {
		return diag.Errorf("error creating VpcSubnet client: %s", err)
	}

	// set dhcp extra opts ntp, addresstime, ipv6_addresstime, domainname
	for _, val := range n.ExtraDhcpOpts {
		switch val.OptName {
		case "ntp":
			mErr = multierror.Append(mErr, d.Set("ntp_server_address", val.OptValue))
		case "addresstime":
			mErr = multierror.Append(mErr, d.Set("dhcp_lease_time", val.OptValue))
		case "ipv6_addresstime":
			mErr = multierror.Append(mErr, d.Set("dhcp_ipv6_lease_time", val.OptValue))
		case "domainname":
			mErr = multierror.Append(mErr, d.Set("dhcp_domain_name", val.OptValue))
		}
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting VPC subnet fields: %s", err)
	}

	return nil
}

func resourceVpcSubnetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	subnetClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	if d.HasChanges("name", "description", "dhcp_enable", "primary_dns", "secondary_dns", "dns_list",
		"ipv6_enable", "dhcp_lease_time", "ntp_server_address", "dhcp_ipv6_lease_time", "dhcp_domain_name") {
		var updateOpts subnets.UpdateOpts

		// name is mandatory while updating subnet
		updateOpts.Name = d.Get("name").(string)
		// always setting dhcp in updateOpts as the field defauts to be false in golangsdk
		updateOpts.EnableDHCP = d.Get("dhcp_enable").(bool)

		if d.HasChange("ipv6_enable") {
			if d.Get("ipv6_enable").(bool) {
				enable := d.Get("ipv6_enable").(bool)
				updateOpts.EnableIPv6 = &enable
			} else {
				return diag.Errorf("parameter cannot be disabled after IPv6 enable")
			}
		}
		if d.HasChange("description") {
			description := d.Get("description").(string)
			updateOpts.Description = &description
		}
		if d.HasChange("primary_dns") {
			updateOpts.PRIMARY_DNS = d.Get("primary_dns").(string)
		}
		if d.HasChange("secondary_dns") {
			updateOpts.SECONDARY_DNS = d.Get("secondary_dns").(string)
		}
		if d.HasChange("dns_list") {
			dnsList := utils.ExpandToStringList(d.Get("dns_list").([]interface{}))
			updateOpts.DnsList = &dnsList
		}
		if d.HasChanges("dhcp_lease_time", "ntp_server_address", "dhcp_ipv6_lease_time", "dhcp_domain_name") {
			updateOpts.ExtraDhcpOpts = buildDhcpOpts(d, true)
		}

		log.Printf("[DEBUG] Update VPC subnet options: %#v", updateOpts)
		vpcID := d.Get("vpc_id").(string)
		_, err = subnets.Update(subnetClient, vpcID, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating VPC Subnet: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		vpcSubnetV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating VpcSubnet client: %s", err)
		}

		tagErr := utils.UpdateResourceTags(vpcSubnetV2Client, d, "subnets", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPC subnet %s: %s", d.Id(), tagErr)
		}
	}

	return resourceVpcSubnetRead(ctx, d, meta)
}

func resourceVpcSubnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*config.Config)
	subnetClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	vpcID := d.Get("vpc_id").(string)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE"},
		Target:       []string{"DELETED"},
		Refresh:      waitForVpcSubnetDelete(subnetClient, vpcID, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting Subnet: %s", err)
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

		// If subnet status is other than Active, send error
		if n.Status == "DOWN" || n.Status == "ERROR" {
			return nil, "", fmt.Errorf("subnet status: '%s'", n.Status)
		}

		return n, "UNKNOWN", nil
	}
}

func waitForVpcSubnetDelete(subnetClient *golangsdk.ServiceClient, vpcId string, subnetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := subnets.Get(subnetClient, subnetId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault500); ok {
				log.Printf("[DEBUG] Got 500 error when delting subnet %s, it should be stream control on API server, try again later", subnetId)
				return r, "ACTIVE", nil
			}

			// due to api problem, under enterprise project permissions, when the subnet is deleted,
			// there will be a temporary 403 error, try again later will get 404 error which is excepted,
			// add this in retry can avoid temporary 403 error
			// but the real 403 error can only be returned after the timeout period is reached
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				log.Printf("[DEBUG] Got 403 error when delting subnet %s, it should be temporary permission error, try again later", subnetId)
				return r, "ACTIVE", nil
			}
			return r, "ACTIVE", err
		}
		err = subnets.Delete(subnetClient, vpcId, subnetId).ExtractErr()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				log.Printf("[INFO] Successfully deleted subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault500); ok {
				log.Printf("[DEBUG] Got 500 error when delting subnet %s, it should be stream control on API server, try again later", subnetId)
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
