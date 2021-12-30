package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/monitors"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceMonitorV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorV2Create,
		Read:   resourceMonitorV2Read,
		Update: resourceMonitorV2Update,
		Delete: resourceMonitorV2Delete,

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
				Optional: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"delay": {
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

			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"expected_codes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "tenant_id is deprecated",
			},
		},
	}
}

func resourceMonitorV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := monitors.CreateOpts{
		PoolID:        d.Get("pool_id").(string),
		TenantID:      d.Get("tenant_id").(string),
		Type:          d.Get("type").(string),
		Delay:         d.Get("delay").(int),
		Timeout:       d.Get("timeout").(int),
		MaxRetries:    d.Get("max_retries").(int),
		URLPath:       d.Get("url_path").(string),
		HTTPMethod:    d.Get("http_method").(string),
		ExpectedCodes: d.Get("expected_codes").(string),
		Name:          d.Get("name").(string),
		MonitorPort:   d.Get("port").(int),
		AdminStateUp:  &adminStateUp,
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	poolID := createOpts.PoolID
	err = waitForLBV2viaPool(lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return err
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	monitor, err := monitors.Create(lbClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to create monitor: %s", err)
	}

	err = waitForLBV2viaPool(lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return err
	}

	d.SetId(monitor.ID)

	return resourceMonitorV2Read(d, meta)
}

func resourceMonitorV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	monitor, err := monitors.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "monitor")
	}

	logp.Printf("[DEBUG] Retrieved monitor %s: %#v", d.Id(), monitor)

	d.Set("tenant_id", monitor.TenantID)
	d.Set("type", monitor.Type)
	d.Set("delay", monitor.Delay)
	d.Set("timeout", monitor.Timeout)
	d.Set("max_retries", monitor.MaxRetries)
	d.Set("url_path", monitor.URLPath)
	d.Set("http_method", monitor.HTTPMethod)
	d.Set("expected_codes", monitor.ExpectedCodes)
	d.Set("admin_state_up", monitor.AdminStateUp)
	d.Set("name", monitor.Name)
	d.Set("region", GetRegion(d, config))
	if monitor.MonitorPort != 0 {
		d.Set("port", monitor.MonitorPort)
	}

	return nil
}

func resourceMonitorV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts monitors.UpdateOpts
	if d.HasChange("url_path") {
		updateOpts.URLPath = d.Get("url_path").(string)
	}
	if d.HasChange("expected_codes") {
		updateOpts.ExpectedCodes = d.Get("expected_codes").(string)
	}
	if d.HasChange("delay") {
		updateOpts.Delay = d.Get("delay").(int)
	}
	if d.HasChange("timeout") {
		updateOpts.Timeout = d.Get("timeout").(int)
	}
	if d.HasChange("max_retries") {
		updateOpts.MaxRetries = d.Get("max_retries").(int)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("port") {
		updateOpts.MonitorPort = d.Get("port").(int)
	}
	if d.HasChange("http_method") {
		updateOpts.HTTPMethod = d.Get("http_method").(string)
	}

	logp.Printf("[DEBUG] Updating monitor %s with options: %#v", d.Id(), updateOpts)
	timeout := d.Timeout(schema.TimeoutUpdate)
	poolID := d.Get("pool_id").(string)
	err = waitForLBV2viaPool(lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return err
	}
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = monitors.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.Errorf("Unable to update monitor %s: %s", d.Id(), err)
	}

	// Wait for LB to become active before continuing
	err = waitForLBV2viaPool(lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return err
	}

	return resourceMonitorV2Read(d, meta)
}

func resourceMonitorV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting monitor %s", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	poolID := d.Get("pool_id").(string)
	err = waitForLBV2viaPool(lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return err
	}
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = monitors.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.Errorf("Unable to delete monitor %s: %s", d.Id(), err)
	}

	err = waitForLBV2viaPool(lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return err
	}

	return nil
}
