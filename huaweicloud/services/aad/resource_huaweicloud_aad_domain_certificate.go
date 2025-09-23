package aad

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var domainCertificateNonUpdatableParams = []string{
	"domain_id",
	"op_type",
	"cert_name",
	"cert_file",
	"cert_key_file",
}

// Due to limited testing conditions, the API in the creation method of this resource was not successfully called.

// @API AAD POST /v1/{project_id}/aad/external/domains/certificate
// @API AAD GET /v2/aad/domains/{domain_id}/certificate
func ResourceDomainCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCertificateCreate,
		ReadContext:   resourceDomainCertificateRead,
		UpdateContext: resourceDomainCertificateUpdate,
		DeleteContext: resourceDomainCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(domainCertificateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the domain ID.",
			},
			"op_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies the operation type. `0` represents upload, `1` represents replace.",
			},
			"cert_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the certificate name.",
			},
			"cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the certificate file content.",
			},
			"cert_key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the private key file content.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain name.",
			},
			"cert_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The certificate information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the certificate name.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate ID.",
						},
						"apply_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The applicable domain.",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The expiration time.",
						},
						"expire_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The expiration status.",
						},
					},
				},
			},
		},
	}
}

func buildDomainCertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"domain_id":     d.Get("domain_id"),
		"op_type":       d.Get("op_type"),
		"cert_name":     d.Get("cert_name"),
		"cert_file":     utils.ValueIgnoreEmpty(d.Get("cert_file")),
		"cert_key_file": utils.ValueIgnoreEmpty(d.Get("cert_key_file")),
	}
}

func resourceDomainCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v1/{project_id}/aad/external/domains/certificate"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildDomainCertificateBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AAD domain certificate: %s", err)
	}

	d.SetId(d.Get("domain_id").(string))

	return resourceDomainCertificateRead(ctx, d, meta)
}

func resourceDomainCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/domains/{domain_id}/certificate"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// When the `domain_id` does not exist, the API returns 400 and errCode is **AAD.50010037**, which needs to be
		// handled as 404 to be handled by CheckDeletedDiag.
		convertedErr := common.ConvertExpected400ErrInto404Err(err, "errCode", "AAD.50010037")
		return common.CheckDeletedDiag(d, convertedErr, "error retrieving AAD domain certificate")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certInfo := utils.PathSearch("cert_info", getRespBody, nil)

	mErr := multierror.Append(
		d.Set("domain_id", utils.PathSearch("domain_id", getRespBody, nil)),
		d.Set("cert_name", utils.PathSearch("cert_name", certInfo, nil)),
		d.Set("domain_name", utils.PathSearch("domain_name", getRespBody, nil)),
		d.Set("cert_info", flattenDomainCertificateInfo(certInfo)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDomainCertificateInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"cert_name":     utils.PathSearch("cert_name", resp, nil),
			"id":            utils.PathSearch("id", resp, nil),
			"apply_domain":  utils.PathSearch("apply_domain", resp, nil),
			"expire_time":   utils.PathSearch("expire_time", resp, nil),
			"expire_status": utils.PathSearch("expire_status", resp, nil),
		},
	}
}

func resourceDomainCertificateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDomainCertificateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to upload domain certificate. Deleting this resource
    will not change the current request record, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
