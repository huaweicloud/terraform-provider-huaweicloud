package iec

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
	ieccommon "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/ports"
	"github.com/chnsz/golangsdk/openstack/iec/v1/subnets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC GET /v1/subnets/{subnet_id}
// @API IEC DELETE /v1/ports/{port_id}
// @API IEC GET /v1/ports/{port_id}
// @API IEC PUT /v1/ports/{port_id}
// @API IEC POST /v1/ports
func ResourceVip() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVipCreate,
		UpdateContext: resourceVipUpdate,
		ReadContext:   resourceVipRead,
		DeleteContext: resourceVipDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"subnet_id": {
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
			"port_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allowed_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func getVipPortIDs(d *schema.ResourceData) []string {
	rawPortIDs := d.Get("port_ids").([]interface{})
	portids := make([]string, len(rawPortIDs))
	for i, raw := range rawPortIDs {
		portids[i] = raw.(string)
	}
	return portids
}

func resourceVipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	iecClient, err := conf.IECV1Client(region)
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	networkID := d.Get("subnet_id").(string)
	n, err := subnets.Get(iecClient, networkID).Extract()
	if err != nil {
		return diag.Errorf("error retrieving IEC subnet %s: %s", networkID, err)
	}

	createOpts := ports.CreateOpts{
		NetworkId:   networkID,
		DeviceOwner: "neutron:VIP_PORT",
	}

	// Construct fixed ip
	if fixedIP := d.Get("ip_address").(string); fixedIP != "" {
		fixip := make([]ports.FixIPEntity, 1)
		fixip[0] = ports.FixIPEntity{
			SubnetID:  n.NeutronSubnetID,
			IPAddress: fixedIP,
		}
		createOpts.FixedIPs = fixip
	}

	p, err := ports.Create(iecClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IEC port: %s", err)
	}

	d.SetId(p.ID)
	log.Printf("[DEBUG] Waiting for IEC Port (%s) to become available.", p.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitingForIECVIPActive(iecClient, p.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error Waiting for IEC Port (%s) to become available: %s", p.ID, err)
	}

	// associate ports with the vip
	portids := getVipPortIDs(d)
	if len(portids) > 0 {
		if err = updateVipAssociate(iecClient, p.ID, portids); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVipRead(ctx, d, meta)
}

func resourceVipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	iecClient, err := conf.IECV1Client(region)
	var mErr *multierror.Error
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	vip, err := ports.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IEC VPC")
	}

	var ipAddr string
	if len(vip.FixedIPs) > 0 {
		ipAddr = vip.FixedIPs[0].IpAddress
	}

	allPortAddrs := make([]string, len(vip.AllowedAddressPairs))
	for i, pair := range vip.AllowedAddressPairs {
		allPortAddrs[i] = pair.IpAddress
	}
	mErr = multierror.Append(
		mErr,
		d.Set("subnet_id", vip.NetworkID),
		d.Set("mac_address", vip.MacAddress),
		d.Set("ip_address", ipAddr),
		d.Set("allowed_addresses", allPortAddrs),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	iecClient, err := conf.IECV1Client(region)
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	portids := getVipPortIDs(d)
	if err = updateVipAssociate(iecClient, d.Id(), portids); err != nil {
		return diag.Errorf("error updateVipAssociate: %s", err)
	}

	return resourceVipRead(ctx, d, meta)
}

func resourceVipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	iecClient, err := conf.IECV1Client(region)
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	if len(getVipPortIDs(d)) > 0 {
		// disassociate ports
		if err := updateVipAssociate(iecClient, d.Id(), []string{}); err != nil {
			return diag.Errorf("error updateVipAssociate: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitingForIECVIPDelete(iecClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for IEC vip ports (%s) to become deleted: %s", d.Id(), err)
	}

	return nil
}

func waitingForIECVIPActive(client *golangsdk.ServiceClient, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, err := ports.Get(client, portID).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Neutron Port: %+v", p)
		if p.Status == "DOWN" || p.Status == "ACTIVE" {
			return p, "ACTIVE", nil
		}

		return p, p.Status, nil
	}
}

func waitingForIECVIPDelete(client *golangsdk.ServiceClient, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete IEC Port %s", portID)
		port, err := ports.Get(client, portID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted IEC Port %s", portID)
				return port, "DELETED", nil
			}
			return port, "ACTIVATE", err
		}
		err = ports.Delete(client, portID).ExtractErr()

		// remote service will return code 204 when delete success
		if err == nil {
			log.Printf("[DEBUG] Successfully deleted IEC Port %s", portID)
			return port, "DELETED", nil
		}

		log.Printf("[DEBUG] IEC Port %s still active.\n", portID)
		return port, "ACTIVE", nil
	}
}

func updateVipAssociate(client *golangsdk.ServiceClient, vipID string, portIDs []string) error {
	allAddrs := make([]string, len(portIDs))
	action := "associate"
	if len(portIDs) == 0 {
		action = "disassociate"
	}

	// check the port id and get ip address
	for i, portid := range portIDs {
		port, err := ports.Get(client, portid).Extract()
		if err != nil {
			return fmt.Errorf("error fetching port %s: %s", portid, err)
		}

		if len(port.FixedIPs) == 0 {
			return fmt.Errorf("port %s has no ip address, Error associate it", portid)
		}
		allAddrs[i] = port.FixedIPs[0].IpAddress
	}

	// construct allowed address pairs
	allowedPairs := make([]ieccommon.AllowedAddressPair, len(allAddrs))
	for i, addr := range allAddrs {
		allowedPairs[i] = ieccommon.AllowedAddressPair{
			IpAddress: addr,
		}
	}
	// associate/disassociate ports with the vip
	associateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &allowedPairs,
	}
	log.Printf("[DEBUG] VIP %s %s with options: %#v", action, vipID, associateOpts)
	_, err := ports.Update(client, vipID, associateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error %s vip: %s", action, err)
	}

	// Update the allowed-address-pairs of the port to 1.1.1.1/0
	// to disable the source/destination check
	portpairs := make([]ieccommon.AllowedAddressPair, 1)
	portpairs[0] = ieccommon.AllowedAddressPair{
		IpAddress: "1.1.1.1/0",
	}
	portUpdateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &portpairs,
	}

	for _, portid := range portIDs {
		_, err = ports.Update(client, portid, portUpdateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error update port %s: %s", portid, err)
		}
	}

	return nil
}
