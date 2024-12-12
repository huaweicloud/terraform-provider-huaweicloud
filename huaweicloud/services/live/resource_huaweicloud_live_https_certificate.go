package live

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LIVE PUT /v1/{project_id}/guard/https-cert
// @API LIVE GET /v1/{project_id}/guard/https-cert
// @API LIVE DELETE /v1/{project_id}/guard/https-cert

// ResourceHTTPSCertificate Due to lack of testing conditions, this resource has only been tested in limited scenarios.
// It is not yet certain whether this resource supports editing, so all fields are currently `forceNew`.
// There is currently no open documentation for this resource.
// Due to lack of test conditions, there is no support for `404` validation after successful certificate deletion.
func ResourceHTTPSCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHTTPSCertificateCreate,
		ReadContext:   resourceHTTPSCertificateRead,
		DeleteContext: resourceHTTPSCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceHTTPSCertificateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the streaming domain name`,
			},
			"certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the certificate body.`,
			},
			"certificate_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `Specifies the private key`,
			},
			"certificate_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the certificate format.`,
			},
			"force_redirect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether to enable redirection.`,
			},
			"gm_certificate": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        gmCertificateSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the GM certificate configuration.`,
			},
			"tls_certificate": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        tlsCertificateSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the TLS certificate configuration.`,
			},
		},
	}
}

func gmCertificateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"cert_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the SCM certificate name.`,
			},
			"cert_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the SCM certificate ID.`,
			},
			"sign_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the Chinese (SM) signature certificate body`,
			},
			"sign_certificate_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `Specifies the Chinese (SM) signature private key`,
			},
			"enc_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the Chinese (SM) encryption certificate body`,
			},
			"enc_certificate_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `Specifies the Chinese (SM) encryption private key`,
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the certificate source.`,
			},
		},
	}
	return &sc
}

func tlsCertificateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"cert_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the SCM certificate name.`,
			},
			"cert_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the SCM certificate ID.`,
			},
			"certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the certificate body.`,
			},
			"certificate_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `Specifies the private key.`,
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the certificate source.`,
			},
		},
	}
	return &sc
}

func buildGMCertificateBodyParams(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"source":               utils.ValueIgnoreEmpty(raw["source"]),
		"cert_name":            utils.ValueIgnoreEmpty(raw["cert_name"]),
		"cert_id":              utils.ValueIgnoreEmpty(raw["cert_id"]),
		"sign_certificate":     utils.ValueIgnoreEmpty(raw["sign_certificate"]),
		"sign_certificate_key": utils.ValueIgnoreEmpty(raw["sign_certificate_key"]),
		"enc_certificate":      utils.ValueIgnoreEmpty(raw["enc_certificate"]),
		"enc_certificate_key":  utils.ValueIgnoreEmpty(raw["enc_certificate_key"]),
	}
}

func buildTLSCertificateBodyParams(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"source":          utils.ValueIgnoreEmpty(raw["source"]),
		"cert_name":       utils.ValueIgnoreEmpty(raw["cert_name"]),
		"cert_id":         utils.ValueIgnoreEmpty(raw["cert_id"]),
		"certificate":     utils.ValueIgnoreEmpty(raw["certificate"]),
		"certificate_key": utils.ValueIgnoreEmpty(raw["certificate_key"]),
	}
}

func buildUpdateHTTPSCertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"certificate_format": utils.ValueIgnoreEmpty(d.Get("certificate_format")),
		"certificate":        utils.ValueIgnoreEmpty(d.Get("certificate")),
		"certificate_key":    utils.ValueIgnoreEmpty(d.Get("certificate_key")),
		"force_redirect":     utils.ValueIgnoreEmpty(d.Get("force_redirect")),
		"gm_certificate":     buildGMCertificateBodyParams(d.Get("gm_certificate")),
		"tls_certificate":    buildTLSCertificateBodyParams(d.Get("tls_certificate")),
	}
}

func resourceHTTPSCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/guard/https-cert"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	resourceID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating Live HTTPS certificate resource ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildHTTPSCertificateQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateHTTPSCertificateBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating Live HTTPS certificate: %s", err)
	}

	d.SetId(resourceID)

	return resourceHTTPSCertificateRead(ctx, d, meta)
}

func buildHTTPSCertificateQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?domain=%s", d.Get("domain_name").(string))
}

func resourceHTTPSCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/guard/https-cert"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildHTTPSCertificateQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "LIVE.100011001"),
			"error retrieving Live HTTPS certificate")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("certificate_format", utils.PathSearch("certificate_format", respBody, nil)),
		d.Set("certificate", utils.PathSearch("certificate", respBody, nil)),
		d.Set("certificate_key", utils.PathSearch("certificate_key", respBody, nil)),
		d.Set("force_redirect", utils.PathSearch("force_redirect", respBody, nil)),
		d.Set("gm_certificate", flattenGMCertificateResponseBody(utils.PathSearch("gm_certificate", respBody, nil))),
		d.Set("tls_certificate", flattenTLSCertificateResponseBody(utils.PathSearch("tls_certificate", respBody, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGMCertificateResponseBody(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"source":               utils.PathSearch("source", respBody, nil),
		"cert_name":            utils.PathSearch("cert_name", respBody, nil),
		"cert_id":              utils.PathSearch("cert_id", respBody, nil),
		"sign_certificate":     utils.PathSearch("sign_certificate", respBody, nil),
		"sign_certificate_key": utils.PathSearch("sign_certificate_key", respBody, nil),
		"enc_certificate":      utils.PathSearch("enc_certificate", respBody, nil),
		"enc_certificate_key":  utils.PathSearch("enc_certificate_key", respBody, nil),
	}}
}

func flattenTLSCertificateResponseBody(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"source":          utils.PathSearch("source", respBody, nil),
		"cert_name":       utils.PathSearch("cert_name", respBody, nil),
		"cert_id":         utils.PathSearch("cert_id", respBody, nil),
		"certificate":     utils.PathSearch("certificate", respBody, nil),
		"certificate_key": utils.PathSearch("certificate_key", respBody, nil),
	}}
}

func resourceHTTPSCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/guard/https-cert"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildHTTPSCertificateQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "LIVE.100011001"),
			"error deleting Live HTTPS certificate")
	}

	return nil
}

func resourceHTTPSCertificateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)
	return []*schema.ResourceData{d}, d.Set("domain_name", importedId)
}
