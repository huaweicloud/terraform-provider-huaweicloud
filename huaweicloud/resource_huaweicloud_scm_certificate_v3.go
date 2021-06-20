package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/scm/v3/certificates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func resourceScmCertificateV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceScmCertificateV3Create,
		Read:   resourceScmCertificateV3Read,
		Delete: resourceScmCertificateV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
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
				Required: true,
				ForceNew: true,
			},
			"certificate_chain": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
			},
			"private_key": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
			},
			"certificate": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"brand": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"push_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"revoke_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sans": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_algrithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issue_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"not_before": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"not_after": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"validation_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"validity_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wildcard_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"authentifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
					}},
			},
		},
	}
}

func resourceScmCertificateV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ScmV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SCM client: %s", err)
	}

	importOpts := certificates.ImportOpts{
		Name:             d.Get("name").(string),
		PrivateKey:       d.Get("private_key").(string),
		Certificate:      d.Get("certificate").(string),
		CertificateChain: d.Get("certificate_chain").(string),
	}

	c, err := certificates.Import(elbClient, importOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error importing certificate: %s", err)
	}
	// If all has been successful, set the ID on the resource
	d.SetId(c.CertificateId)
	log.Printf("[DEBUG] Imported certificate %s: %#v", d.Id(), importOpts)

	return resourceScmCertificateV3Read(d, meta)
}

func resourceScmCertificateV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	scmClient, err := config.ScmV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SCM client: %s", err)
	}
	certDetail, err := certificates.Get(scmClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error obtain certificate information")
	}
	log.Printf("[DEBUG] Retrieved certificate %s: %#v", d.Id(), certDetail)

	d.Set("region", GetRegion(d, config))
	d.Set("status", certDetail.Status)
	d.Set("orderId", certDetail.OrderId)
	d.Set("name", certDetail.Name)
	d.Set("certificate_type", certDetail.CertificateType)
	d.Set("brand", certDetail.Brand)
	d.Set("push_support", certDetail.PushSupport)
	d.Set("revoke_reason", certDetail.RevokeReason)
	d.Set("signature_algrithm", certDetail.SignatureAlgrithm)
	d.Set("issue_time", certDetail.IssueTime)
	d.Set("not_before", certDetail.NotBefore)
	d.Set("not_after", certDetail.NotAfter)

	d.Set("validity_period", certDetail.ValidityPeriod)
	d.Set("validation_method", certDetail.ValidationMethod)
	d.Set("domain_type", certDetail.DomainType)
	d.Set("domain", certDetail.Domain)
	d.Set("sans", certDetail.Sans)
	d.Set("domain_count", certDetail.DomainCount)
	d.Set("wildcard_count", certDetail.WildcardCount)

	// convert the type of 'certDetail.Authentifications' to TypeList
	auths := convertAuthToArray(certDetail.Authentifications)
	d.Set("authentifications", auths)

	return nil
}

func convertAuthToArray(authArr []certificates.Authentification) []map[string]interface{} {
	auths := make([]map[string]interface{}, 0, len(authArr))
	for _, v := range authArr {
		auth := map[string]interface{}{
			"record_name":  v.RecordName,
			"record_type":  v.RecordType,
			"record_value": v.RecordValue,
			"domain":       v.Domain,
		}
		auths = append(auths, auth)
	}
	return auths
}

func resourceScmCertificateV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	scmClient, err := config.ScmV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SCM client: %s", err)
	}

	log.Printf("[DEBUG] Deleting certificate: %s", d.Id())
	err = certificates.Delete(scmClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting certificate error %s: %s", d.Id(), err)
	}

	return nil
}
