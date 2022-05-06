package gaussdb

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

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

		Schema: map[string]*schema.Schema{ //request and response parameters
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
	config := meta.(*config.Config)
	client, err := config.GaussdbV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	createOpts := instances.ProxyOpts{
		Flavor:  d.Get("flavor").(string),
		NodeNum: d.Get("node_num").(int),
	}

	instance_id := d.Get("instance_id").(string)
	n, err := instances.EnableProxy(client, instance_id, createOpts).ExtractJobResponse()
	if err != nil {
		return fmtp.DiagErrorf("Error creating gaussdb_mysql_proxy: %s", err)
	}
	d.SetId(instance_id)

	if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return fmtp.DiagErrorf("Error waiting for gaussdb_mysql_proxy job: %s", err)
	}

	return resourceGaussDBProxyRead(ctx, d, meta)
}

func resourceGaussDBProxyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.GaussdbV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	proxy, err := instances.GetProxy(client, d.Id()).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error fetchting gaussdb_mysql_proxy: %s", err)
	}
	mErr := multierror.Append(nil,
		d.Set("instance_id", d.Id()),
		d.Set("flavor", proxy.Flavor),
		d.Set("node_num", proxy.NodeNum),
		d.Set("address", proxy.Address),
		d.Set("port", proxy.Port),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGaussDBProxyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.GaussdbV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	if d.HasChange("node_num") {
		oldnum, newnum := d.GetChange("node_num")
		if newnum.(int) < oldnum.(int) {
			return fmtp.DiagErrorf("Error updating gaussdb_mysql_proxy %s: new node num should be greater than old num", d.Id())
		}

		enlarge_size := newnum.(int) - oldnum.(int)
		enlargeProxyOpts := instances.EnlargeProxyOpts{
			NodeNum: enlarge_size,
		}

		lp, err := instances.EnlargeProxy(client, d.Id(), enlargeProxyOpts).ExtractJobResponse()
		if err != nil {
			return fmtp.DiagErrorf("Error enlarging gaussdb_mysql_proxy: %s", err)
		}

		if err = instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), lp.JobID); err != nil {
			return fmtp.DiagErrorf("Error waiting for gaussdb_mysql_proxy job: %s", err)
		}
	}

	return resourceGaussDBProxyRead(ctx, d, meta)
}

func resourceGaussDBProxyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.GaussdbV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	dp, err := instances.DeleteProxy(client, d.Id()).ExtractJobResponse()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting gaussdb_mysql_proxy: %s", err)
	}

	if err = instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutDelete)/time.Second), dp.JobID); err != nil {
		return fmtp.DiagErrorf("Error waiting for gaussdb_mysql_proxy job: %s", err)
	}

	d.SetId("")
	return nil
}
