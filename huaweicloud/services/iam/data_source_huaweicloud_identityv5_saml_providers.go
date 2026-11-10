package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/saml-providers
// @API IAM GET /v5/saml-providers/{provider_id}
func DataSourceV5SamlProviders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5SamlProvidersRead,

		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the SAML identity provider used for filtering. If this parameter is set, only the provider with the matched ID is returned.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the SAML identity provider used for filtering.`,
			},
			"providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the SAML identity provider.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the SAML identity provider.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the SAML identity provider.`,
						},
						"saml_metadata_document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SAML metadata document of the SAML identity provider.`,
						},
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The uniform resource name of the SAML identity provider.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the SAML identity provider.`,
						},
						"expires_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The expiration time of the SAML identity provider.`,
						},
						"private_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the private key.`,
									},
									"timestamp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The upload time of the private key in ISO 8601 format.`,
									},
								},
							},
							Description: `The list of private keys managed by the SAML identity provider.`,
						},
						"assertion_encryption_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The assertion encryption mode of the SAML identity provider.`,
						},
					},
				},
				Description: `The list of SAML identity providers that matched the filter parameters.`,
			},
		},
	}
}

func dataSourceV5SamlProvidersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var providers []interface{}

	if providerId, ok := d.GetOk("provider_id"); ok {
		provider, err := GetV5SamlProviderById(client, providerId.(string))
		if err != nil {
			return diag.Errorf("error querying SAML identity provider (ID: %s): %s", providerId, err)
		}
		providers = []interface{}{provider}
	} else {
		providers, err = listV5SamlProviders(client)
		if err != nil {
			return diag.Errorf("error querying SAML identity providers: %s", err)
		}

		if name, ok := d.GetOk("name"); ok {
			providers = utils.PathSearch(fmt.Sprintf("[?name=='%s']", name.(string)), providers, make([]interface{}, 0)).([]interface{})
		}
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("providers", flattenV5SamlProviders(providers)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SAML identity providers fields: %s", err)
	}

	return nil
}

func flattenV5SamlProviders(providers []interface{}) []interface{} {
	if len(providers) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(providers))
	for _, provider := range providers {
		result = append(result, map[string]interface{}{
			"provider_id":               utils.PathSearch("provider_id", provider, nil),
			"name":                      utils.PathSearch("name", provider, nil),
			"description":               utils.PathSearch("description", provider, nil),
			"saml_metadata_document":    utils.PathSearch("saml_metadata_document", provider, nil),
			"urn":                       utils.PathSearch("urn", provider, nil),
			"created_at":                utils.PathSearch("created_at", provider, nil),
			"expires_at":                utils.PathSearch("expires_at", provider, nil),
			"private_keys":              flattenV5SamlProviderPrivateKeys(utils.PathSearch("private_keys", provider, nil)),
			"assertion_encryption_mode": utils.PathSearch("assertion_encryption_mode", provider, nil),
		})
	}
	return result
}
