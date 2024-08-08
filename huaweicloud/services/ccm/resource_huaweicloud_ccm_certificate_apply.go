package ccm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SCM POST /v3/scm/certificates/{certificate_id}/apply
func ResourceCertificateApply() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateApplyCreate,
		ReadContext:   resourceCertificateApplyRead,
		DeleteContext: resourceCertificateApplyDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the CCM SSL certificate ID.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the domain name bound to the certificate.`,
			},
			"applicant_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the applicant.`,
			},
			"applicant_phone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the phone number of the applicant.`,
			},
			"applicant_email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the email of the applicant.`,
			},
			"domain_method": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the domain name verification method.`,
			},
			"sans": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies additional domain names bound to multi-domain type certificates.`,
			},
			"csr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the certificate CSR string, which must match the domain name.`,
			},
			"company_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the company name.`,
			},
			"company_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the department name.`,
			},
			"company_province": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the province where the company is located.`,
			},
			"company_city": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the city where the company is located.`,
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the country code.`,
			},
			"contact_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the technical contact name.`,
			},
			"contact_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the technical contact phone number.`,
			},
			"contact_email": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the technical contact email.`,
			},
			"auto_dns_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to push DNS verification information to HuaweiCloud resolution service.`,
			},
			"key_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the key algorithm.`,
			},
			"ca_hash_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the signature algorithm.`,
			},
		},
	}
}

func buildCertificateApplyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":                   d.Get("domain"),
		"applicant_name":           d.Get("applicant_name"),
		"applicant_phone":          d.Get("applicant_phone"),
		"applicant_email":          d.Get("applicant_email"),
		"domain_method":            d.Get("domain_method"),
		"sans":                     utils.ValueIgnoreEmpty(d.Get("sans")),
		"csr":                      utils.ValueIgnoreEmpty(d.Get("csr")),
		"company_name":             utils.ValueIgnoreEmpty(d.Get("company_name")),
		"company_unit":             utils.ValueIgnoreEmpty(d.Get("company_unit")),
		"company_province":         utils.ValueIgnoreEmpty(d.Get("company_province")),
		"company_city":             utils.ValueIgnoreEmpty(d.Get("company_city")),
		"country":                  utils.ValueIgnoreEmpty(d.Get("country")),
		"contact_name":             utils.ValueIgnoreEmpty(d.Get("contact_name")),
		"contact_phone":            utils.ValueIgnoreEmpty(d.Get("contact_phone")),
		"contact_email":            utils.ValueIgnoreEmpty(d.Get("contact_email")),
		"auto_dns_auth":            utils.ValueIgnoreEmpty(d.Get("auto_dns_auth")),
		"key_algorithm":            utils.ValueIgnoreEmpty(d.Get("key_algorithm")),
		"ca_hash_algorithm":        utils.ValueIgnoreEmpty(d.Get("ca_hash_algorithm")),
		"agree_privacy_protection": true,
	}
	return bodyParams
}

func resourceCertificateApplyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v3/scm/certificates/{certificate_id}/apply"
		product       = "scm"
		certificateID = d.Get("certificate_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{certificate_id}", certificateID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCertificateApplyBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error applying for CCM SSL certificate: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	result := utils.PathSearch("request_info", createRespBody, "").(string)
	if result != "success" {
		return diag.Errorf("error applying for CCM SSL certificate: the `request_info` in API response is not `success`")
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceCertificateApplyRead(ctx, d, meta)
}

func resourceCertificateApplyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCertificateApplyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
