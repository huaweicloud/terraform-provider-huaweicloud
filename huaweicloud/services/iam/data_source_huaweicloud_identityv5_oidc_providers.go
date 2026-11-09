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

// @API IAM GET /v5/oidc-providers
func DataSourceV5OidcProviders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5OidcProvidersRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the OIDC identity provider used for filtering.`,
			},
			"providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the OIDC identity provider.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the OIDC identity provider.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the OIDC identity provider.`,
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URL of the OIDC identity provider.`,
						},
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The uniform resource name of the OIDC identity provider.`,
						},
						"client_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of client IDs registered in the OIDC identity provider.`,
						},
						"thumbprints": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of server certificate fingerprints of the OIDC identity provider.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the OIDC identity provider.`,
						},
						"expires_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The expiration time of the OIDC identity provider.`,
						},
					},
				},
				Description: `The list of OIDC identity providers that matched the filter parameters.`,
			},
		},
	}
}

func dataSourceV5OidcProvidersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	providers, err := listV5OidcProviders(client)
	if err != nil {
		return diag.Errorf("error querying OIDC identity providers: %s", err)
	}

	if name, ok := d.GetOk("name"); ok {
		providers = utils.PathSearch(fmt.Sprintf("[?name=='%s']", name.(string)), providers, make([]interface{}, 0)).([]interface{})
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("providers", flattenV5OidcProviders(providers)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OIDC identity providers fields: %s", err)
	}

	return nil
}

func flattenV5OidcProviders(providers []interface{}) []interface{} {
	if len(providers) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(providers))
	for _, provider := range providers {
		result = append(result, map[string]interface{}{
			"provider_id": utils.PathSearch("provider_id", provider, nil),
			"name":        utils.PathSearch("name", provider, nil),
			"description": utils.PathSearch("description", provider, nil),
			"url":         utils.PathSearch("url", provider, nil),
			"urn":         utils.PathSearch("urn", provider, nil),
			"client_ids":  utils.PathSearch("client_ids", provider, make([]interface{}, 0)),
			"thumbprints": utils.PathSearch("thumbprints", provider, make([]interface{}, 0)),
			"created_at":  utils.PathSearch("created_at", provider, nil),
			"expires_at":  utils.PathSearch("expires_at", provider, nil),
		})
	}
	return result
}
