package eps

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EPS GET /v1.0/enterprise-projects/quotas
func DataSourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource quotas.`,
			},

			// Attributes.
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total number of the resource quota.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type corresponding to quota.`,
						},
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The used quota number.`,
						},
					},
				},
				Description: `The list of the resource quotas.`,
			},
		},
	}
}

func listResourceQuotas(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v1.0/enterprise-projects/quotas"
	listPath := client.Endpoint + httpUrl

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

func flattenResourceQuotas(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"quota": utils.PathSearch("quota", resource, nil),
			"type":  utils.PathSearch("type", resource, nil),
			"used":  utils.PathSearch("used", resource, nil),
		})
	}

	return result
}

func dataSourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	staticRoutes, err := listResourceQuotas(client)
	if err != nil {
		return diag.Errorf("error querying resource quotas: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattenResourceQuotas(staticRoutes)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the EPS resource quotas: %s", mErr)
	}
	return nil
}
