package elb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/elb/v3/certificates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceCertificateV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateV3Create,
		ReadContext:   resourceCertificateV3Read,
		UpdateContext: resourceCertificateV3Update,
		DeleteContext: resourceCertificateV3Delete,
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
			},

			"certificate": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
			},

			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"expire_time": {
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

func resourceCertificateV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	createOpts := certificates.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Type:                d.Get("type").(string),
		Domain:              d.Get("domain").(string),
		PrivateKey:          d.Get("private_key").(string),
		Certificate:         d.Get("certificate").(string),
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	c, err := certificates.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating Certificate: %s", err)
	}

	// If all has been successful, set the ID on the resource
	d.SetId(c.ID)

	return resourceCertificateV3Read(ctx, d, meta)
}

func resourceCertificateV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	c, err := certificates.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "certificate")
	}
	log.Printf("[DEBUG] Retrieved certificate %s: %#v", d.Id(), c)

	mErr := multierror.Append(nil,
		d.Set("name", c.Name),
		d.Set("description", c.Description),
		d.Set("type", c.Type),
		d.Set("domain", c.Domain),
		d.Set("certificate", c.Certificate),
		d.Set("private_key", c.PrivateKey),
		d.Set("create_time", c.CreateTime),
		d.Set("update_time", c.UpdateTime),
		d.Set("expire_time", c.ExpireTime),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB Certificate fields: %s", err)
	}

	return nil
}

func resourceCertificateV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	var updateOpts certificates.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		desc := d.Get("description").(string)
		updateOpts.Description = &desc
	}
	if d.HasChange("domain") {
		updateOpts.Domain = d.Get("domain").(string)
	}
	if d.HasChange("private_key") {
		updateOpts.PrivateKey = d.Get("private_key").(string)
	}
	if d.HasChange("certificate") {
		updateOpts.Certificate = d.Get("certificate").(string)
	}

	log.Printf("[DEBUG] Updating certificate %s with options: %#v", d.Id(), updateOpts)

	_, err = certificates.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating certificate %s: %s", d.Id(), err)
	}

	return resourceCertificateV3Read(ctx, d, meta)
}

func resourceCertificateV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	log.Printf("[DEBUG] Deleting certificate %s", d.Id())
	err = certificates.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		if utils.IsResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable certificate: %s", d.Id())
			return nil
		}
		return diag.Errorf("error deleting certificate %s: %s", d.Id(), err)
	}

	return nil
}
