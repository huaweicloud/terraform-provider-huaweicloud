package eg

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

// @API EG GET /v1/{project_id}/quotas
func DataSourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEgQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The quota of the resource type to be queried.`,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of resource quotas that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The quota name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The quota type.`,
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The quota of current user.`,
						},
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The quota used by the current user.`,
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum quota.`,
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The minimum quota.`,
						},
					},
				},
			},
		},
	}
}

func buildQuotasQueryParams(d *schema.ResourceData) string {
	res := ""
	if quotaType, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s?type=%v", res, quotaType)
	}

	return res
}

func queryQuotas(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/quotas"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildQuotasQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("quotas.resources", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenQuotas(quotas []interface{}) []interface{} {
	if len(quotas) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(quotas))
	for _, resource := range quotas {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", resource, nil),
			"type":  utils.PathSearch("type", resource, nil),
			"quota": utils.PathSearch("quota", resource, nil),
			"used":  utils.PathSearch("used", resource, nil),
			"max":   utils.PathSearch("max", resource, nil),
			"min":   utils.PathSearch("min", resource, nil),
		})
	}

	return result
}

func dataSourceEgQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	quotas, err := queryQuotas(client, d)
	if err != nil {
		return diag.Errorf("error getting EG quotas: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	// Create the quotas structure with nested resources
	flattened := flattenQuotas(quotas)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattened),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
