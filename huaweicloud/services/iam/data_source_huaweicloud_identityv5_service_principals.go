package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/service-principals
func DataSourceV5ServicePrincipals() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5ServicePrincipalsRead,

		Schema: map[string]*schema.Schema{
			"language": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "zh-cn",
				Description: `The language of the information returned by the interface.`,
			},
			"service_principals": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_principal": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the service principal.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the service principal.`,
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The display name of the service principal.`,
						},
						"service_catalog": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cloud service name.`,
						},
					},
				},
				Description: `The list of service principals.`,
			},
		},
	}
}

func buildV5ServicePrincipalsQueryParams(marker string) string {
	if marker != "" {
		return fmt.Sprintf("&marker=%v", marker)
	}
	return ""
}

func listV5ServicePrincipals(client *golangsdk.ServiceClient, language string) ([]interface{}, error) {
	var (
		httpUrl = "v5/service-principals"
		result  = make([]interface{}, 0)
		marker  = ""
		// Default limit is 100, maximum is 200.
		limit = 200
	)

	listPath := client.Endpoint + fmt.Sprintf("%s?limit=%d", httpUrl, limit)
	for {
		listPathWithMarker := listPath + buildV5ServicePrincipalsQueryParams(marker)
		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"X-Language": language,
			},
		}
		resp, err := client.Request("GET", listPathWithMarker, reqOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		principals := utils.PathSearch("service_principals", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, principals...)
		if len(principals) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}
	return result, nil
}

func dataSourceV5ServicePrincipalsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	principals, err := listV5ServicePrincipals(client, d.Get("language").(string))
	if err != nil {
		return diag.Errorf("error retrieving service principals: %s", err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomId)

	return diag.FromErr(d.Set("service_principals", flattenV5ServicePrincipals(principals)))
}

func flattenV5ServicePrincipals(principals []interface{}) []interface{} {
	if len(principals) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(principals))
	for _, principal := range principals {
		result = append(result, map[string]interface{}{
			"service_principal": utils.PathSearch("service_principal", principal, nil),
			"description":       utils.PathSearch("description", principal, nil),
			"display_name":      utils.PathSearch("display_name", principal, nil),
			"service_catalog":   utils.PathSearch("service_catalog", principal, nil),
		})
	}
	return result
}
