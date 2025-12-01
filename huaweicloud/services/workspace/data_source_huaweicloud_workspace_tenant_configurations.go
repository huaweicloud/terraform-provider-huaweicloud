package workspace

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

// @API Workspace GET /v2/{project_id}/tenant-configs
func DataSourceTenantConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTenantConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the tenant configurations are located.`,
			},
			"configuration_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the configuration.`,
			},

			// Attributes
			"configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of tenant configurations that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the configuration.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the configuration.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the configuration.`,
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The configuration values.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The key of the configuration item.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of the configuration item.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildTenantConfigurationsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("configuration_name"); ok {
		res += fmt.Sprintf("&name=%v", v)
	}

	if len(res) > 1 {
		return "?" + res[1:]
	}
	return res
}

func dataSourceTenantConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	httpUrl := "v2/{project_id}/tenant-configs"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildTenantConfigurationsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving tenant configurations: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	configurations := utils.PathSearch("function_configs", respBody, make([]interface{}, 0)).([]interface{})

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("configurations", flattenTenantConfigurations(configurations)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTenantConfigurations(configurations []interface{}) []map[string]interface{} {
	if len(configurations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(configurations))
	for _, configuration := range configurations {
		result = append(result, map[string]interface{}{
			"id":     utils.PathSearch("id", configuration, nil),
			"name":   utils.PathSearch("name", configuration, nil),
			"status": utils.PathSearch("status", configuration, nil),
			"values": flattenConfigurationValues(utils.PathSearch("values", configuration,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenConfigurationValues(values []interface{}) []map[string]interface{} {
	if len(values) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(values))
	for _, value := range values {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", value, nil),
			"value": utils.PathSearch("value", value, nil),
		})
	}

	return result
}
