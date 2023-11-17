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
)

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
			"ip_addresses": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 6,
				MinItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
