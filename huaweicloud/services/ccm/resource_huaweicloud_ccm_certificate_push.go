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

// @API CCM POST /v3/scm/certificates/{certificate_id}/batch-push
// @API WAF GET /v1/{project_id}/waf/certificate/{certificate_id}
// @API WAF DELETE /v1/{project_id}/waf/certificate/{certificate_id}
// @API ELB GET /v3/{project_id}/elb/certificates/{certificate_id}
// @API ELB DELETE /v3/{project_id}/elb/certificates/{certificate_id}
func ResourceCertificatePush() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificatePushCreate,
		ReadContext:   resourceCertificatePushRead,
		UpdateContext: resourceCertificatePushUpdate,
		DeleteContext: resourceCertificatePushDelete,

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

func buildTargetProject(projectList []interface{}) []string {
	projects := make([]string, 0, len(projectList))
	for _, v := range projectList {
		projectName := utils.PathSearch("project_name", v, "").(string)
		if projectName == "" {
			continue
		}

		projects = append(projects, projectName)
	}
	return projects
}

func buildPushCertificateBodyParams(service string, projectList []interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"target_service":  service,
		"target_projects": buildTargetProject(projectList),
	}
	return bodyParams
}

func pushCertificateToService(scmClient *golangsdk.ServiceClient, d *schema.ResourceData, projects []interface{}) (interface{}, error) {
	var (
		pushCertificateHttpUrl = "v3/scm/certificates/{certificate_id}/batch-push"
		certId                 = d.Get("certificate_id").(string)
		service                = d.Get("service").(string)
	)

	pushCertificatePath := scmClient.Endpoint + pushCertificateHttpUrl
	pushCertificatePath = strings.ReplaceAll(pushCertificatePath, "{certificate_id}", certId)
	pushCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildPushCertificateBodyParams(service, projects)),
	}

	pushCertificateResp, err := scmClient.Request("POST", pushCertificatePath, &pushCertificateOpt)
	if err != nil {
		return nil, fmt.Errorf("error pushing certificate: %s", err)
	}
	return utils.FlattenResponse(pushCertificateResp)
}

// flattenCreateResult using to obtain the target certificate ID from the response of the creation API and record it.
// The response body example is as follows: `{"results":[{"project_name":"cn-north-4","cert_id":"XXX","message":"success"},
// {"project_name":"cn-error-4","cert_id":"","message":"Invalid region id."}]}`
// If `cert_id` is empty, it means that the certificate push failed.
func flattenCreateResult(resp interface{}, d *schema.ResourceData) ([]interface{}, error) {
	curJson := utils.PathSearch("results", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	var mErr *multierror.Error
	for _, v := range curArray {
		projectName := utils.PathSearch("project_name", v, "").(string)
		certId := utils.PathSearch("cert_id", v, "").(string)
		if certId == "" {
			message := utils.PathSearch("message", v, "").(string)
			mErr = multierror.Append(
				mErr,
				fmt.Errorf("error pushing certificate (%s) to project (%s): %s", d.Id(), projectName, message),
			)
			continue
		}

		rst = append(rst, map[string]interface{}{
			"project_name": projectName,
			"cert_id":      certId,
		})
	}

	return rst, mErr.ErrorOrNil()
}

func resourceCertificatePushCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		certId = d.Get("certificate_id").(string)
	)

	scmClient, err := cfg.ScmV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	pushCertificateRespBody, err := pushCertificateToService(scmClient, d, d.Get("targets").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(certId)

	rstTargets, err := flattenCreateResult(pushCertificateRespBody, d)
	if len(rstTargets) == 0 {
		return diag.Errorf("all attempts to push the certificate to the services failed in creation: %s.\n "+
			"Please check whether the push service and projects are accurate.", err)
	}

	if err != nil {
		log.Printf("[WARM] some error occurred when pushing certificate to services in creation: %s", err)
	}

	if err := d.Set("targets", rstTargets); err != nil {
		return diag.Errorf("error setting `targets` attributes in creation operation: %s", err)
	}

	return resourceCertificatePushRead(ctx, d, meta)
}

func FlattenElbCertificateDetail(cfg *config.Config, projectName string, certId string) map[string]interface{} {
	elbClient, err := cfg.ElbV3Client(projectName)
	if err != nil {
		log.Printf("[WARN] error creating ELB client: %s", err)
		return nil
	}

	cert, err := certificates.Get(elbClient, certId).Extract()
	if err != nil {
		log.Printf("[WARN] error retrieving ELB certificate (%s): %s", certId, err)
		return nil
	}

	if cert == nil {
		log.Printf("[WARN] error retrieving ELB certificate (%s): The certificate in API response is empty", certId)
		return nil
	}
	return map[string]interface{}{
		"project_name": projectName,
		"cert_id":      certId,
		"cert_name":    cert.Name,
	}
}

func FlattenWafCertificateDetail(cfg *config.Config, projectName string, certId string) map[string]interface{} {
	wafClient, err := cfg.WafV1Client(projectName)
	if err != nil {
		log.Printf("[WARN] error creating WAF client: %s", err)
		return nil
	}

	cert, err := wafcertificates.Get(wafClient, certId).Extract()
	if err != nil {
		log.Printf("[WARN] error retrieving WAF certificate (%s): %s", certId, err)
		return nil
	}

	if cert == nil {
		log.Printf("[WARN] error retrieving WAF certificate (%s): The certificate in API response is empty", certId)
		return nil
	}
	return map[string]interface{}{
		"project_name": projectName,
		"cert_id":      certId,
		"cert_name":    cert.Name,
	}
}

func resourceCertificatePushRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		projectList   = d.Get("targets").(*schema.Set).List()
		targetService = d.Get("service").(string)
		rst           = make([]interface{}, 0, len(projectList))
		cfg           = meta.(*config.Config)
	)

	for _, val := range projectList {
		projectName := utils.PathSearch("project_name", val, "").(string)
		certId := utils.PathSearch("cert_id", val, "").(string)
		if certId == "" {
			continue
		}

		switch targetService {
		case "ELB":
			if elbRstMap := FlattenElbCertificateDetail(cfg, projectName, certId); elbRstMap != nil {
				rst = append(rst, elbRstMap)
			}
		case "WAF":
			if wafRstMap := FlattenWafCertificateDetail(cfg, projectName, certId); wafRstMap != nil {
				rst = append(rst, wafRstMap)
			}
		}
	}

	if len(rst) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving certificate")
	}
	mErr := multierror.Append(nil,
		d.Set("targets", rst),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func removePushedCertificate(cfg *config.Config, oRaws []interface{}, d *schema.ResourceData) error {
	var mErr *multierror.Error
	for _, val := range oRaws {
		projectName := utils.PathSearch("project_name", val, "").(string)
		certId := utils.PathSearch("cert_id", val, "").(string)
		if certId == "" {
			continue
		}

		switch d.Get("service").(string) {
		case "ELB":
			elbClient, err := cfg.ElbV3Client(projectName)
			if err != nil {
				log.Printf("[WARN] error creating ELB client: %s", err)
				continue
			}

			err = certificates.Delete(elbClient, certId).ExtractErr()
			if err != nil && !utils.IsResourceNotFound(err) {
				mErr = multierror.Append(mErr, fmt.Errorf("error deleting ELB certificate (%s): %s", certId, err))
			}
		case "WAF":
			client, err := cfg.WafV1Client(projectName)
			if err != nil {
				log.Printf("[WARN] error creating WAF client: %s", err)
				continue
			}
			epsID := cfg.GetEnterpriseProjectID(d)
			err = wafcertificates.DeleteWithEpsID(client, certId, epsID).ExtractErr()
			if err != nil && !utils.IsResourceNotFound(err) {
				mErr = multierror.Append(mErr, fmt.Errorf("error deleting WAF certificate (%s): %s", certId, err))
			}
		}
	}

	return mErr.ErrorOrNil()
}

// flattenUpdateResult using to handle response values in the update API.
// The response body example is as follows: `{"results":[{"project_name":"cn-north-4","cert_id":"XXX","message":"success"},
// {"project_name":"cn-error-4","cert_id":"","message":"Invalid region id."}]}`
// 1. Try to get the object from the API response through `project_name`.
// 2. If it cannot be obtained, directly set the value in `local state`.
// 3. If the object can be obtained from the API response, try to obtain whether the `cert_id` has a value.
// 4. If the `cert_id` in the API response has a value, get and set values from response.
// 5. Otherwise, it will be considered that the push failed, the failed push operation is saved in the error message.
func flattenUpdateResult(d *schema.ResourceData, resp interface{}) ([]interface{}, error) {
	targets := d.Get("targets").(*schema.Set).List()
	rst := make([]interface{}, 0, len(targets))
	var mErr *multierror.Error
	for _, v := range targets {
		projectName := utils.PathSearch("project_name", v, "").(string)
		if projectName == "" {
			continue
		}

		publishResult := utils.PathSearch(fmt.Sprintf("results[?project_name=='%s']|[0]", projectName), resp, nil)
		if publishResult == nil {
			rst = append(rst, v)
			continue
		}

		certId := utils.PathSearch("cert_id", publishResult, "").(string)
		if certId != "" {
			rst = append(rst, map[string]interface{}{
				"project_name": projectName,
				"cert_id":      certId,
			})
			continue
		}

		message := utils.PathSearch("message", publishResult, "").(string)
		mErr = multierror.Append(
			mErr,
			fmt.Errorf("error pushing certificate (%s) to project (%s): %s", d.Id(), projectName, message),
		)
	}

	return rst, mErr.ErrorOrNil()
}

func resourceCertificatePushUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	if d.HasChange("targets") {
		oRaw, nRaw := d.GetChange("targets")
		addSet := nRaw.(*schema.Set).Difference(oRaw.(*schema.Set))
		rmSet := oRaw.(*schema.Set).Difference(nRaw.(*schema.Set))
		if rmSet.Len() > 0 {
			err := removePushedCertificate(cfg, rmSet.List(), d)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if addSet.Len() > 0 {
			scmClient, err := cfg.ScmV3Client(region)
			if err != nil {
				return diag.Errorf("error creating CCM client: %s", err)
			}
			pushCertificateRespBody, err := pushCertificateToService(scmClient, d, addSet.List())
			if err != nil {
				return diag.FromErr(err)
			}

			rstTargets, err := flattenUpdateResult(d, pushCertificateRespBody)
			if len(rstTargets) == 0 {
				return diag.Errorf("all attempts to push the certificate to the services failed in update"+
					" operation: %s.\n Please check whether the push service and projects are accurate.", err)
			}

			if err != nil {
				log.Printf("[WARM] some error occurred when pushing certificate to services in update operation: %s", err)
			}

			if err := d.Set("targets", rstTargets); err != nil {
				return diag.Errorf("error setting `targets` attributes in update operation: %s", err)
			}
		}
	}

	return resourceCertificatePushRead(ctx, d, meta)
}

func resourceCertificatePushDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg   = meta.(*config.Config)
		oRaws = d.Get("targets").(*schema.Set).List()
	)

	if len(oRaws) == 0 {
		return nil
	}
	return diag.FromErr(removePushedCertificate(cfg, oRaws, d))
}
