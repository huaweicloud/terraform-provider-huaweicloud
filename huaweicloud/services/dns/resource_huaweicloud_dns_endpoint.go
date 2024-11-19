package dns

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dns/v2/endpoints"
	"github.com/chnsz/golangsdk/openstack/dns/v2/ipaddress"
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS DELETE /v2.1/endpoints/{endpoint_id}/ipaddresses/{ipaddress_id}
// @API DNS GET /v2.1/endpoints/{endpoint_id}/ipaddresses
// @API DNS POST /v2.1/endpoints/{endpoint_id}/ipaddresses
// @API DNS DELETE /v2.1/endpoints/{endpoint_id}
// @API DNS GET /v2.1/endpoints/{endpoint_id}
// @API DNS PUT /v2.1/endpoints/{endpoint_id}
// @API DNS POST /v2.1/endpoints
// @API VPC GET /v1/{project_id}/subnets
func ResourceDNSEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSEndpointCreate,
		UpdateContext: resourceDNSEndpointUpdate,
		ReadContext:   resourceDNSEndpointRead,
		DeleteContext: resourceDNSEndpointDelete,
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
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// Multiple IP addresses can be assigned to a subnet, so the Set type cannot be used.
			"ip_addresses": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 6,
				MinItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_address_id": {
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
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resolver_rule_count": {
				Type:     schema.TypeInt,
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
		},
	}
}

func resourceDNSEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	dnsClient, err := conf.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}
	name := d.Get("name").(string)
	direction := d.Get("direction").(string)
	ipAddresses := d.Get("ip_addresses").([]interface{})
	ipAddressesList := make([]endpoints.IPAddresses, len(ipAddresses))
	for i, ipObject := range ipAddresses {
		ip := ipObject.(map[string]interface{})
		ipAddressesList[i] = endpoints.IPAddresses{
			SubnetID: ip["subnet_id"].(string),
			IP:       ip["ip"].(string),
		}
	}
	opt := endpoints.CreateOpt{
		Name:        name,
		Direction:   direction,
		Region:      region,
		IPAddresses: ipAddressesList,
	}

	endpoint, err := endpoints.Create(dnsClient, opt).Extract()
	if err != nil {
		return diag.Errorf("err creating DNS endpoint: %s", err)
	}
	endpointId := endpoint.ID
	d.SetId(endpointId)
	log.Printf("[DEBUG] Waiting for DNS endpoint (%s) to become available", endpointId)
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING"},
		Refresh:      refreshEndpointStatus(dnsClient, endpointId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for DNS endpoint (%s) to become ACTIVE for creation: %s",
			endpointId, err)
	}

	return resourceDNSEndpointRead(ctx, d, meta)
}

func refreshEndpointStatus(client *golangsdk.ServiceClient, endpointId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		endpoint, err := endpoints.Get(client, endpointId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return endpoint, "DELETED", nil
			}
			return nil, "", err
		}

		log.Printf("[DEBUG] DNS endpoint (%s) current status: %s", endpoint.ID, endpoint.Status)
		return endpoint, parseStatus(endpoint.Status), nil
	}
}

func resourceDNSEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	dnsClient, err := conf.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	name := d.Get("name").(string)

	if d.HasChange("name") {
		opts := endpoints.UpdateOpts{Name: name}
		_, err = endpoints.Update(dnsClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("ip_addresses") {
		err = updateIpAddresses(dnsClient, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] Waiting for DNS endpoint (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING"},
		Refresh:      refreshEndpointStatus(dnsClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for DNS endpoint (%s) to become ACTIVE for update: %s",
			d.Id(), err)
	}

	return resourceDNSEndpointRead(ctx, d, meta)
}

func updateIpAddresses(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		oldRaws, newRaws = d.GetChange("ip_addresses")
		num              = len(oldRaws.([]interface{}))
		endpointId       = d.Id()
	)
	log.Printf("[DEBUG] update IP Address oldRaws (%s), newRaws (%s)", oldRaws, newRaws)
	addRaws := getDiffIpAddresses(newRaws.([]interface{}), oldRaws.([]interface{}))
	log.Printf("[DEBUG] update IP Address addRaws (%s)", addRaws)
	removeRaws := getDiffIpAddresses(oldRaws.([]interface{}), newRaws.([]interface{}))
	log.Printf("[DEBUG] update IP Address removeRaws (%s)", removeRaws)

	// The `ip_addresses` parameter length limit range form `2` to `6`.
	// The following logic is to handle the critical value scenario.
	// For example, the original number is `3` IP addresses, `4` IP addresses are added, and `1` IP address is deleted.
	var err error
	for {
		if len(addRaws) == 0 && len(removeRaws) == 0 {
			return nil
		}
		if num < 6 && len(addRaws) > 0 {
			addRaws, err = addIpAddress(client, addRaws, endpointId)
			if err != nil {
				return err
			}
			num++
			continue
		}
		removeRaws, err = removeIpAddress(client, removeRaws, endpointId)
		if err != nil {
			return err
		}
		num--
	}
}

// The getDiffIpAddresses method is used to get the part in "from" array but not in "to" array.
func getDiffIpAddresses(from, to []interface{}) []interface{} {
	rest := make([]interface{}, 0)
	for i, fRaw := range from {
		for j, tRaw := range to {
			isMatched := utils.PathSearch("is_matched", tRaw, false).(bool)
			if !isMatched && utils.PathSearch("ip", fRaw, "").(string) == utils.PathSearch("ip", tRaw, "").(string) {
				if utils.PathSearch("subnet_id", fRaw, "").(string) == utils.PathSearch("subnet_id", tRaw, "").(string) {
					from[i].(map[string]interface{})["is_matched"] = true
					to[j].(map[string]interface{})["is_matched"] = true
				} else {
					// When the IP is not specified in the new value, the IP of the old value will be obtained,
					// which will cause the IP to not correspond to the subnet, so special processing to set the IP to empty.
					from[i].(map[string]interface{})["ip"] = ""
				}
			}
		}

		if !utils.PathSearch("is_matched", fRaw, false).(bool) {
			rest = append(rest, fRaw)
		}
	}

	return rest
}

func addIpAddress(client *golangsdk.ServiceClient, list []interface{}, endpointId string) ([]interface{}, error) {
	opts := ipaddress.CreateOpts{IPAddress: ipaddress.IPAddress{
		SubnetID: utils.PathSearch("subnet_id", list[0], "").(string),
		IP:       utils.PathSearch("ip", list[0], "").(string),
	}}
	_, err := ipaddress.Create(client, opts, endpointId).Extract()
	if err != nil {
		return nil, err
	}
	return list[1:], nil
}

func removeIpAddress(client *golangsdk.ServiceClient, list []interface{}, endpointId string) ([]interface{}, error) {
	ipAddressId := utils.PathSearch("ip_address_id", list[0], "").(string)
	err := ipaddress.Delete(client, endpointId, ipAddressId).ExtractErr()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return list[1:], nil
		}
		return nil, err
	}
	return list[1:], nil
}

func resourceDNSEndpointRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	dnsClient, err := conf.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	id := d.Id()
	endpointInfo, err := endpoints.Get(dnsClient, id).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS endpoint")
	}
	log.Printf("[DEBUG] Retrieved endpoint %s: %#v", id, endpointInfo)

	ipAddress, err := ipaddress.List(dnsClient, id).Extract()
	if err != nil {
		return diag.Errorf("error retrieving IP addresses: %s", err)
	}
	log.Printf("[DEBUG] Retrieved IP addresses %s: %#v", id, ipAddress)

	subnetClient, err := conf.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	subnetList, err := subnets.List(subnetClient, subnets.ListOpts{VPC_ID: endpointInfo.VpcID})
	if err != nil {
		return diag.Errorf("unable to retrieve subnets: %s", err)
	}

	if len(subnetList) == 0 {
		return diag.Errorf("unable to find subnet from API response")
	}

	log.Printf("[DEBUG] Retrieved ip address %s: %#v", id, ipAddress)

	mErr := multierror.Append(nil,
		d.Set("name", endpointInfo.Name),
		d.Set("direction", endpointInfo.Direction),
		d.Set("status", endpointInfo.Status),
		d.Set("vpc_id", endpointInfo.VpcID),
		d.Set("resolver_rule_count", endpointInfo.ResolverRuleCount),
		d.Set("created_at", endpointInfo.CreateTime),
		d.Set("updated_at", endpointInfo.UpdateTime),
		d.Set("ip_addresses", flattenIpAddresses(ipAddress, subnetList)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}

	return nil
}

func flattenIpAddresses(ipAddresses []ipaddress.ListObject, subnetList []subnets.Subnet) []interface{} {
	var rest []interface{}
	for _, ipObject := range ipAddresses {
		ipAddress := map[string]interface{}{
			"ip_address_id": ipObject.ID,
			"ip":            ipObject.IP,
			"status":        ipObject.Status,
			"created_at":    ipObject.CreateTime,
			"updated_at":    ipObject.UpdateTime,
		}
		for _, subnet := range subnetList {
			if subnet.SubnetId == ipObject.SubnetID {
				ipAddress["subnet_id"] = subnet.ID
				break
			}
		}
		rest = append(rest, ipAddress)
	}
	return rest
}

func resourceDNSEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	dnsClient, err := conf.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	err = endpoints.Delete(dnsClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting DNS endpoint: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS endpoint (%s) to become DELETED", d.Id())
	stateConf := &resource.StateChangeConf{
		Target: []string{"DELETED"},
		// we allow to try to delete ERROR endpoint
		Pending:      []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:      refreshEndpointStatus(dnsClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for DNS endpoint (%s) to delete: %s",
			d.Id(), err)
	}

	return nil
}
