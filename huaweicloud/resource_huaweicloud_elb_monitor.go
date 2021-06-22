package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk/openstack/elb/v3/monitors"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceMonitorV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorV3Create,
		Read:   resourceMonitorV3Read,
		Update: resourceMonitorV3Update,
		Delete: resourceMonitorV3Delete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"interval": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"max_retries": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceMonitorV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	createOpts := monitors.CreateOpts{
		PoolID:      d.Get("pool_id").(string),
		Type:        d.Get("protocol").(string),
		Delay:       d.Get("interval").(int),
		Timeout:     d.Get("timeout").(int),
		MaxRetries:  d.Get("max_retries").(int),
		URLPath:     d.Get("url_path").(string),
		DomainName:  d.Get("domain_name").(string),
		MonitorPort: d.Get("port").(int),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	monitor, err := monitors.Create(lbClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to create monitor: %s", err)
	}

	d.SetId(monitor.ID)

	return resourceMonitorV3Read(d, meta)
}

func resourceMonitorV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	monitor, err := monitors.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "monitor")
	}

	logp.Printf("[DEBUG] Retrieved monitor %s: %#v", d.Id(), monitor)

	d.Set("protocol", monitor.Type)
	d.Set("interval", monitor.Delay)
	d.Set("timeout", monitor.Timeout)
	d.Set("max_retries", monitor.MaxRetries)
	d.Set("url_path", monitor.URLPath)
	d.Set("domain_name", monitor.DomainName)
	d.Set("region", GetRegion(d, config))
	if monitor.MonitorPort != 0 {
		d.Set("port", monitor.MonitorPort)
	}

	return nil
}

func resourceMonitorV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts monitors.UpdateOpts
	if d.HasChange("url_path") {
		updateOpts.URLPath = d.Get("url_path").(string)
	}
	if d.HasChange("interval") {
		updateOpts.Delay = d.Get("interval").(int)
	}
	if d.HasChange("timeout") {
		updateOpts.Timeout = d.Get("timeout").(int)
	}
	if d.HasChange("max_retries") {
		updateOpts.MaxRetries = d.Get("max_retries").(int)
	}
	if d.HasChange("domain_name") {
		updateOpts.DomainName = d.Get("domain_name").(string)
	}
	if d.HasChange("port") {
		updateOpts.MonitorPort = d.Get("port").(int)
	}

	logp.Printf("[DEBUG] Updating monitor %s with options: %#v", d.Id(), updateOpts)
	_, err = monitors.Update(lbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to update monitor %s: %s", d.Id(), err)
	}

	return resourceMonitorV3Read(d, meta)
}

func resourceMonitorV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting monitor %s", d.Id())
	err = monitors.Delete(lbClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Unable to delete monitor %s: %s", d.Id(), err)
	}

	return nil
}
