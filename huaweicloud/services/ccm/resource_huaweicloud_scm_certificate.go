package ccm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/scm/v3/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	maxErrorMessageLen = 200
	ellipsisString     = "..."

	targetServiceCdn = "CDN"
)

// @API CCM POST /v3/scm/certificates/import
// @API CCM POST /v3/scm/certificates/{certificate_id}/push
// @API CCM DELETE /v3/scm/certificates/{certificate_id}
// @API CCM GET /v3/scm/certificates/{certificate_id}
func ResourceScmCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScmCertificateCreate,
		UpdateContext: resourceScmCertificateUpdate,
		ReadContext:   resourceScmCertificateRead,
		DeleteContext: resourceScmCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"certificate_chain": {
				Type:             schema.TypeString,
				Optional:         true,
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

func resourceScmCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	scmClient, err := conf.ScmV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SCM client: %s", err)
	}

	importOpts := certificates.ImportOpts{
		Name:             d.Get("name").(string),
		Certificate:      d.Get("certificate").(string),
		PrivateKey:       "***",
		CertificateChain: d.Get("certificate_chain").(string),
	}
	log.Printf("[DEBUG] Imported certificate options %s: %#v", d.Id(), importOpts)
	importOpts.PrivateKey = d.Get("private_key").(string)

	c, err := certificates.Import(scmClient, importOpts).Extract()
	if err != nil {
		return diag.Errorf("error importing certificate: %s", err)
	}
	// If all has been successful, set the ID on the resource
	d.SetId(c.CertificateId)

	// Get targets and push certificate to the target service.
	targets := d.Get("target").([]interface{})
	log.Printf("[DEBUG] SCM Certificate target: %#v", targets)
	parseTargetsAndPush(scmClient, d, targets)

	return resourceScmCertificateRead(ctx, d, meta)
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
			log.Printf("[DEBUG] Push certificate to CDN Service. %#v", pushOpts)
			err := pushCertificateToService(d.Id(), pushOpts, c)
			if err == nil {
				t["service"] = service
				t["project"] = []string{}
				tag = append(tag, t)
			} else {
				log.Printf("[WARN] Push to CDN failed: %#v", pushOpts)
			}
		} else {
			var proj = make([]string, 0, len(projects))

			for _, p := range projects {
				pushOpts := certificates.PushOpts{
					TargetProject: p.(string),
					TargetService: service,
				}
				log.Printf("[DEBUG] Push certificate to services. %#v", pushOpts)
				err := pushCertificateToService(d.Id(), pushOpts, c)
				if err == nil {
					proj = append(proj, p.(string))
				} else {
					log.Printf("[WARN] Push failed: %#v", pushOpts)
				}
			}
			t["service"] = service
			t["project"] = proj
			tag = append(tag, t)
		}
	}
	d.Set("target", tag)
}

func resourceScmCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	scmClient, err := conf.ScmV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SCM client: %s", err)
	}
	oldVal, newVal := d.GetChange("target")
	newPushCert, err := parsePushCertificateToMap(newVal.([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	oldPushCert, _ := parsePushCertificateToMap(oldVal.([]interface{}))

	// extract the new push service
	for service, newProjects := range newPushCert {
		oldProjects, ok := oldPushCert[service]
		if strings.Compare(service, targetServiceCdn) == 0 && !ok {
			pushOpts := certificates.PushOpts{
				TargetService: service,
			}
			log.Printf("[DEBUG] Find new services and start to push. %#v", pushOpts)
			err := pushCertificateToService(d.Id(), pushOpts, scmClient)
			if err != nil {
				d.Set("target", oldVal)
				return diag.FromErr(err)
			}
		} else {
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
				err := pushCertificateToService(d.Id(), pushOpts, scmClient)
				if err != nil {
					d.Set("target", oldVal)
					return diag.FromErr(err)
				}
				log.Printf("[DEBUG] Successfully push the certificate to the %s of %s.", service, project)
			}
		}
	}

	return resourceScmCertificateRead(ctx, d, meta)
}

func pushCertificateToService(id string, pushOpts certificates.PushOpts, scmClient *golangsdk.ServiceClient) error {
	if strings.Compare(pushOpts.TargetService, targetServiceCdn) != 0 && len(pushOpts.TargetProject) == 0 {
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

func resourceScmCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	scmClient, err := conf.ScmV3Client(region)
	if err != nil {
		return diag.Errorf("error creating SCM client: %s", err)
	}
	certDetail, err := certificates.Get(scmClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error obtain certificate information")
	}
	log.Printf("[DEBUG] Retrieved certificate %s: %#v", d.Id(), certDetail)

	// convert the type of 'certDetail.Authentifications' to TypeList
	auths := buildAuthtificatesAttribute(certDetail.Authentifications)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("status", certDetail.Status),
		d.Set("name", certDetail.Name),
		d.Set("push_support", certDetail.PushSupport),
		d.Set("authentifications", auths),
		d.Set("domain", certDetail.Domain),
		d.Set("domain_count", certDetail.DomainCount),
		d.Set("not_before", certDetail.NotBefore),
		d.Set("not_after", certDetail.NotAfter),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SCM certificate attributes: %s", err)
	}

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

func resourceScmCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	scmClient, err := conf.ScmV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SCM client: %s", err)
	}

	log.Printf("[DEBUG] Deleting certificate: %s", d.Id())
	err = certificates.Delete(scmClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting certificate")
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
				return nil, fmt.Errorf("there are duplicate projects for the same service, service = %s, project = %s", targetService, projectName)
			}
			projects.Add(projectName)
		}
		serviceMapping[targetService] = projects

		log.Printf("[DEBUG] Push certificate service mapping: %#v", serviceMapping)
	}
	return serviceMapping, nil
}

func processErr(err error) string {
	var errMsg string
	if err500, ok := err.(golangsdk.ErrDefault500); ok {
		errBody := string(err500.Body)
		// Maybe the text in the body is very long, only 200 characters printed。
		if len(errBody) >= maxErrorMessageLen {
			errBody = errBody[0:maxErrorMessageLen] + ellipsisString
		}
		// If 'err' is an ErrDefault500 object, the following information will be printed.
		log.Printf("[ERROR] Push certificate service error. URL: %s, Body: %s",
			err500.URL, errBody)
		errMsg = fmt.Sprintf("push certificate service error: "+
			"Bad request with: [%s %s], error message: %s", err500.Method, err500.URL, errBody)
	} else {
		// If 'err' is other error object, the default information will be printed.
		log.Printf("[ERROR] Push certificate service error: %s, \n%v", err.Error(), err)
		errMsg = fmt.Sprintf("push certificate service error: %s", err)
	}
	return errMsg
}
