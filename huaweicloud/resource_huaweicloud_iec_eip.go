package huaweicloud

import (
	"strconv"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/publicips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceIecNetworkEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceIecEipV1Create,
		Read:   resourceIecEipV1Read,
		Update: resourceIecEipV1Update,
		Delete: resourceIecEipV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceIecEipV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	eipClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
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

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := publicips.Create(eipClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC public ip: %s", err)
	}

	logp.Printf("[DEBUG] IEC publicips ID: %s", n.ID)
	d.SetId(n.ID)

	logp.Printf("[DEBUG] Waiting for public ip (%s) to become active", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE", "UNBOUND"},
		Refresh:    waitForIecEipStatus(eipClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmtp.Errorf(
			"Error waiting for public ip (%s) to become ACTIVE: %s",
			d.Id(), stateErr)
	}

	if bindPort := d.Get("port_id").(string); bindPort != "" {
		logp.Printf("[DEBUG] bind public ip %s to port %s", d.Id(), bindPort)
		if err := operateOnPort(d, eipClient, bindPort); err != nil {
			return err
		}
	}

	return resourceIecEipV1Read(d, config)
}

func resourceIecEipV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	eipClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
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

		return fmtp.Errorf("Error retrieving Huaweicloud IEC public ip: %s", err)
	}

	logp.Printf("[DEBUG] IEC public ip %s: %+v", d.Id(), n)

	d.Set("site_id", n.SiteID)
	d.Set("line_id", n.Type)
	d.Set("port_id", n.PortID)
	d.Set("public_ip", n.PublicIpAddress)
	d.Set("private_ip", n.PrivateIpAddress)
	d.Set("ip_version", n.IPVersion)
	d.Set("bandwidth_id", n.BandwidthID)
	d.Set("bandwidth_name", n.BandwidthName)
	d.Set("bandwidth_size", n.BandwidthSize)
	d.Set("bandwidth_share_type", n.BandwidthShareType)
	d.Set("site_info", n.SiteInfo)
	d.Set("status", eip.NormalizeEIPStatus(n.Status))

	return nil
}

func resourceIecEipV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	eipClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	if d.HasChange("port_id") {
		var opErr error
		oPort, nPort := d.GetChange("port_id")
		if oldPort := oPort.(string); oldPort != "" {
			logp.Printf("[DEBUG] unbind public ip %s from port %s", d.Id(), oldPort)
			opErr = operateOnPort(d, eipClient, "")
		}

		if newPort := nPort.(string); newPort != "" {
			logp.Printf("[DEBUG] bind public ip %s to port %s", d.Id(), newPort)
			opErr = operateOnPort(d, eipClient, newPort)
		}

		if opErr != nil {
			return opErr
		}
	}

	return resourceIecEipV1Read(d, meta)
}

func resourceIecEipV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	eipClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	// unbound the port before deleting the publicips
	if port := d.Get("port_id").(string); port != "" {
		logp.Printf("[DEBUG] unbind public ip %s from port %s", d.Id(), port)
		if err := operateOnPort(d, eipClient, ""); err != nil {
			return err
		}
	}

	err = publicips.Delete(eipClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting Huaweicloud IEC public ip")
	}

	// waiting for public ip to become deleted
	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Refresh:    waitForIecEipStatus(eipClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmtp.Errorf(
			"Error waiting for Subnet (%s) to become deleted: %s",
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
		var action string = "binding"
		if port == "" {
			action = "unbinding"
		}
		return fmtp.Errorf("Error %s Huaweicloud IEC public ip: %s", action, err)
	}
	return nil
}

func waitForIecEipStatus(subnetClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := publicips.Get(subnetClient, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				logp.Printf("[INFO] Successfully deleted Huaweicloud IEC public ip %s", id)
				return n, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[INFO] Successfully deleted Huaweicloud IEC public ip %s", id)
				return n, "DELETED", nil
			}

			return n, "ERROR", err
		}

		if n.Status == "ERROR" || n.Status == "BIND_ERROR" {
			return n, n.Status, fmtp.Errorf("got error status with the public ip")
		}

		// "DOWN" means the publicips is active but unbound
		if n.Status == "DOWN" {
			return n, "UNBOUND", nil
		}

		return n, n.Status, nil
	}
}
