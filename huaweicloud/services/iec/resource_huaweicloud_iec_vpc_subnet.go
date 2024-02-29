package iec

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/subnets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func buildSubnetDNSList(d *schema.ResourceData) []string {
	rawDNSN := d.Get("dns_list").([]interface{})

	// set the default DNS if it was not specified
	if len(rawDNSN) == 0 {
		return []string{"114.114.114.114", "8.8.8.8"}
	}

	dnsn := make([]string, len(rawDNSN))
	for i, raw := range rawDNSN {
		dnsn[i] = raw.(string)
	}
	return dnsn
}

// @API IEC POST /v1/subnets
// @API IEC GET /v1/subnets/{subnet_id}
// @API IEC PUT /v1/subnets/{subnet_id}
// @API IEC DELETE /v1/subnets/{subnet_id}
func ResourceSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetCreate,
		ReadContext:   resourceSubnetRead,
		UpdateContext: resourceSubnetUpdate,
		DeleteContext: resourceSubnetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_id": {
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
			"dhcp_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"dns_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"site_info": {
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

func resourceSubnetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	subnetClient, err := conf.IECV1Client(conf.GetRegion(d))

	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	dhcp := d.Get("dhcp_enable").(bool)
	createOpts := subnets.CreateOpts{
		Name:       d.Get("name").(string),
		Cidr:       d.Get("cidr").(string),
		VpcID:      d.Get("vpc_id").(string),
		SiteID:     d.Get("site_id").(string),
		GatewayIP:  d.Get("gateway_ip").(string),
		DhcpEnable: &dhcp,
		DNSList:    buildSubnetDNSList(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := subnets.Create(subnetClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IEC subnets: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] Waiting for IEC subnets (%s) to become active", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"UNKNOWN"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForSubnetStatus(subnetClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"error waiting for IEC subnets (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	return resourceSubnetRead(ctx, d, conf)
}

func resourceSubnetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	subnetClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	n, err := subnets.Get(subnetClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IEC subnets")
	}

	log.Printf("[DEBUG] IEC subnets %s: %+v", d.Id(), n)
	mErr := multierror.Append(
		nil,
		d.Set("name", n.Name),
		d.Set("cidr", n.Cidr),
		d.Set("vpc_id", n.VpcID),
		d.Set("site_id", n.SiteID),
		d.Set("gateway_ip", n.GatewayIP),
		d.Set("dhcp_enable", n.DhcpEnable),
		d.Set("dns_list", n.DNSList),
		d.Set("site_info", n.SiteInfo),
		d.Set("status", n.Status),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSubnetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	subnetClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	var updateOpts subnets.UpdateOpts

	// name is mandatory while updating subnets
	updateOpts.Name = d.Get("name").(string)

	if d.HasChange("dhcp_enable") {
		dhcp := d.Get("dhcp_enable").(bool)
		updateOpts.DhcpEnable = &dhcp
	}
	if d.HasChange("dns_list") {
		dnsList := utils.ExpandToStringList(d.Get("dns_list").([]interface{}))
		updateOpts.DNSList = &dnsList
	}

	_, err = subnets.Update(subnetClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating IEC subnets: %s", err)
	}

	return resourceSubnetRead(ctx, d, meta)
}

func resourceSubnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	subnetClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	err = subnets.Delete(subnetClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IEC subnets")
	}

	// waiting for subnets to become deleted
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "UNKNOWN"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSubnetStatus(subnetClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"error waiting for IEC subnets (%s) to become deleted: %s",
			d.Id(), stateErr)
	}

	return nil
}

func waitForSubnetStatus(subnetClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := subnets.Get(subnetClient, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted IEC subnets %s", id)
				return n, "DELETED", nil
			}
			return n, "ERROR", err
		}

		return n, n.Status, nil
	}
}
