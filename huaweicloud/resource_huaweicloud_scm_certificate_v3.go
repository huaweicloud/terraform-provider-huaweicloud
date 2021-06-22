package huaweicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/scm/v3/certificates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	MAX_ERROR_MESSAGE_LEN = 200
	ELLIPSIS_STRING       = "..."

	TARGET_SERVICE_CDN         = "CDN"
	TARGET_SERVICE_WAF         = "WAF"
	TARGET_SERVICE_ENHANCE_ELB = "Enhance_ELB"
)

func resourceScmCertificateV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceScmCertificateV3Create,
		Update: resourceScmCertificateV3Update,
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
			"push_certificate": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_service": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								TARGET_SERVICE_CDN, TARGET_SERVICE_WAF, TARGET_SERVICE_ENHANCE_ELB,
							}, false),
						},
						"target_project": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					}},
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

	pushCert := d.Get("push_certificate").([]interface{})
	for _, pushInfo := range pushCert {
		targetService := pushInfo.(map[string]interface{})["target_service"].(string)
		targetProjectArr := pushInfo.(map[string]interface{})["target_project"].([]interface{})

		for _, targetProject := range targetProjectArr {
			pushOpts := certificates.PushOpts{
				TargetProject: targetProject.(string),
				TargetService: targetService,
			}
			log.Printf("[DEBUG] Push certificate to services. %#v", pushOpts)
			err := doPushCertificateToService(d.Id(), pushOpts, elbClient)
			if err != nil {
				return err
			}
		}
	}

	return resourceScmCertificateV3Read(d, meta)
}

func resourceScmCertificateV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ScmV3Client(GetRegion(d, config))

	oldVal, newVal := d.GetChange("push_certificate")
	newPushCert, err := parsePushCertificateToMap(newVal.([]interface{}))
	if err != nil {
		return err
	}

	oldPushCert, _ := parsePushCertificateToMap(oldVal.([]interface{}))

	// extract the new push service
	for targetService, newProjects := range newPushCert {
		oldProjects, ok := oldPushCert[targetService]
		if strings.Compare(targetService, TARGET_SERVICE_CDN) == 0 {
			if !ok {
				pushOpts := certificates.PushOpts{
					TargetService: targetService,
				}
				log.Printf("[DEBUG] Find new services and start to push. %#v", pushOpts)
				err := doPushCertificateToService(d.Id(), pushOpts, elbClient)
				if err != nil {
					return err
				}
			}
		} else if len(newProjects) == 0 {
			return fmt.Errorf("the argument of \"target_project\" cannot be empty, "+
				"it can be empty when pushed to the CDN service only. "+
				"\r\ncertificate_id: %s, target_service: %s", d.Id(), targetService)
		}
		for project, _ := range newProjects {
			if _, ok := oldProjects[project]; !ok {
				pushOpts := certificates.PushOpts{
					TargetProject: project,
					TargetService: targetService,
				}
				log.Printf("[DEBUG] Find new services and start to push. %#v", pushOpts)
				err := doPushCertificateToService(d.Id(), pushOpts, elbClient)
				if err != nil {
					return err
				}
			}
		}
	}

	//return fmt.Errorf("Error creating HuaweiCloud SCM client") //
	return resourceScmCertificateV3Read(d, meta)
}

func doPushCertificateToService(id string, pushOpts certificates.PushOpts, elbClient *golangsdk.ServiceClient) error {
	if strings.Compare(pushOpts.TargetService, TARGET_SERVICE_CDN) != 0 && len(pushOpts.TargetProject) == 0 {
		return fmt.Errorf("the argument of \"target_project\" cannot be empty, "+
			"it can be empty when pushed to the CDN service only. "+
			"\r\ncertificate_id: %s, target_service: %s", id, pushOpts.TargetService)
	}
	err := certificates.Push(elbClient, id, pushOpts).ExtractErr()
	if err != nil {
		// Parse 'err' to print more error messages.
		errMsg := processErr(err)
		return fmt.Errorf(errMsg)
	}
	return nil
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

func parsePushCertificateToMap(pushCertificate []interface{}) (map[string]map[string]bool, error) {
	serviceMapping := make(map[string]map[string]bool)

	for _, pushInfo := range pushCertificate {
		targetService := pushInfo.(map[string]interface{})["target_service"].(string)
		targetProjectArr := pushInfo.(map[string]interface{})["target_project"].([]interface{})

		projects, ok := serviceMapping[targetService]
		if !ok {
			projects = make(map[string]bool)
		}
		for _, proj := range targetProjectArr {
			projectName := proj.(string)
			if _, ok := projects[projectName]; ok {
				return nil, fmt.Errorf("There are duplicate projects for the same service!\n"+
					"target_service = %s, target_project = %s.", targetService, projectName)
			} else {
				projects[projectName] = true
			}
		}
		serviceMapping[targetService] = projects

		log.Printf("[DEBUG] Push certificate service mapping: %#v", serviceMapping)
	}
	return serviceMapping, nil
}
