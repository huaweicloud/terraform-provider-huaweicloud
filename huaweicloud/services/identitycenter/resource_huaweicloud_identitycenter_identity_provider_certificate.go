package identitycenter

import (
	"context"
	"errors"
	"fmt"
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

var identityCenterIdentityProviderCertificateNonUpdateParams = []string{"identity_store_id", "idp_id", "x509_certificate_in_pem", "certificate_use"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate/{certificate_id}
func ResourceIdentityCenterIdentityProviderCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterIdentityProviderCertificateCreate,
		UpdateContext: resourceIdentityCenterIdentityProviderCertificateUpdate,
		ReadContext:   resourceIdentityCenterIdentityProviderCertificateRead,
		DeleteContext: resourceIdentityCenterIdentityProviderCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterIdentityProviderCertificateImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterIdentityProviderCertificateNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"idp_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"x509_certificate_in_pem": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate_use": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"issuer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"not_after": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"not_before": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial_number_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_algorithm_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterIdentityProviderCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		importHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate"
		importProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(importProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	importPath := client.Endpoint + importHttpUrl
	importPath = strings.ReplaceAll(importPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	importPath = strings.ReplaceAll(importPath, "{idp_id}", fmt.Sprintf("%v", d.Get("idp_id")))

	importOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildImportExternalIdPCertificateBodyParams(d)),
	}

	importResp, err := client.Request("POST", importPath, &importOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center identity provider certificate: %s", err)
	}

	importRespBody, err := utils.FlattenResponse(importResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateId := utils.PathSearch("certificate_id", importRespBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("unable to find the Identity Center identity provider certificate ID from the API response")
	}
	d.SetId(certificateId)

	return resourceIdentityCenterIdentityProviderCertificateRead(ctx, d, meta)
}

func buildImportExternalIdPCertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"x509_certificate_in_pem": d.Get("x509_certificate_in_pem"),
		"certificate_use":         d.Get("certificate_use"),
	}
	return bodyParams
}

func resourceIdentityCenterIdentityProviderCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate"
		listProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	listPath = strings.ReplaceAll(listPath, "{idp_id}", fmt.Sprintf("%v", d.Get("idp_id")))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center identity provider certificate.")
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificate := utils.PathSearch(fmt.Sprintf("idp_certificates[?certificate_id=='%s']|[0]", d.Id()), listRespBody, nil)
	if certificate == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no idp certificates found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("issuer_name", utils.PathSearch("issuer_name", certificate, nil)),
		d.Set("not_after", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("not_after", certificate, float64(0)).(float64))/1000, false)),
		d.Set("not_before", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("not_before", certificate, float64(0)).(float64))/1000, false)),
		d.Set("public_key", utils.PathSearch("public_key", certificate, nil)),
		d.Set("serial_number_string", utils.PathSearch("serial_number_string", certificate, nil)),
		d.Set("subject_name", utils.PathSearch("subject_name", certificate, nil)),
		d.Set("version", utils.PathSearch("version", certificate, nil)),
		d.Set("signature_algorithm_name", utils.PathSearch("signature_algorithm_name", certificate, nil)),
		d.Set("x509_certificate_in_pem", utils.PathSearch("x509_Certificate_in_pem", certificate, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterIdentityProviderCertificateUpdate(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterIdentityProviderCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate/{certificate_id}"
		deleteProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(deleteProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	deletePath = strings.ReplaceAll(deletePath, "{idp_id}", fmt.Sprintf("%v", d.Get("idp_id")))
	deletePath = strings.ReplaceAll(deletePath, "{certificate_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center identity provider certificate: %s", err)
	}

	return nil
}

func resourceIdentityCenterIdentityProviderCertificateImport(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, errors.New("invalid id format, must be <identity_store_id>/<idp_id>/<certificate_id>")
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("identity_store_id", parts[0]),
		d.Set("idp_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
