package elb

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB POST /v3/{project_id}/elb/certificates
// @API ELB GET /v3/{project_id}/elb/certificates/{certificate_id}
// @API ELB PUT /v3/{project_id}/elb/certificates/{certificate_id}
// @API ELB DELETE /v3/{project_id}/elb/certificates/{certificate_id}
func ResourceCertificateV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateV3Create,
		ReadContext:   resourceCertificateV3Read,
		UpdateContext: resourceCertificateV3Update,
		DeleteContext: resourceCertificateV3Delete,
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
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
			"enc_certificate": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
			},
			"enc_private_key": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
			},
			"scm_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"common_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_alternative_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCertificateV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/certificates"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateCertificateBodyParams(d, cfg))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB certificate: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error retrieving ELB certificate: %s", err)
	}
	certificateId := utils.PathSearch("certificate.id", createRespBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("error creating ELB certificate: ID is not found in API response")
	}

	d.SetId(certificateId)

	return resourceCertificateV3Read(ctx, d, meta)
}

func buildCreateCertificateBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"type":                  utils.ValueIgnoreEmpty(d.Get("type")),
		"domain":                utils.ValueIgnoreEmpty(d.Get("domain")),
		"private_key":           utils.ValueIgnoreEmpty(d.Get("private_key")),
		"certificate":           utils.ValueIgnoreEmpty(d.Get("certificate")),
		"enc_certificate":       utils.ValueIgnoreEmpty(d.Get("enc_certificate")),
		"enc_private_key":       utils.ValueIgnoreEmpty(d.Get("enc_private_key")),
		"scm_certificate_id":    utils.ValueIgnoreEmpty(d.Get("scm_certificate_id")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	bodyParams := map[string]interface{}{
		"certificate": params,
	}
	return bodyParams
}

func resourceCertificateV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/elb/certificates/{certificate_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{certificate_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB certificate")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("certificate.name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("certificate.description", getRespBody, nil)),
		d.Set("type", utils.PathSearch("certificate.type", getRespBody, nil)),
		d.Set("domain", utils.PathSearch("certificate.domain", getRespBody, nil)),
		d.Set("certificate", utils.PathSearch("certificate.certificate", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("certificate.create_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("certificate.update_time", getRespBody, nil)),
		d.Set("expire_time", utils.PathSearch("certificate.expire_time", getRespBody, nil)),
		d.Set("enc_certificate", utils.PathSearch("certificate.enc_certificate", getRespBody, nil)),
		d.Set("scm_certificate_id", utils.PathSearch("certificate.scm_certificate_id", getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("certificate.enterprise_project_id", getRespBody, nil)),
		d.Set("common_name", utils.PathSearch("certificate.common_name", getRespBody, nil)),
		d.Set("fingerprint", utils.PathSearch("certificate.fingerprint", getRespBody, nil)),
		d.Set("subject_alternative_names", utils.PathSearch("certificate.subject_alternative_names",
			getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCertificateV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/certificates/{certificate_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{certificate_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateCertificateBodyParams(d))
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating ELB certificate: %s", err)
	}

	return resourceCertificateV3Read(ctx, d, meta)
}

func buildUpdateCertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":               utils.ValueIgnoreEmpty(d.Get("name")),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
		"domain":             utils.ValueIgnoreEmpty(d.Get("domain")),
		"private_key":        utils.ValueIgnoreEmpty(d.Get("private_key")),
		"certificate":        utils.ValueIgnoreEmpty(d.Get("certificate")),
		"enc_certificate":    utils.ValueIgnoreEmpty(d.Get("enc_certificate")),
		"enc_private_key":    utils.ValueIgnoreEmpty(d.Get("enc_private_key")),
		"scm_certificate_id": utils.ValueIgnoreEmpty(d.Get("scm_certificate_id")),
	}
	bodyParams := map[string]interface{}{
		"certificate": params,
	}
	return bodyParams
}

func resourceCertificateV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/certificates/{certificate_id}"
		product = "elb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{certificate_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting ELB certificate")
	}

	return nil
}
