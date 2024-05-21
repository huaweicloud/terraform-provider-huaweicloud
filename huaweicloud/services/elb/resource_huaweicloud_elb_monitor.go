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

// @API ELB POST /v3/{project_id}/elb/healthmonitors
// @API ELB GET /v3/{project_id}/elb/healthmonitors/{healthmonitor_id}
// @API ELB PUT /v3/{project_id}/elb/healthmonitors/{healthmonitor_id}
// @API ELB DELETE /v3/{project_id}/elb/healthmonitors/{healthmonitor_id}
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
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
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
			"max_retries_down": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMonitorV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	createOpts := monitors.CreateOpts{
		PoolID:         d.Get("pool_id").(string),
		Type:           d.Get("protocol").(string),
		Delay:          d.Get("interval").(int),
		Timeout:        d.Get("timeout").(int),
		MaxRetries:     d.Get("max_retries").(int),
		MaxRetriesDown: d.Get("max_retries_down").(int),
		Name:           d.Get("name").(string),
		URLPath:        d.Get("url_path").(string),
		DomainName:     d.Get("domain_name").(string),
		MonitorPort:    d.Get("port").(int),
		ExpectedCodes:  d.Get("status_code").(string),
		HTTPMethod:     d.Get("http_method").(string),
		AdminStateUp:   &enabled,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	monitor, err := monitors.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to create monitor: %s", err)
	}

	d.SetId(monitor.ID)

	return resourceMonitorV3Read(ctx, d, meta)
}

func resourceMonitorV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	monitor, err := monitors.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "monitor")
	}

	log.Printf("[DEBUG] Retrieved monitor %s: %#v", d.Id(), monitor)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("protocol", monitor.Type),
		d.Set("interval", monitor.Delay),
		d.Set("timeout", monitor.Timeout),
		d.Set("max_retries", monitor.MaxRetries),
		d.Set("url_path", monitor.URLPath),
		d.Set("domain_name", monitor.DomainName),
		d.Set("status_code", monitor.ExpectedCodes),
		d.Set("name", monitor.Name),
		d.Set("max_retries_down", monitor.MaxRetriesDown),
		d.Set("enabled", monitor.AdminStateUp),
		d.Set("http_method", monitor.HTTPMethod),
		d.Set("created_at", monitor.CreatedAt),
		d.Set("updated_at", monitor.UpdatedAt),
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
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
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
	if d.HasChange("status_code") {
		updateOpts.ExpectedCodes = d.Get("status_code").(string)
	}
	if d.HasChange("protocol") {
		updateOpts.Type = d.Get("protocol").(string)
	}
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("max_retries_down") {
		updateOpts.MaxRetriesDown = d.Get("max_retries_down").(int)
	}
	if d.HasChange("http_method") {
		updateOpts.HTTPMethod = d.Get("http_method").(string)
	}
	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.AdminStateUp = &enabled
	}

	log.Printf("[DEBUG] Updating monitor %s with options: %#v", d.Id(), updateOpts)
	_, err = monitors.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to update monitor %s: %s", d.Id(), err)
	}

	return resourceMonitorV3Read(ctx, d, meta)
}

func resourceMonitorV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	log.Printf("[DEBUG] Deleting monitor %s", d.Id())
	err = monitors.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("unable to delete monitor %s: %s", d.Id(), err)
	}

	return nil
}
