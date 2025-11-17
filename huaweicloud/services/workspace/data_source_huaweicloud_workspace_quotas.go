package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/quotas
func DataSourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the quotas are located.`,
			},

			// Attributes.
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of quota resources.`,
							Elem:        quotaResourceSchema(),
						},
					},
				},
				Description: `The list of common quotas.`,
			},
			"site_quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The site ID.`,
						},
						"resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of quota resources.`,
							Elem:        quotaResourceSchema(),
						},
					},
				},
				Description: `The list of site quotas.`,
			},
		},
	}
}

func quotaResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource type.`,
			},
			"quota": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The quota value.`,
			},
			"used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The used quota value.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The quota unit.`,
			},
		},
	}
}

func getQuotas(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/quotas"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenQuotaResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"type":  utils.PathSearch("type", resource, nil),
			"quota": utils.PathSearch("quota", resource, nil),
			"used":  utils.PathSearch("used", resource, nil),
			"unit":  utils.PathSearch("unit", resource, nil),
		})
	}

	return result
}

func flattenQuotas(quotas map[string]interface{}) []map[string]interface{} {
	if len(quotas) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"resources": flattenQuotaResources(utils.PathSearch("resources", quotas, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenSiteQuotas(siteQuotas []interface{}) []map[string]interface{} {
	if len(siteQuotas) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(siteQuotas))
	for _, siteQuota := range siteQuotas {
		result = append(result, map[string]interface{}{
			"site_id":   utils.PathSearch("site_id", siteQuota, nil),
			"resources": flattenQuotaResources(utils.PathSearch("resources", siteQuota, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	resp, err := getQuotas(client)
	if err != nil {
		return diag.Errorf("error querying Workspace quotas: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattenQuotas(utils.PathSearch("quotas", resp, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("site_quotas", flattenSiteQuotas(utils.PathSearch("site_quotas", resp, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
