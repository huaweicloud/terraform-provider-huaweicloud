package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/certificates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceCertificateV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateV2Create,
		Read:   resourceCertificateV2Read,
		Update: resourceCertificateV2Update,
		Delete: resourceCertificateV2Delete,

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

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCertificateV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	createOpts := certificates.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Domain:      d.Get("domain").(string),
		PrivateKey:  d.Get("private_key").(string),
		Certificate: d.Get("certificate").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	c, err := certificates.Create(elbClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating Certificate: %s", err)
	}

	// If all has been successful, set the ID on the resource
	d.SetId(c.ID)

	return resourceCertificateV2Read(d, meta)
}

func resourceCertificateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	c, err := certificates.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "certificate")
	}
	logp.Printf("[DEBUG] Retrieved certificate %s: %#v", d.Id(), c)

	d.Set("name", c.Name)
	d.Set("description", c.Description)
	d.Set("type", c.Type)
	d.Set("domain", c.Domain)
	d.Set("certificate", c.Certificate)
	d.Set("private_key", c.PrivateKey)
	d.Set("create_time", c.CreateTime)
	d.Set("update_time", c.UpdateTime)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCertificateV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
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
	if d.HasChange("private_key") {
		updateOpts.PrivateKey = d.Get("private_key").(string)
	}
	if d.HasChange("certificate") {
		updateOpts.Certificate = d.Get("certificate").(string)
	}

	logp.Printf("[DEBUG] Updating certificate %s with options: %#v", d.Id(), updateOpts)

	timeout := d.Timeout(schema.TimeoutUpdate)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err := certificates.Update(elbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmtp.Errorf("Error updating certificate %s: %s", d.Id(), err)
	}

	return resourceCertificateV2Read(d, meta)
}

func resourceCertificateV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting certificate %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := certificates.Delete(elbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if utils.IsResourceNotFound(err) {
			logp.Printf("[INFO] deleting an unavailable certificate: %s", d.Id())
			return nil
		}
		return fmtp.Errorf("Error deleting certificate %s: %s", d.Id(), err)
	}

	return nil
}
