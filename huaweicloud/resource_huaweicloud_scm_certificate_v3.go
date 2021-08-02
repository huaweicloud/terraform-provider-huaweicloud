package huaweicloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/scm/v3/certificates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	maxErrorMessageLen = 200
	ellipsisString     = "..."

	targetServiceCdn        = "CDN"
	targetServiceWaf        = "WAF"
	targetServiceEnhanceElb = "Enhance_ELB"
)

func resourceScmCertificateV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceScmCertificateV3Create,
		Update: resourceScmCertificateV3Update,
		Read:   resourceScmCertificateV3Read,
		Delete: resourceScmCertificateV3Delete,

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
								targetServiceCdn, targetServiceWaf, targetServiceEnhanceElb,
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
		return fmtp.Errorf("Error creating HuaweiCloud SCM client: %s", err)
	}

	importOpts := certificates.ImportOpts{
		Name:             d.Get("name").(string),
		PrivateKey:       d.Get("private_key").(string),
		Certificate:      d.Get("certificate").(string),
		CertificateChain: d.Get("certificate_chain").(string),
	}

	c, err := certificates.Import(scmClient, importOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error importing certificate: %s", err)
	}
	// If all has been successful, set the ID on the resource
	d.SetId(c.CertificateId)
	logp.Printf("[DEBUG] Imported certificate %s: %#v", d.Id(), importOpts)

	// Get targets and push certificate to the target service.
	targets := d.Get("target").([]interface{})
	logp.Printf("[DEBUG] SCM Certificate target: %#v", targets)
	parseTargetsAndPush(scmClient, d, targets)

	return resourceScmCertificateV3Read(d, meta)
}

// parseTargetsAndPush pushes the certificate to the service.
// If the push fails, only the target attributes are updated.
func parseTargetsAndPush(c *golangsdk.ServiceClient, d *schema.ResourceData, targets []interface{}) {
	var tag = make([]map[string]interface{}, 0, len(targets))

	for _, pushInfo := range targets {
		service := pushInfo.(map[string]interface{})["service"].(string)
		projects := pushInfo.(map[string]interface{})["project"].([]interface{})

		t := map[string]interface{}{}

		if strings.Compare(service, targetServiceCdn) == 0 {
			pushOpts := certificates.PushOpts{
				TargetService: service,
			}
			logp.Printf("[DEBUG] Push certificate to CDN Service. %#v", pushOpts)
			err := pushCertificateToService(d.Id(), pushOpts, c)
			if err == nil {
				t["service"] = service
				t["project"] = []string{}
				tag = append(tag, t)
			} else {
				logp.Printf("[WARN] Push to CDN failed: %#v", pushOpts)
			}
		} else {
			var proj = make([]string, 0, len(projects))

			for _, p := range projects {
				pushOpts := certificates.PushOpts{
					TargetProject: p.(string),
					TargetService: service,
				}
				logp.Printf("[DEBUG] Push certificate to services. %#v", pushOpts)
				err := pushCertificateToService(d.Id(), pushOpts, c)
				if err == nil {
					proj = append(proj, p.(string))
				} else {
					logp.Printf("[WARN] Push failed: %#v", pushOpts)
				}
			}
			t["service"] = service
			t["project"] = proj
			tag = append(tag, t)
		}
	}
	d.Set("target", tag)
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
		if strings.Compare(service, targetServiceCdn) == 0 && !ok {
			pushOpts := certificates.PushOpts{
				TargetService: service,
			}
			logp.Printf("[DEBUG] Find new services and start to push. %#v", pushOpts)
			err := pushCertificateToService(d.Id(), pushOpts, scmClient)
			if err != nil {
				d.Set("target", oldVal)
				return err
			}
		} else {
			projectToAdd := newProjects
			if oldProjects != nil {
				projectToAdd = newProjects.Difference(oldProjects)
			}
			logp.Printf("[DEBUG] Find new services to push. %s: %#v", service, projectToAdd)
			for _, project := range projectToAdd.List() {
				pushOpts := certificates.PushOpts{
					TargetProject: project.(string),
					TargetService: service,
				}
				err := pushCertificateToService(d.Id(), pushOpts, scmClient)
				if err != nil {
					d.Set("target", oldVal)
					return err
				}
				logp.Printf("[DEBUG] Successfully push the certificate to the %s of %s.", service, project)
			}
		}
	}

	return resourceScmCertificateV3Read(d, meta)
}

func pushCertificateToService(id string, pushOpts certificates.PushOpts, scmClient *golangsdk.ServiceClient) error {
	if strings.Compare(pushOpts.TargetService, targetServiceCdn) != 0 && len(pushOpts.TargetProject) == 0 {
		return fmtp.Errorf("the argument of \"project\" cannot be empty, "+
			"it can be empty when pushed to the CDN service only. "+
			"\r\ncertificate_id: %s, service: %s", id, pushOpts.TargetService)
	}
	err := certificates.Push(scmClient, id, pushOpts).ExtractErr()
	if err != nil {
		// Parse 'err' to print more error messages.
		errMsg := processErr(err)
		return fmtp.Errorf(errMsg)
	}
	return nil
}

func resourceScmCertificateV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	scmClient, err := config.ScmV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud SCM client: %s", err)
	}
	certDetail, err := certificates.Get(scmClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error obtain certificate information")
	}
	logp.Printf("[DEBUG] Retrieved certificate %s: %#v", d.Id(), certDetail)

	d.Set("region", GetRegion(d, config))
	d.Set("status", certDetail.Status)
	d.Set("name", certDetail.Name)
	d.Set("push_support", certDetail.PushSupport)
	d.Set("not_before", certDetail.NotBefore)
	d.Set("not_after", certDetail.NotAfter)
	d.Set("domain", certDetail.Domain)
	d.Set("domain_count", certDetail.DomainCount)

	// convert the type of 'certDetail.Authentifications' to TypeList
	auths := buildAuthtificatesAttribute(certDetail.Authentifications)
	d.Set("authentifications", auths)

	return nil
}

func buildAuthtificatesAttribute(authentifications []certificates.Authentification) []map[string]interface{} {
	auth := make([]map[string]interface{}, 0, len(authentifications))
	for _, v := range authentifications {
		a := map[string]interface{}{
			"record_name":  v.RecordName,
			"record_type":  v.RecordType,
			"record_value": v.RecordValue,
			"domain":       v.Domain,
		}
		auth = append(auth, a)
	}
	return auth
}

func resourceScmCertificateV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	scmClient, err := config.ScmV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud SCM client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting certificate: %s", d.Id())
	err = certificates.Delete(scmClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting certificate error %s: %s", d.Id(), err)
	}

	return nil
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
				// if _, ok := projects[projectName]; ok {
				return nil, fmtp.Errorf("There are duplicate projects for the same service!\n"+
					"service = %s, project = %s.", targetService, projectName)
			}
			projects.Add(projectName)
		}
		serviceMapping[targetService] = projects

		logp.Printf("[DEBUG] Push certificate service mapping: %#v", serviceMapping)
	}
	return serviceMapping, nil
}

func processErr(err error) string {
	// errMsg: The error message to be printed.
	errMsg := fmt.Sprintf("Push certificate service error: %s", err)
	if err500, ok := err.(golangsdk.ErrDefault500); ok {
		errBody := string(err500.Body)
		// Maybe the text in the body is very long, only 200 characters printed。
		if len(errBody) >= maxErrorMessageLen {
			errBody = errBody[0:maxErrorMessageLen] + ellipsisString
		}
		// If 'err' is an ErrDefault500 object, the following information will be printed.
		logp.Printf("[ERROR] Push certificate service error. URL: %s, Body: %s",
			err500.URL, errBody)
		errMsg = fmt.Sprintf("Push certificate service error: "+
			"Bad request with: [%s %s], error message: %s", err500.Method, err500.URL, errBody)
	} else {
		// If 'err' is other error object, the default information will be printed.
		logp.Printf("[ERROR] Push certificate service error: %s, \n%#v", err.Error(), err)
		errMsg = fmt.Sprintf("Push certificate service error: %s", err)
	}
	return errMsg
}
