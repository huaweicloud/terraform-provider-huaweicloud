package iec

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/publicips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
)

// @API IEC DELETE /v1/publicips/{publicip_id}
// @API IEC GET /v1/publicips/{publicip_id}
// @API IEC PUT /v1/publicips/{publicip_id}
// @API IEC POST /v1/publicips
func ResourceEip() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipCreate,
		ReadContext:   resourceEipRead,
		UpdateContext: resourceEipUpdate,
		DeleteContext: resourceEipDelete,
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
				Computed: true,
				ForceNew: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"line_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{4}),
				Description:  "schema: Computed",
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth_share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	eipClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	createOpts := publicips.CreateOpts{
		Publicip: publicips.PublicIPRequest{
			SiteID: d.Get("site_id").(string),
			Type:   d.Get("line_id").(string),
		},
	}

	ipVersion := d.Get("ip_version").(int)
	if ipVersion != 0 {
		createOpts.Publicip.IPVersion = strconv.Itoa(ipVersion)
	}

	log.Printf("[DEBUG] create Options: %#v", createOpts)
	n, err := publicips.Create(eipClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IEC public IP: %s", err)
	}

	log.Printf("[DEBUG] IEC publicips ID: %s", n.ID)
	d.SetId(n.ID)

	log.Printf("[DEBUG] waiting for public IP (%s) to become active", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE", "UNBOUND"},
		Refresh:    waitForEipStatus(eipClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"error waiting for public IP (%s) to become ACTIVE: %s",
			d.Id(), stateErr)
	}

	if bindPort := d.Get("port_id").(string); bindPort != "" {
		log.Printf("[DEBUG] bind public IP %s to port %s", d.Id(), bindPort)
		if err := operateOnPort(d, eipClient, bindPort); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceEipRead(ctx, d, cfg)
}

func resourceEipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	eipClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	n, err := publicips.Get(eipClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		if _, ok := err.(golangsdk.ErrDefault400); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error retrieving IEC public IP: %s", err)
	}

	log.Printf("[DEBUG] IEC public IP %s: %+v", d.Id(), n)
	mErr := multierror.Append(
		nil,
		d.Set("site_id", n.SiteID),
		d.Set("line_id", n.Type),
		d.Set("port_id", n.PortID),
		d.Set("public_ip", n.PublicIpAddress),
		d.Set("private_ip", n.PrivateIpAddress),
		d.Set("ip_version", n.IPVersion),
		d.Set("bandwidth_id", n.BandwidthID),
		d.Set("bandwidth_name", n.BandwidthName),
		d.Set("bandwidth_size", n.BandwidthSize),
		d.Set("bandwidth_share_type", n.BandwidthShareType),
		d.Set("site_info", n.SiteInfo),
		d.Set("status", eip.NormalizeEipStatus(n.Status)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}
	return nil
}

func resourceEipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	eipClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	if d.HasChange("port_id") {
		var opErr error
		oPort, nPort := d.GetChange("port_id")
		if oldPort := oPort.(string); oldPort != "" {
			log.Printf("[DEBUG] unbind public IP %s from port %s", d.Id(), oldPort)
			opErr = operateOnPort(d, eipClient, "")
		}

		if newPort := nPort.(string); newPort != "" {
			log.Printf("[DEBUG] bind public IP %s to port %s", d.Id(), newPort)
			opErr = operateOnPort(d, eipClient, newPort)
		}

		if opErr != nil {
			return diag.FromErr(err)
		}
	}

	return resourceEipRead(ctx, d, meta)
}

func resourceEipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	eipClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	// unbound the port before deleting the publicips
	if port := d.Get("port_id").(string); port != "" {
		log.Printf("[DEBUG] unbind public IP %s from port %s", d.Id(), port)
		if err := operateOnPort(d, eipClient, ""); err != nil {
			return diag.FromErr(err)
		}
	}

	err = publicips.Delete(eipClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IEC public IP")
	}

	// waiting for public IP to become deleted
	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Refresh:    waitForEipStatus(eipClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"error waiting for EIP (%s) to become deleted: %s",
			d.Id(), stateErr)
	}

	d.SetId("")
	return nil
}

func operateOnPort(d *schema.ResourceData, client *golangsdk.ServiceClient, port string) error {
	updateOpts := publicips.UpdateOpts{
		PortId: port,
	}
	_, err := publicips.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		var action = "binding"
		if port == "" {
			action = "unbinding"
		}
		return fmt.Errorf("error %s IEC public IP: %s", action, err)
	}
	return nil
}

func waitForEipStatus(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := publicips.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				log.Printf("[INFO] successfully deleted IEC public IP %s", id)
				return n, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] successfully deleted IEC public IP %s", id)
				return n, "DELETED", nil
			}

			return n, "ERROR", err
		}

		if n.Status == "ERROR" || n.Status == "BIND_ERROR" {
			return n, n.Status, fmt.Errorf("got error status with the public IP")
		}

		// "DOWN" means the publicips is active but unbound
		if n.Status == "DOWN" {
			return n, "UNBOUND", nil
		}

		return n, n.Status, nil
	}
}
