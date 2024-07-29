package ccm

import (
	"context"
	"fmt"
	"log"

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
	targetServiceCdn = "CDN"
)

// @API CCM POST /v3/scm/certificates/import
// @API CCM POST /v3/scm/certificates/{certificate_id}/push
// @API CCM DELETE /v3/scm/certificates/{certificate_id}
// @API CCM GET /v3/scm/certificates/{certificate_id}
func ResourceCertificateImport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateImportCreate,
		UpdateContext: resourceCertificateImportUpdate,
		ReadContext:   resourceCertificateImportRead,
		DeleteContext: resourceCertificateImportDelete,
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
			"enc_certificate": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
			},
			"enc_private_key": {
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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

func buildPushOpts(targetService, targetProject string) certificates.PushOpts {
	return certificates.PushOpts{
		TargetService: targetService,
		TargetProject: targetProject,
	}
}

func pushCertificateToCDNService(client *golangsdk.ServiceClient, d *schema.ResourceData) map[string]interface{} {
	err := certificates.Push(client, d.Id(), buildPushOpts(targetServiceCdn, "")).ExtractErr()
	if err != nil {
		log.Printf("[WARN] error pushing certificate (%s) to CDN service: %s", d.Id(), err)
		return nil
	}

	return map[string]interface{}{
		"service": targetServiceCdn,
		"project": []string{},
	}
}

func pushCertificateToNonCDNService(client *golangsdk.ServiceClient, d *schema.ResourceData,
	targetMap map[string]interface{}) map[string]interface{} {
	var (
		service           = targetMap["service"].(string)
		projects          = targetMap["project"].([]interface{})
		projectAttributes = make([]string, 0, len(projects))
	)

	for _, project := range projects {
		if project.(string) == "" {
			log.Printf("[WARN] the argument `project` cannot be empty when pushing certificate to Non-CDN service (%s)", service)
			continue
		}

		err := certificates.Push(client, d.Id(), buildPushOpts(service, project.(string))).ExtractErr()
		if err != nil {
			log.Printf("[WARN] error pushing certificate (%s) to %s service: %s", d.Id(), service, err)
			continue
		}

		projectAttributes = append(projectAttributes, project.(string))
	}

	return map[string]interface{}{
		"service": service,
		"project": projectAttributes,
	}
}

func pushCertificateAndSetTargetAttribute(client *golangsdk.ServiceClient, d *schema.ResourceData) {
	var (
		targets          = d.Get("target").([]interface{})
		targetAttributes = make([]map[string]interface{}, 0, len(targets))
	)

	for _, target := range targets {
		targetMap := target.(map[string]interface{})
		if targetMap["service"].(string) == targetServiceCdn {
			targetAttributes = append(targetAttributes, pushCertificateToCDNService(client, d))
			continue
		}
		targetAttributes = append(targetAttributes, pushCertificateToNonCDNService(client, d, targetMap))
	}
	if err := d.Set("target", targetAttributes); err != nil {
		log.Printf("[ERROR] error setting `target` attribute to local state: %s", err)
	}
}

func resourceCertificateImportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.ScmV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	importOpts := certificates.ImportOpts{
		Name:                d.Get("name").(string),
		Certificate:         d.Get("certificate").(string),
		PrivateKey:          d.Get("private_key").(string),
		CertificateChain:    d.Get("certificate_chain").(string),
		EnterpriseProjectID: conf.GetEnterpriseProjectID(d),
		EncCertificate:      d.Get("enc_certificate").(string),
		EncPrivateKey:       d.Get("enc_private_key").(string),
	}
	c, err := certificates.Import(client, importOpts).Extract()
	if err != nil {
		return diag.Errorf("error importing CCM certificate: %s", err)
	}
	d.SetId(c.CertificateId)
	pushCertificateAndSetTargetAttribute(client, d)

	return resourceCertificateImportRead(ctx, d, meta)
}

func buildUpdatePushOpts(newPushCert, oldPushCert map[string]*schema.Set) ([]certificates.PushOpts, error) {
	pushOptResults := make([]certificates.PushOpts, 0, len(newPushCert))

	for service, newProjects := range newPushCert {
		oldProjects, ok := oldPushCert[service]
		if service == targetServiceCdn && !ok {
			pushOptResults = append(pushOptResults, buildPushOpts(targetServiceCdn, ""))
			continue
		}

		projectToAdd := newProjects
		if oldProjects != nil {
			projectToAdd = newProjects.Difference(oldProjects)
		}
		for _, project := range projectToAdd.List() {
			if project == "" {
				return nil, fmt.Errorf("the argument `project` cannot be empty when pushing certificate to Non-CDN service (%s)", service)
			}
			pushOptResults = append(pushOptResults, buildPushOpts(service, project.(string)))
		}
	}
	return pushOptResults, nil
}

func resourceCertificateImportUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.ScmV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}
	oldVal, newVal := d.GetChange("target")
	newPushCert, err := parsePushCertificateToMap(newVal.([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	oldPushCert, _ := parsePushCertificateToMap(oldVal.([]interface{}))

	// extract the new push service
	pushOptArrays, err := buildUpdatePushOpts(newPushCert, oldPushCert)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, pushOpts := range pushOptArrays {
		err := certificates.Push(client, d.Id(), pushOpts).ExtractErr()
		if err != nil {
			_ = d.Set("target", oldVal)
			return diag.Errorf("error pushing certificate (%s) to CDN service in update operation: %s", d.Id(), err)
		}
	}
	return resourceCertificateImportRead(ctx, d, meta)
}

func resourceCertificateImportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.ScmV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}
	certDetail, err := certificates.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error obtain certificate information")
	}

	auths := flattenAuthenticationAttribute(certDetail.Authentifications)
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
		d.Set("enterprise_project_id", certDetail.EnterpriseProjectID),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CCM certificate attributes: %s", err)
	}

	return nil
}

func flattenAuthenticationAttribute(authentications []certificates.Authentification) []map[string]interface{} {
	auth := make([]map[string]interface{}, 0, len(authentications))
	for _, v := range authentications {
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

func resourceCertificateImportDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.ScmV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	err = certificates.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting certificate")
	}

	return nil
}

// Parse the pushed service configuration into a map structure.
// The key indicates service name, and the value indicates the project set
func parsePushCertificateToMap(pushCertificate []interface{}) (map[string]*schema.Set, error) {
	serviceMapping := map[string]*schema.Set{}

	for _, pushInfo := range pushCertificate {
		targetService := pushInfo.(map[string]interface{})["service"].(string)
		targetProjectArr := pushInfo.(map[string]interface{})["project"].([]interface{})

		projects, ok := serviceMapping[targetService]
		if !ok {
			projects = schema.NewSet(schema.HashString, nil)
		}
		for _, proj := range targetProjectArr {
			projectName := proj.(string)
			if projects.Contains(projectName) {
				return nil, fmt.Errorf("duplicate project (%s) for the same service (%s)", projectName, targetService)
			}
			projects.Add(projectName)
		}
		serviceMapping[targetService] = projects
	}
	return serviceMapping, nil
}
