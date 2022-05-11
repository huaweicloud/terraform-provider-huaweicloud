package lb

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk/openstack/elb/v2/certificates"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceCertificateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateV2Create,
		ReadContext:   resourceCertificateV2Read,
		UpdateContext: resourceCertificateV2Update,
		DeleteContext: resourceCertificateV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "server",
				ValidateFunc: validation.StringInSlice([]string{
					"server", "client",
				}, true),
			},

			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"private_key": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
				Sensitive:        true,
			},

			"certificate": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
				Sensitive:        true,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCertificateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	createOpts := certificates.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Type:                d.Get("type").(string),
		Domain:              d.Get("domain").(string),
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Add certificate and private_key here so they wouldn't go in the above log entry
	createOpts.Certificate = d.Get("certificate").(string)
	createOpts.PrivateKey = d.Get("private_key").(string)

	c, err := certificates.Create(elbClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating Certificate: %s", err)
	}

	// If all has been successful, set the ID on the resource
	d.SetId(c.Id)

	return resourceCertificateV2Read(ctx, d, meta)
}

func resourceCertificateV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	c, err := certificates.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving certificate")
	}
	logp.Printf("[DEBUG] Retrieved certificate %s: %#v", d.Id(), c)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", c.Name),
		d.Set("description", c.Description),
		d.Set("type", c.Type),
		d.Set("domain", c.Domain),
		d.Set("certificate", c.Certificate),
		d.Set("private_key", c.PrivateKey),
		d.Set("create_time", c.CreateTime),
		d.Set("update_time", c.UpdateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting certificate fields: %s", err)
	}

	return nil
}

func resourceCertificateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts certificates.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("domain") {
		updateOpts.Domain = d.Get("domain").(string)
	}

	logp.Printf("[DEBUG] Updating certificate %s with options: %#v", d.Id(), updateOpts)

	// Add certificate and private_key here so they wouldn't go in the above log entry
	if d.HasChange("private_key") {
		updateOpts.PrivateKey = d.Get("private_key").(string)
	}
	if d.HasChange("certificate") {
		updateOpts.Certificate = d.Get("certificate").(string)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err := certificates.Update(elbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmtp.DiagErrorf("Error updating certificate %s: %s", d.Id(), err)
	}

	return resourceCertificateV2Read(ctx, d, meta)
}

func resourceCertificateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting certificate %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := certificates.Delete(elbClient, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if utils.IsResourceNotFound(err) {
			logp.Printf("[INFO] deleting an unavailable certificate: %s", d.Id())
			return nil
		}
		return fmtp.DiagErrorf("Error deleting certificate %s: %s", d.Id(), err)
	}

	return nil
}
