package lb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB POST /v2/{project_id}/elb/certificates
// @API ELB GET /v2/{project_id}/elb/certificates/{certificate_id}
// @API ELB PUT /v2/{project_id}/elb/certificates/{certificate_id}
// @API ELB DELETE /v2/{project_id}/elb/certificates/{certificate_id}
func ResourceCertificateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateV2Create,
		ReadContext:   resourceCertificateV2Read,
		UpdateContext: resourceCertificateV2Update,
		DeleteContext: resourceCertificateV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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
				Sensitive:        true,
			},
			"certificate": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: utils.SuppressNewLineDiffs,
				Sensitive:        true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
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
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCertificateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/elb/certificates"
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
	certificateId := utils.PathSearch("id", createRespBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("error creating ELB certificate: ID is not found in API response")
	}

	d.SetId(certificateId)

	return resourceCertificateV2Read(ctx, d, meta)
}

func buildCreateCertificateBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"type":                  utils.ValueIgnoreEmpty(d.Get("type")),
		"domain":                utils.ValueIgnoreEmpty(d.Get("domain")),
		"private_key":           utils.ValueIgnoreEmpty(d.Get("private_key")),
		"certificate":           utils.ValueIgnoreEmpty(d.Get("certificate")),
		"protection_reason":     utils.ValueIgnoreEmpty(d.Get("protection_reason")),
		"protection_status":     utils.ValueIgnoreEmpty(d.Get("protection_status")),
		"source":                utils.ValueIgnoreEmpty(d.Get("source")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	return bodyParams
}

func resourceCertificateV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/{project_id}/elb/certificates/{certificate_id}"
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
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("domain", utils.PathSearch("domain", getRespBody, nil)),
		d.Set("certificate", utils.PathSearch("certificate", getRespBody, nil)),
		d.Set("protection_reason", utils.PathSearch("protection_reason", getRespBody, nil)),
		d.Set("protection_status", utils.PathSearch("protection_status", getRespBody, nil)),
		d.Set("source", utils.PathSearch("source", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
		d.Set("expire_time", utils.PathSearch("expire_time", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCertificateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/elb/certificates/{certificate_id}"
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

	return resourceCertificateV2Read(ctx, d, meta)
}

func buildUpdateCertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":              d.Get("name"),
		"description":       d.Get("description"),
		"domain":            d.Get("domain"),
		"private_key":       d.Get("private_key"),
		"certificate":       d.Get("certificate"),
		"protection_reason": d.Get("protection_reason"),
		"protection_status": d.Get("protection_status"),
		"source":            d.Get("source"),
	}
	return bodyParams
}

func resourceCertificateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/elb/certificates/{certificate_id}"
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
