package identitycenter

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IDENTITYCENTER GET /v1/application-providers
func DataSourceIdentityCenterApplicationProviders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterApplicationProvidersRead,

		Schema: map[string]*schema.Schema{
			"application_providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_provider_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `the urn of application.`,
						},
						"display_data": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"icon_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"federation_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_provider_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListApplicationProvidersParams(marker string) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func listApplicationProviders(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/application-providers"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl

	queryParams := buildListApplicationProvidersParams("")

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestResp, err := client.Request("GET", listPath+queryParams, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Identity Center application providers: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		providers := utils.PathSearch("application_providers", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, providers...)

		marker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}

		queryParams = buildListApplicationProvidersParams(marker)
	}
	return result, nil
}

func dataSourceIdentityCenterApplicationProvidersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	providers, err := listApplicationProviders(client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("application_providers", flattenApplicationProviders(providers)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenApplicationProviders(providers []interface{}) []interface{} {
	if len(providers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(providers))
	for _, provider := range providers {
		result = append(result, map[string]interface{}{
			"application_provider_urn": utils.PathSearch("application_provider_urn", provider, nil),
			"application_provider_id":  utils.PathSearch("application_provider_id", provider, nil),
			"federation_protocol":      utils.PathSearch("federation_protocol", provider, nil),
			"display_data":             flattenDisplayData(utils.PathSearch("display_data", provider, nil)),
		})
	}
	return result
}

func flattenDisplayData(display interface{}) []map[string]interface{} {
	if display == nil || len(display.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"display_name": utils.PathSearch("display_name", display, nil),
			"description":  utils.PathSearch("description", display, nil),
			"icon_url":     utils.PathSearch("icon_url", display, nil),
		},
	}
}
