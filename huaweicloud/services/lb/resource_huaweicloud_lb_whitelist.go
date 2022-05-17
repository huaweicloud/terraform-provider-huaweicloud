package lb

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk/openstack/elb/v2/whitelists"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceWhitelistV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWhitelistV2Create,
		ReadContext:   resourceWhitelistV2Read,
		UpdateContext: resourceWhitelistV2Update,
		DeleteContext: resourceWhitelistV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "tenant_id is deprecated",
			},

			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enable_whitelist": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"whitelist": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressLBWhitelistDiffs,
			},
		},
	}
}

func resourceWhitelistV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	enableWhitelist := d.Get("enable_whitelist").(bool)
	createOpts := whitelists.CreateOpts{
		TenantId:        d.Get("tenant_id").(string),
		ListenerId:      d.Get("listener_id").(string),
		EnableWhitelist: &enableWhitelist,
		Whitelist:       d.Get("whitelist").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	wl, err := whitelists.Create(elbClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud Whitelist: %s", err)
	}

	d.SetId(wl.ID)
	return resourceWhitelistV2Read(ctx, d, meta)
}

func resourceWhitelistV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	wl, err := whitelists.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving whitelist")
	}

	logp.Printf("[DEBUG] Retrieved whitelist %s: %#v", d.Id(), wl)

	d.SetId(wl.ID)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("tenant_id", wl.TenantId),
		d.Set("listener_id", wl.ListenerId),
		d.Set("enable_whitelist", wl.EnableWhitelist),
		d.Set("whitelist", wl.Whitelist),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting whitelist fields: %s", err)
	}

	return nil
}

func resourceWhitelistV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts whitelists.UpdateOpts
	if d.HasChange("enable_whitelist") {
		ew := d.Get("enable_whitelist").(bool)
		updateOpts.EnableWhitelist = &ew
	}
	if d.HasChange("whitelist") {
		updateOpts.Whitelist = d.Get("whitelist").(string)
	}

	logp.Printf("[DEBUG] Updating whitelist %s with options: %#v", d.Id(), updateOpts)
	_, err = whitelists.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to update whitelist %s: %s", d.Id(), err)
	}

	return resourceWhitelistV2Read(ctx, d, meta)
}

func resourceWhitelistV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	logp.Printf("[DEBUG] Attempting to delete whitelist %s", d.Id())
	err = whitelists.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud whitelist: %s", err)
	}
	d.SetId("")
	return nil
}
