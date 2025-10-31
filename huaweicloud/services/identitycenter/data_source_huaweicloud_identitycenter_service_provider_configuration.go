package identitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/sp-config
func DataSourceIdentityCenterServiceProviderConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterServiceProviderConfigurationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sp_oidc_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redirect_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"sp_saml_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acs_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityCenterServiceProviderConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/sp-config"
		product = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET",
		path, &opt)

	if err != nil {
		return diag.FromErr(err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("sp_oidc_config", flattenSpOidcConfig(utils.PathSearch("sp_oidc_config",
			respBody, nil))),
		d.Set("sp_saml_config", flattenSpSamlConfig(utils.PathSearch("sp_saml_config",
			respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSpOidcConfig(spOidcConfig interface{}) []map[string]interface{} {
	if spOidcConfig == nil || len(spOidcConfig.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"redirect_url": utils.PathSearch("redirect_url", spOidcConfig, nil),
		},
	}
}

func flattenSpSamlConfig(spSamlConfig interface{}) []map[string]interface{} {
	if spSamlConfig == nil || len(spSamlConfig.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"issuer":   utils.PathSearch("issuer", spSamlConfig, nil),
			"acs_url":  utils.PathSearch("acs_url", spSamlConfig, nil),
			"metadata": utils.PathSearch("metadata", spSamlConfig, nil),
		},
	}
}
