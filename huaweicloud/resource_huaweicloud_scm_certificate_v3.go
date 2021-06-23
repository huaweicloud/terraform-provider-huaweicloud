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
			"target": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								TARGET_SERVICE_CDN, TARGET_SERVICE_WAF, TARGET_SERVICE_ENHANCE_ELB,
							}, false),
						},
						"project": {
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
			"domain_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"push_support": {
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
	scmClient, err := config.ScmV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SCM client: %s", err)
	}

	importOpts := certificates.ImportOpts{
		Name:             d.Get("name").(string),
		PrivateKey:       d.Get("private_key").(string),
		Certificate:      d.Get("certificate").(string),
		CertificateChain: d.Get("certificate_chain").(string),
	}

	c, err := certificates.Import(scmClient, importOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error importing certificate: %s", err)
	}
	// If all has been successful, set the ID on the resource
	d.SetId(c.CertificateId)
	log.Printf("[DEBUG] Imported certificate %s: %#v", d.Id(), importOpts)

	pushCert := d.Get("target").([]interface{})
	for _, pushInfo := range pushCert {
		targetService := pushInfo.(map[string]interface{})["service"].(string)
		targetProjectArr := pushInfo.(map[string]interface{})["project"].([]interface{})

		for _, targetProject := range targetProjectArr {
			pushOpts := certificates.PushOpts{
				TargetProject: targetProject.(string),
				TargetService: targetService,
			}
			log.Printf("[DEBUG] Push certificate to services. %#v", pushOpts)
			err := doPushCertificateToService(d.Id(), pushOpts, scmClient)
			if err != nil {
				d.SetId("")
				return err
			}
		}
	}

	return resourceScmCertificateV3Read(d, meta)
}

func resourceScmCertificateV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	scmClient, err := config.ScmV3Client(GetRegion(d, config))

	oldVal, newVal := d.GetChange("target")
	newPushCert, err := parsePushCertificateToMap(newVal.([]interface{}))
	if err != nil {
		return err
	}
	oldPushCert, _ := parsePushCertificateToMap(oldVal.([]interface{}))

	// extract the new push service
	for service, newProjects := range newPushCert {
		oldProjects, ok := oldPushCert[service]
		if strings.Compare(service, TARGET_SERVICE_CDN) == 0 && !ok {
			pushOpts := certificates.PushOpts{
				TargetService: service,
			}
			log.Printf("[DEBUG] Find new services and start to push. %#v", pushOpts)
			err := doPushCertificateToService(d.Id(), pushOpts, scmClient)
			if err != nil {
				d.SetId("")
				return err
			}
		} else {
			log.Printf("[DEBUG] Difference newProjects. %s, %#v", service, newProjects)
			log.Printf("[DEBUG] Difference oldProjects. %s, %#v", service, oldProjects)
			projectToAdd := newProjects
			if oldProjects != nil {
				projectToAdd = newProjects.Difference(oldProjects)
			}
			log.Printf("[DEBUG] Find new services to push. %s: %#v", service, projectToAdd)
			for _, project := range projectToAdd.List() {
				pushOpts := certificates.PushOpts{
					TargetProject: project.(string),
					TargetService: service,
				}
				err := doPushCertificateToService(d.Id(), pushOpts, scmClient)
				if err != nil {
					d.SetId("")
					return err
				}
				log.Printf("[DEBUG] Successfully push the certificate to the %s of %s.", service, project)
			}
		}
	}

	return resourceScmCertificateV3Read(d, meta)
}

func doPushCertificateToService(id string, pushOpts certificates.PushOpts, scmClient *golangsdk.ServiceClient) error {
	if strings.Compare(pushOpts.TargetService, TARGET_SERVICE_CDN) != 0 && len(pushOpts.TargetProject) == 0 {
		return fmt.Errorf("the argument of \"project\" cannot be empty, "+
			"it can be empty when pushed to the CDN service only. "+
			"\r\ncertificate_id: %s, service: %s", id, pushOpts.TargetService)
	}
	err := certificates.Push(scmClient, id, pushOpts).ExtractErr()
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
	d.Set("name", certDetail.Name)
	d.Set("push_support", certDetail.PushSupport)
	d.Set("not_before", certDetail.NotBefore)
	d.Set("not_after", certDetail.NotAfter)

	d.Set("domain", certDetail.Domain)
	d.Set("domain_count", certDetail.DomainCount)

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

func parsePushCertificateToMap(pushCertificate []interface{}) (map[string]*schema.Set, error) {
	serviceMapping := map[string]*schema.Set{}

	for _, pushInfo := range pushCertificate {
		targetService := pushInfo.(map[string]interface{})["service"].(string)
		targetProjectArr := pushInfo.(map[string]interface{})["project"].([]interface{})

		projects, ok := serviceMapping[targetService]
		if !ok {
			projects = &schema.Set{F: schema.HashString}
		}
		for _, proj := range targetProjectArr {
			projectName := proj.(string)
			if projects.Contains(projectName) {
				//if _, ok := projects[projectName]; ok {
				return nil, fmt.Errorf("There are duplicate projects for the same service!\n"+
					"service = %s, project = %s.", targetService, projectName)
			} else {
				projects.Add(projectName)
			}
		}
		serviceMapping[targetService] = projects

		log.Printf("[DEBUG] Push certificate service mapping: %#v", serviceMapping)
	}
	return serviceMapping, nil
}

func processErr(err error) string {
	// errMsg: The error message to be printed.
	errMsg := fmt.Sprintf("Push certificate service error: %s", err)
	if err500, ok := err.(golangsdk.ErrDefault500); ok {
		errBody := string(err500.Body)
		// Maybe the text in the body is very long, only 200 characters printed。
		if len(errBody) >= MAX_ERROR_MESSAGE_LEN {
			errBody = errBody[0:MAX_ERROR_MESSAGE_LEN] + ELLIPSIS_STRING
		}
		// If 'err' is an ErrDefault500 object, the following information will be printed.
		log.Printf("[ERROR] Push certificate service error. URL: %s, Body: %s",
			err500.URL, errBody)
		errMsg = fmt.Sprintf("Push certificate service error: "+
			"Bad request with: [%s %s], error message: %s", err500.Method, err500.URL, errBody)
	} else {
		// If 'err' is other error object, the default information will be printed.
		log.Printf("[ERROR] Push certificate service error: %s, \n%#v", err.Error(), err)
		errMsg = fmt.Sprintf("Push certificate service error: %s", err)
	}
	return errMsg
}
