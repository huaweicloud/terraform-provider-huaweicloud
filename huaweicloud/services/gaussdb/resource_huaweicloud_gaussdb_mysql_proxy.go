package gaussdb

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy/enlarge
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/proxy
func ResourceGaussDBProxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBProxyCreate,
		ReadContext:   resourceGaussDBProxyRead,
		UpdateContext: resourceGaussDBProxyUpdate,
		DeleteContext: resourceGaussDBProxyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceGaussDBProxyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GaussdbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s ", err)
	}

	createOpts := instances.ProxyOpts{
		Flavor:  d.Get("flavor").(string),
		NodeNum: d.Get("node_num").(int),
	}

	instanceId := d.Get("instance_id").(string)
	n, err := instances.EnableProxy(client, instanceId, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating gaussdb_mysql_proxy: %s", err)
	}
	d.SetId(instanceId)

	if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return diag.Errorf("error waiting for gaussdb_mysql_proxy job: %s", err)
	}

	return resourceGaussDBProxyRead(ctx, d, meta)
}

func resourceGaussDBProxyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBProxyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GaussdbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s ", err)
	}

	if d.HasChange("node_num") {
		oldnum, newnum := d.GetChange("node_num")
		if newnum.(int) < oldnum.(int) {
			return diag.Errorf("error updating gaussdb_mysql_proxy %s: new node num should be greater than old num", d.Id())
		}

		enlargeSize := newnum.(int) - oldnum.(int)
		enlargeProxyOpts := instances.EnlargeProxyOpts{
			NodeNum: enlargeSize,
		}

		lp, err := instances.EnlargeProxy(client, d.Id(), enlargeProxyOpts).ExtractJobResponse()
		if err != nil {
			return diag.Errorf("error enlarging gaussdb_mysql_proxy: %s", err)
		}

		if err = instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), lp.JobID); err != nil {
			return diag.Errorf("error waiting for gaussdb_mysql_proxy job: %s", err)
		}
	}

	return resourceGaussDBProxyRead(ctx, d, meta)
}

func resourceGaussDBProxyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GaussdbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s ", err)
	}

	dp, err := instances.DeleteProxy(client, d.Id()).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error deleting gaussdb_mysql_proxy: %s", err)
	}

	if err = instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutDelete)/time.Second), dp.JobID); err != nil {
		return diag.Errorf("error waiting for gaussdb_mysql_proxy job: %s", err)
	}

	return nil
}
