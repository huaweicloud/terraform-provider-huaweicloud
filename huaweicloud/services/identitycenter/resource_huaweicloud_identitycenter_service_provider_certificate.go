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

var identityCenterServiceProviderCertificateNonUpdateParams = []string{"identity_store_id"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/saml-certificates
// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}/active
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/saml-certificates
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}
func ResourceIdentityCenterServiceProviderCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterServiceProviderCertificateCreate,
		UpdateContext: resourceIdentityCenterServiceProviderCertificateUpdate,
		ReadContext:   resourceIdentityCenterServiceProviderCertificateRead,
		DeleteContext: resourceIdentityCenterServiceProviderCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterServiceProviderCertificateImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterServiceProviderCertificateNonUpdateParams),

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
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"x509certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterServiceProviderCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates"
		activeHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}/active"
		product       = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center service provider certificate: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateId := utils.PathSearch("certificate_id", createRespBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("unable to find the Identity Center service provider certificate ID from the API response")
	}
	d.SetId(certificateId)

	if status, ok := d.GetOk("state"); ok && status == "ACTIVE" {
		activePath := client.Endpoint + activeHttpUrl
		activePath = strings.ReplaceAll(activePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
		activePath = strings.ReplaceAll(activePath, "{certificate_id}", d.Id())

		activeOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("POST", activePath, &activeOpt)
		if err != nil {
			return diag.Errorf("error activing Identity Center identity provider certificate: %s", err)
		}
	}

	return resourceIdentityCenterServiceProviderCertificateRead(ctx, d, meta)
}

func resourceIdentityCenterServiceProviderCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates"
		getProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(getProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center service provider certificate.")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificate := utils.PathSearch(fmt.Sprintf("[?certificate_id =='%s']|[0]", d.Id()), getRespBody, nil)

	if certificate == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no sp certificate found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("state", utils.PathSearch("state", certificate, nil)),
		d.Set("x509certificate", utils.PathSearch("x509certificate", certificate, nil)),
		d.Set("expiry_date", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("expiry_date", certificate, float64(0)).(float64))/1000, false)),
		d.Set("algorithm", utils.PathSearch("algorithm", certificate, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterServiceProviderCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		activeHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}/active"
		activeProduct = "identitystore"
	)

	if d.HasChange("state") {
		if d.Get("state").(string) != "ACTIVE" {
			return diag.Errorf("state must be set to 'ACTIVE'")
		}

		client, err := cfg.NewServiceClient(activeProduct, region)
		if err != nil {
			return diag.Errorf("error creating Identity Center Client: %s", err)
		}

		activePath := client.Endpoint + activeHttpUrl
		activePath = strings.ReplaceAll(activePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
		activePath = strings.ReplaceAll(activePath, "{certificate_id}", d.Id())

		activeOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("POST", activePath, &activeOpt)
		if err != nil {
			return diag.Errorf("error activing Identity Center identity provider certificate: %s", err)
		}
	}
	return resourceIdentityCenterServiceProviderCertificateRead(ctx, d, meta)
}

func resourceIdentityCenterServiceProviderCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}"
		deleteProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(deleteProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	deletePath = strings.ReplaceAll(deletePath, "{certificate_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center service provider certificate: %s", err)
	}

	return nil
}

func resourceIdentityCenterServiceProviderCertificateImport(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid id format, must be <identity_store_id>/<certificate_id>")
	}
	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("identity_store_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
