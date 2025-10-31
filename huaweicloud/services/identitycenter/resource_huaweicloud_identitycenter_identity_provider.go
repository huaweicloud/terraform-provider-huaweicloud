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

var identityCenterIdentityProviderNonUpdateParams = []string{"identity_store_id", "idp_saml_metadata", "idp_certificate"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/external-idp
// @API IdentityCenter PUT /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/external-idp
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}
// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/enable
// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/disable
func ResourceIdentityCenterIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterIdentityProviderCreate,
		UpdateContext: resourceIdentityCenterIdentityProviderUpdate,
		ReadContext:   resourceIdentityCenterIdentityProviderRead,
		DeleteContext: resourceIdentityCenterIdentityProviderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterIdentityProviderImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterIdentityProviderNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"idp_saml_metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idp_certificate": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"entity_id", "login_url"},
			},
			"entity_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"login_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"want_request_signed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"idp_certificate_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIdentityCenterIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createHttpUrl        = "v1/identity-stores/{identity_store_id}/external-idp"
		enableHttpUrl        = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/enable"
		identityStoreProduct = "identitystore"
	)

	client, err := cfg.NewServiceClient(identityStoreProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateExternalIdPConfigurationForDirectoryBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating identity provider: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	idpId := utils.PathSearch("idp_id", createRespBody, "").(string)
	if idpId == "" {
		return diag.Errorf("unable to find the Identity Center identity provider ID from the API response")
	}

	d.SetId(idpId)

	enablePath := client.Endpoint + enableHttpUrl
	enablePath = strings.ReplaceAll(enablePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	enablePath = strings.ReplaceAll(enablePath, "{idp_id}", d.Id())

	enableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", enablePath, &enableOpt)
	if err != nil {
		return diag.Errorf("error enabling identity provider: %s", err)
	}

	return resourceIdentityCenterIdentityProviderRead(ctx, d, meta)
}

func buildCreateExternalIdPConfigurationForDirectoryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"idp_saml_metadata": utils.ValueIgnoreEmpty(d.Get("idp_saml_metadata")),
		"idp_certificate":   utils.ValueIgnoreEmpty(d.Get("idp_certificate")),
		"idp_saml_config": map[string]interface{}{
			"entity_id": utils.ValueIgnoreEmpty(d.Get("entity_id")),
			"login_url": utils.ValueIgnoreEmpty(d.Get("login_url")),
		},
	}
	return bodyParams
}

func resourceIdentityCenterIdentityProviderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp"
		listProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving identity provider.")
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	association := utils.PathSearch(fmt.Sprintf("associations[?idp_id=='%s']|[0]", d.Id()), listRespBody, nil)
	if association == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no identity provider found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("entity_id", utils.PathSearch("idp_saml_config.entity_id", association, nil)),
		d.Set("login_url", utils.PathSearch("idp_saml_config.login_url", association, nil)),
		d.Set("is_enabled", utils.PathSearch("is_enabled", association, nil)),
		d.Set("want_request_signed", utils.PathSearch("idp_saml_config.want_request_signed", association, nil)),
		d.Set("idp_certificate_ids", flattenIdentityCenterIdentityProviderCertificate(association)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenIdentityCenterIdentityProviderCertificate(data interface{}) []interface{} {
	if data == nil {
		return nil
	}
	curJson := utils.PathSearch("idp_certificate_ids", data, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	certs := make([]interface{}, len(curArray))
	for i, v := range curArray {
		certs[i] = map[string]interface{}{
			"certificate_id": utils.PathSearch("certificate_id", v, ""),
			"status":         utils.PathSearch("status", v, ""),
		}
	}
	return certs
}

func resourceIdentityCenterIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}"
		updateProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(updateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	updateChanges := []string{
		"entity_id",
		"login_url",
	}

	if d.HasChanges(updateChanges...) {
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
		updatePath = strings.ReplaceAll(updatePath, "{idp_id}", d.Id())

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateExternalIdPConfigurationsForDirectoryBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating identity provider: %s", err)
		}
	}

	return resourceIdentityCenterIdentityProviderRead(ctx, d, meta)
}

func buildUpdateExternalIdPConfigurationsForDirectoryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"entity_id": utils.ValueIgnoreEmpty(d.Get("entity_id")),
		"login_url": utils.ValueIgnoreEmpty(d.Get("login_url")),
	}
	return bodyParams
}

func resourceIdentityCenterIdentityProviderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHttpUrl        = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}"
		disableHttpUrl       = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/disable"
		identityStoreProduct = "identitystore"
	)

	client, err := cfg.NewServiceClient(identityStoreProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	disablePath := client.Endpoint + disableHttpUrl
	disablePath = strings.ReplaceAll(disablePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	disablePath = strings.ReplaceAll(disablePath, "{idp_id}", d.Id())

	disableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("POST", disablePath, &disableOpt)
	if err != nil {
		return diag.Errorf("error deleting identity provider: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	deletePath = strings.ReplaceAll(deletePath, "{idp_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting identity provider: %s", err)
	}

	return nil
}

func resourceIdentityCenterIdentityProviderImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid id format, must be <identity_store_id>/<idp_id>")
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
