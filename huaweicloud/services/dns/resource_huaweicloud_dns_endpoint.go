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
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the endpoint.`,
			},
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The direction of the endpoint.`,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 6,
				MinItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The subnet ID of the IP address.`,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								"The IP address associated with the endpoint.",
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"ip_address_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the IP address.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of IP address.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the IP address.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the endpoint.`,
						},
					},
				},
				Description: `The list of the IP addresses of the endpoint.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of endpoint.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC to which the subnet belongs.`,
			},
			"resolver_rule_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of bound resolver rules.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the endpoint.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the endpoint.`,
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
	d.SetId(endpoint.ID)
	log.Printf("[DEBUG] Waiting for DNS endpoint (%s) to become available", endpoint.ID)
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING"},
		Refresh:      waitForDNSEndpoint(dnsClient, endpoint.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for DNS endpoint (%s) to become ACTIVE for creation: %s",
			endpoint.ID, err)
	}

	return resourceDNSEndpointRead(ctx, d, meta)
}

func waitForDNSEndpoint(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		endpoint, err := endpoints.Get(client, id).Extract()
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
		Refresh:      waitForDNSEndpoint(dnsClient, d.Id()),
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
		endpointID       = d.Id()
	)
	log.Printf("[DEBUG] update IP Address oldRaws (%s), newRaws (%s)", oldRaws, newRaws)
	addRaws := getDisjointPart(newRaws.([]interface{}), oldRaws.([]interface{}))
	log.Printf("[DEBUG] update IP Address addRaws (%s)", addRaws)
	removeRaws := getDisjointPart(oldRaws.([]interface{}), newRaws.([]interface{}))
	log.Printf("[DEBUG] update IP Address removeRaws (%s)", removeRaws)
	var err error
	for {
		if len(addRaws) == 0 && len(removeRaws) == 0 {
			return nil
		}
		if num < 6 && len(addRaws) > 0 {
			addRaws, err = addIpAddress(client, addRaws, endpointID)
			if err != nil {
				return err
			}
			num++
			continue
		}
		removeRaws, err = removeIpAddress(client, removeRaws, endpointID)
		if err != nil {
			return err
		}
		num--
	}
}

// getDisjointPart get the part in "from" array but not in "to" array
func getDisjointPart(from []interface{}, to []interface{}) []ipaddress.ListObject {
	type IPMsg struct {
		IPAddressID string
		SubnetID    string
		IP          string
		IsMatched   bool
	}

	var fIPList = make([]IPMsg, len(from))
	for i, f := range from {
		fMap := f.(map[string]interface{})
		fIPList[i] = IPMsg{
			IPAddressID: fMap["ip_address_id"].(string),
			SubnetID:    fMap["subnet_id"].(string),
			IP:          fMap["ip"].(string),
		}
	}

	var tIPList = make([]IPMsg, len(to))
	for i, t := range to {
		tMap := t.(map[string]interface{})
		tIPList[i] = IPMsg{
			IPAddressID: tMap["ip_address_id"].(string),
			SubnetID:    tMap["subnet_id"].(string),
			IP:          tMap["ip"].(string),
		}
	}

	for i, fIP := range fIPList {
		for j, tIP := range tIPList {
			if !tIP.IsMatched && fIP.IP == tIP.IP {
				if fIP.SubnetID == tIP.SubnetID {
					fIPList[i].IsMatched = true
					tIPList[j].IsMatched = true
				} else {
					fIPList[i].IP = ""
				}
			}
		}
	}

	var res []ipaddress.ListObject
	for _, ip := range fIPList {
		if !ip.IsMatched {
			res = append(res, ipaddress.ListObject{
				ID:       ip.IPAddressID,
				SubnetID: ip.SubnetID,
				IP:       ip.IP,
			})
		}
	}
	return res
}

func addIpAddress(client *golangsdk.ServiceClient, list []ipaddress.ListObject, endpointID string) ([]ipaddress.ListObject, error) {
	opts := ipaddress.CreateOpts{IPAddress: ipaddress.IPAddress{
		SubnetID: list[0].SubnetID,
		IP:       list[0].IP,
	}}
	_, err := ipaddress.Create(client, opts, endpointID).Extract()
	if err != nil {
		return nil, err
	}
	return list[1:], nil
}

func removeIpAddress(client *golangsdk.ServiceClient, list []ipaddress.ListObject, endpointID string) ([]ipaddress.ListObject, error) {
	err := ipaddress.Delete(client, endpointID, list[0].ID).ExtractErr()
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
		return diag.Errorf("error retrieving ip addresses: %s", err)
	}
	log.Printf("[DEBUG] Retrieved ip address %s: %#v", id, ipAddress)

	subnetClient, err := conf.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	subnetList, err := subnets.List(subnetClient, subnets.ListOpts{VPC_ID: endpointInfo.VpcID})
	if err != nil {
		return diag.Errorf("unable to retrieve subnets: %s", err)
	}

	if len(subnetList) == 0 {
		return diag.Errorf("no subnet found")
	}

	log.Printf("[DEBUG] Retrieved ip address %s: %#v", id, ipAddress)
	var ipAddresses []interface{}
	for _, ipObject := range ipAddress {
		ipAddress := make(map[string]interface{}, 6)
		ipAddress["ip_address_id"] = ipObject.ID
		ipAddress["ip"] = ipObject.IP
		ipAddress["status"] = ipObject.Status
		ipAddress["created_at"] = ipObject.CreateTime
		ipAddress["updated_at"] = ipObject.UpdateTime
		for _, subnet := range subnetList {
			if subnet.SubnetId == ipObject.SubnetID {
				ipAddress["subnet_id"] = subnet.ID
				break
			}
		}
		ipAddresses = append(ipAddresses, ipAddress)
	}
	mErr := multierror.Append(nil,
		d.Set("name", endpointInfo.Name),
		d.Set("direction", endpointInfo.Direction),
		d.Set("status", endpointInfo.Status),
		d.Set("vpc_id", endpointInfo.VpcID),
		d.Set("resolver_rule_count", endpointInfo.ResolverRuleCount),
		d.Set("created_at", endpointInfo.CreateTime),
		d.Set("updated_at", endpointInfo.UpdateTime),
		d.Set("ip_addresses", ipAddresses),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}

	return nil
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
		Refresh:      waitForDNSEndpoint(dnsClient, d.Id()),
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
