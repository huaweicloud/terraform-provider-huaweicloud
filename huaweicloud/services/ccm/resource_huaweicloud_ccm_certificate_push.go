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
	"github.com/chnsz/golangsdk/openstack/elb/v3/certificates"
	wafcertificates "github.com/chnsz/golangsdk/openstack/waf/v1/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	targetServiceWaf = "WAF"
	targetServiceElb = "ELB"
)

// @API CCM POST /v3/scm/certificates/{certificate_id}/batch-push
// @API WAF GET /v1/{project_id}/waf/certificate/{certificate_id}
// @API WAF DELETE /v1/{project_id}/waf/certificate/{certificate_id}
// @API ELB GET /v3/{project_id}/elb/certificates/{certificate_id}
// @API ELB DELETE /v3/{project_id}/elb/certificates/{certificate_id}
func ResourceCcmCertificatePush() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCcmCertificateBatchPush,
		ReadContext:   resourceCcmCertificateRead,
		UpdateContext: resourceCcmCertificateUpdate,
		DeleteContext: resourceCcmCertificateDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"targets": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cert_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCcmCertificateBatchPush(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	scmClient, err := conf.ScmV3Client(conf.GetRegion(d))

	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}
	certId := d.Get("certificate_id").(string)
	pushCertificateHttpUrl := "v3/scm/certificates/{certificate_id}/batch-push"
	pushCertificatePath := scmClient.Endpoint + pushCertificateHttpUrl
	pushCertificatePath = strings.ReplaceAll(pushCertificatePath, "{certificate_id}", certId)

	pushCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	service := d.Get("service").(string)
	projects := d.Get("targets").(*schema.Set).List()
	pushCertificateOpt.JSONBody = utils.RemoveNil(buildPushCertificateBodyParams(service, projects))
	pushCertificateResp, err := scmClient.Request("POST", pushCertificatePath, &pushCertificateOpt)
	if err != nil {
		return diag.Errorf("error pushing certificate: %s", err)
	}
	pushCertificateRespBody, err := utils.FlattenResponse(pushCertificateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(certId)
	d.Set("targets", flattenResult(pushCertificateRespBody))

	return resourceCcmCertificateRead(ctx, d, meta)
}

func buildPushCertificateBodyParams(service string, projectlist []interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"target_service":  service,
		"target_projects": buildTargetProject(projectlist),
	}
	return bodyParams
}

func buildTargetProject(projectlist []interface{}) []string {
	var projects []string
	for _, v := range projectlist {
		projects = append(projects, utils.PathSearch("project_name", v, "").(string))
	}
	return projects
}

func flattenResult(resp interface{}) []interface{} {
	curJson := utils.PathSearch("results", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"project_name": utils.PathSearch("project_name", v, ""),
			"cert_id":      utils.PathSearch("cert_id", v, ""),
		})
	}
	return rst
}

func flattenUpdateResult(d *schema.ResourceData, resp interface{}) []interface{} {
	targets := d.Get("targets").(*schema.Set).List()
	rst := make([]interface{}, 0, len(targets))
	for _, v := range targets {
		certId := utils.PathSearch("cert_id", v, "").(string)
		projectName := utils.PathSearch("project_name", v, "").(string)
		if certId == "" {
			certId = getCertIdFromResult(projectName, resp)
		}
		rst = append(rst, map[string]interface{}{
			"project_name": projectName,
			"cert_id":      certId,
			"cert_name":    utils.PathSearch("cert_name", v, ""),
		})
	}

	return rst
}

func getCertIdFromResult(projectName string, resp interface{}) string {
	certId := ""
	curJson := utils.PathSearch("results", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		resName := utils.PathSearch("project_name", v, "").(string)
		if projectName == resName {
			certId = utils.PathSearch("cert_id", v, "").(string)
			break
		}
	}
	return certId
}
func resourceCcmCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectList := d.Get("targets").(*schema.Set).List()
	targetService := d.Get("service").(string)
	rst := make([]interface{}, 0, len(projectList))
	cfg := meta.(*config.Config)
	// search elb certificate detail and set to targets
	if targetService == targetServiceElb {
		for _, val := range projectList {
			projectName := utils.PathSearch("project_name", val, "").(string)
			certId := utils.PathSearch("cert_id", val, "").(string)
			elbClient, err := cfg.ElbV3Client(projectName)
			if err != nil {
				return diag.Errorf("error creating ELB client: %s", err)
			}

			cert, err := certificates.Get(elbClient, certId).Extract()
			if err != nil {
				log.Printf("[WARN] error retrieving ELB certificate: %s", certId)
			}
			if cert != nil {
				rst = append(rst, map[string]interface{}{
					"project_name": projectName,
					"cert_id":      certId,
					"cert_name":    cert.Name,
				})
			}
		}
	}

	// search waf certificate detail and set to targets
	if targetService == targetServiceWaf {
		for _, val := range projectList {
			projectName := utils.PathSearch("project_name", val, "").(string)
			certId := utils.PathSearch("cert_id", val, "").(string)
			wafClient, err := cfg.WafV1Client(projectName)
			if err != nil {
				return diag.Errorf("error creating WAF client: %s", err)
			}

			cert, err := wafcertificates.Get(wafClient, certId).Extract()
			if err != nil {
				log.Printf("[DEBUG] error retrieving WAF certificate: %s", certId)
			}
			if cert != nil && cert.Id != "" {
				rst = append(rst, map[string]interface{}{
					"project_name": projectName,
					"cert_id":      certId,
					"cert_name":    cert.Name,
				})
			}
		}
	}
	if len(rst) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving certificate")
	}
	mErr := multierror.Append(nil,
		d.Set("targets", rst),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting targets fields: %s", err)
	}
	return nil
}
func resourceCcmCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("targets") {
		err := rePushOrDelete(d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCcmCertificateRead(ctx, d, meta)
}

func rePushOrDelete(d *schema.ResourceData, meta interface{}) error {
	oRaw, nRaw := d.GetChange("targets")
	addSet := nRaw.(*schema.Set).Difference(oRaw.(*schema.Set))
	rmSet := oRaw.(*schema.Set).Difference(nRaw.(*schema.Set))
	if rmSet.Len() > 0 {
		// remve the pushed certificate
		targetService := d.Get("service").(string)
		_, err := rmPushedCert(targetService, rmSet.List(), meta, d)
		if err != nil {
			return err
		}
	}
	if addSet.Len() > 0 {
		// push to new project
		conf := meta.(*config.Config)
		scmClient, err := conf.ScmV3Client(conf.GetRegion(d))

		if err != nil {
			return err
		}

		pushCertificateHttpUrl := "v3/scm/certificates/{certificate_id}/batch-push"
		pushCertificatePath := scmClient.Endpoint + pushCertificateHttpUrl
		pushCertificatePath = strings.ReplaceAll(pushCertificatePath, "{certificate_id}", d.Get("certificate_id").(string))

		pushCertificateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}
		pushCertificateOpt.JSONBody = utils.RemoveNil(buildPushCertificateBodyParams(d.Get("service").(string), addSet.List()))
		pushCertificateResp, err := scmClient.Request("POST", pushCertificatePath, &pushCertificateOpt)
		if err != nil {
			return fmt.Errorf("error pushing certificate: %s", err)
		}
		pushCertificateRespBody, err := utils.FlattenResponse(pushCertificateResp)
		if err != nil {
			return err
		}
		// save push result
		d.Set("targets", flattenUpdateResult(d, pushCertificateRespBody))
	}
	return nil
}

func resourceCcmCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	targetService := d.Get("service").(string)
	oRaws := d.Get("targets").(*schema.Set).List()
	if len(oRaws) == 0 {
		return nil
	}
	_, err := rmPushedCert(targetService, oRaws, meta, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func rmPushedCert(targetService string, oRaws []interface{}, meta interface{}, d *schema.ResourceData) (bool, error) {
	cfg := meta.(*config.Config)
	if targetService == targetServiceElb {
		for _, val := range oRaws {
			strVal := utils.PathSearch("project_name", val, "").(string)
			certId := utils.PathSearch("cert_id", val, "").(string)
			if certId != "" {
				elbClient, err := cfg.ElbV3Client(strVal)
				if err != nil {
					return true, fmt.Errorf("error creating ELB client: %s", err)
				}
				err = certificates.Delete(elbClient, certId).ExtractErr()
				if err != nil {
					if utils.IsResourceNotFound(err) {
						return true, nil
					}
					return true, fmt.Errorf("error deleting certificate %s: %s", certId, err)
				}
			}
		}
	}

	if targetService == targetServiceWaf {
		for _, val := range oRaws {
			strVal := utils.PathSearch("project_name", val, "").(string)
			certId := utils.PathSearch("cert_id", val, "").(string)
			if certId != "" {
				client, err := cfg.WafV1Client(strVal)
				if err != nil {
					return true, fmt.Errorf("error creating WAF client: %s", err)
				}
				epsID := cfg.GetEnterpriseProjectID(d)
				err = wafcertificates.DeleteWithEpsID(client, certId, epsID).ExtractErr()
				if err != nil {
					if utils.IsResourceNotFound(err) {
						return true, nil
					}
					return true, fmt.Errorf("error deleting WAF certificate: %s", err)
				}
			}
		}
	}
	return false, nil
}
