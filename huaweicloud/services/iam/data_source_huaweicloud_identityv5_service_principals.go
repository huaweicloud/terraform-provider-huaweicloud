package iam

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

// @API IAM GET /v5/service-principals
func DataSourceIdentityV5ServicePrincipals() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5ServicePrincipalsRead,

		Schema: map[string]*schema.Schema{
			"language": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "zh-cn",
			},
			"service_principals": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_principal": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_catalog": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5ServicePrincipalsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allServicePrincipals []interface{}
	var marker string
	var path string
	for {
		path = fmt.Sprintf("%sv5/service-principals", client.Endpoint) + buildListServicePrincipalsV5Params(marker)
		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"X-Language": d.Get("language").(string),
			},
		}
		resp, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving service principals: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		servicePrincipals := flattenServicePrincipalsV5Response(respBody)
		allServicePrincipals = append(allServicePrincipals, servicePrincipals...)

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("service_principals", allServicePrincipals),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListServicePrincipalsV5Params(marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func flattenServicePrincipalsV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	principals := utils.PathSearch("service_principals", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(principals))
	for i, principal := range principals {
		result[i] = map[string]interface{}{
			"service_principal": utils.PathSearch("service_principal", principal, nil),
			"description":       utils.PathSearch("description", principal, nil),
			"display_name":      utils.PathSearch("display_name", principal, nil),
			"service_catalog":   utils.PathSearch("service_catalog", principal, nil),
		}
	}
	return result
}
