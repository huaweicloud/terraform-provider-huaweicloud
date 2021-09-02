package waf

import (
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk/openstack/waf/v1/certificates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceWafCertificateV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafCertificateV1Create,
		Read:   resourceWafCertificateV1Read,
		Update: resourceWafCertificateV1Update,
		Delete: resourceWafCertificateV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-]{1,256}$`),
					"The maximum length is 256 characters. Only digits, letters, underscores (_), and hyphens (-) are allowed"),
			},
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWafCertificateV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	createOpts := certificates.CreateOpts{
		Name:    d.Get("name").(string),
		Content: strings.TrimSpace(d.Get("certificate").(string)),
		Key:     strings.TrimSpace(d.Get("private_key").(string)),
	}

	certificate, err := certificates.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("error creating WAF Certificate: %w", err)
	}

	logp.Printf("[DEBUG] Waf certificate created: %#v", certificate)
	d.SetId(certificate.Id)

	return resourceWafCertificateV1Read(d, meta)
}

func resourceWafCertificateV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	n, err := certificates.Get(wafClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Error obtain WAF certificate information")
	}

	expires := time.Unix(int64(n.ExpireTime/1000), 0).UTC().Format("2006-01-02 15:04:05 MST")

	d.Set("region", config.GetRegion(d))
	d.Set("name", n.Name)
	d.Set("expiration", expires)

	return nil
}

func resourceWafCertificateV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	updateOpts := certificates.UpdateOpts{
		Name: d.Get("name").(string),
	}

	_, err = certificates.Update(wafClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("error updating WAF Certificate: %w", err)
	}
	return resourceWafCertificateV1Read(d, meta)
}

func resourceWafCertificateV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	err = certificates.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("error deleting WAF Certificate: %s", err)
	}

	d.SetId("")
	return nil
}
