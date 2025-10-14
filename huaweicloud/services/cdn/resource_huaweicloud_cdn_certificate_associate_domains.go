package cdn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var cdnCertificateAssociateDomainsNonUpdatableParams = []string{"domain_names", "https_switch", "access_origin_way",
	"force_redirect_https", "force_redirect_config", "http2", "cert_name", "certificate", "private_key", "certificate_type"}

// @API CDN PUT /v1.0/cdn/domains/config-https-info
func ResourceCertificateAssociateDomains() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateAssociateDomainsCreate,
		ReadContext:   resourceCertificateAssociateDomainsRead,
		UpdateContext: resourceCertificateAssociateDomainsUpdate,
		DeleteContext: resourceCertificateAssociateDomainsDelete,

		CustomizeDiff: config.FlexibleForceNew(cdnCertificateAssociateDomainsNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"domain_names": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The list of domain names to associate with the certificate.`,
			},
			"https_switch": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The HTTPS certificate configuration switch.`,
			},

			// Optional parameters.
			"access_origin_way": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: `The origin protocol configuration.`,
			},
			"force_redirect_https": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: `Whether to enable HTTPS force redirect.`,
			},
			"force_redirect_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        forceRedirectConfigSchema(),
				Description: `The force redirect configuration.`,
			},
			"http2": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: `The HTTP/2 protocol switch.`,
			},
			"cert_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The certificate name.`,
			},
			"certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The SSL certificate content in PEM format.`,
			},
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The SSL certificate private key content in PEM format.`,
			},
			"certificate_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: `The certificate type.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func forceRedirectConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"switch": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The force redirect switch.`,
			},
			"redirect_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The redirect type.`,
			},
		},
	}
}

func buildForceRedirectConfig(redirectConfigs []interface{}) map[string]interface{} {
	if len(redirectConfigs) < 1 {
		return nil
	}

	return map[string]interface{}{
		"switch":        utils.PathSearch("switch", redirectConfigs[0], nil),
		"redirect_type": utils.PathSearch("redirect_type", redirectConfigs[0], nil),
	}
}

func buildCertificateAssociateDomainsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain_name":           d.Get("domain_names"),
		"https_switch":          d.Get("https_switch"),
		"access_origin_way":     utils.ValueIgnoreEmpty(d.Get("access_origin_way")),
		"force_redirect_https":  utils.ValueIgnoreEmpty(d.Get("force_redirect_https")),
		"force_redirect_config": buildForceRedirectConfig(d.Get("force_redirect_config").([]interface{})),
		"http2":                 utils.ValueIgnoreEmpty(d.Get("http2")),
		"cert_name":             utils.ValueIgnoreEmpty(d.Get("cert_name")),
		"certificate":           utils.ValueIgnoreEmpty(d.Get("certificate")),
		"private_key":           utils.ValueIgnoreEmpty(d.Get("private_key")),
		"certificate_type":      utils.ValueIgnoreEmpty(d.Get("certificate_type")),
	}

	return map[string]interface{}{
		"https": bodyParams,
	}
}

func parseCertificateAssociateDomainsError(respBody interface{}) error {
	status := utils.PathSearch("status", respBody, "").(string)
	if status == "success" {
		return nil
	}

	results := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
	var errorDetails []string
	for _, result := range results {
		domain := utils.PathSearch("domain_name", result, "").(string)
		reason := utils.PathSearch("reason", result, "").(string)
		if reason != "" {
			errorDetails = append(errorDetails, fmt.Sprintf("domain %s: %s", domain, reason))
		}
	}

	if len(errorDetails) > 0 {
		return fmt.Errorf("error associating certificate with domains: %s. Details: %s",
			status, strings.Join(errorDetails, "; "))
	}
	return fmt.Errorf("error associating certificate with domains: %s", status)
}

func resourceCertificateAssociateDomainsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1.0/cdn/domains/config-https-info"
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCertificateAssociateDomainsBodyParams(d),
	}

	resp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error associating certificate with domains: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing response: %s", err)
	}

	if err := parseCertificateAssociateDomainsError(respBody); err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceCertificateAssociateDomainsRead(ctx, d, meta)
}

func resourceCertificateAssociateDomainsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCertificateAssociateDomainsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCertificateAssociateDomainsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for bind certificate on the list of domains. Deleting this 
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
