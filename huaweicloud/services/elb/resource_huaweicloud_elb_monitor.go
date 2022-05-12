package elb

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v3/monitors"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceMonitorV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitorV3Create,
		ReadContext:   resourceMonitorV3Read,
		UpdateContext: resourceMonitorV3Update,
		DeleteContext: resourceMonitorV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

func resourceMonitorV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
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

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	monitor, err := monitors.Create(lbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to create monitor: %s", err)
	}

	d.SetId(monitor.ID)

	return resourceMonitorV3Read(ctx, d, meta)
}

func resourceMonitorV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	monitor, err := monitors.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "monitor")
	}

	log.Printf("[DEBUG] Retrieved monitor %s: %#v", d.Id(), monitor)

	mErr := multierror.Append(nil,
		d.Set("protocol", monitor.Type),
		d.Set("interval", monitor.Delay),
		d.Set("timeout", monitor.Timeout),
		d.Set("max_retries", monitor.MaxRetries),
		d.Set("url_path", monitor.URLPath),
		d.Set("domain_name", monitor.DomainName),
		d.Set("region", config.GetRegion(d)),
	)

	if len(monitor.Pools) != 0 {
		mErr = multierror.Append(mErr, d.Set("pool_id", monitor.Pools[0].ID))
	}

	if monitor.MonitorPort != 0 {
		mErr = multierror.Append(mErr, d.Set("port", monitor.MonitorPort))
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB monitor fields: %s", err)
	}

	return nil
}

func resourceMonitorV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
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

	log.Printf("[DEBUG] Updating monitor %s with options: %#v", d.Id(), updateOpts)
	_, err = monitors.Update(lbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to update monitor %s: %s", d.Id(), err)
	}

	return resourceMonitorV3Read(ctx, d, meta)
}

func resourceMonitorV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	log.Printf("[DEBUG] Deleting monitor %s", d.Id())
	err = monitors.Delete(lbClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("unable to delete monitor %s: %s", d.Id(), err)
	}

	return nil
}
