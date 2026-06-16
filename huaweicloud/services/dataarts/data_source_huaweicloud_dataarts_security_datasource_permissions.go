package dataarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/permission-sets/datasource/configurations
func DataSourceSecurityDatasourcePermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityDatasourcePermissionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the datasource permissions are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the datasource permissions belong.`,
			},

			// Attributes.
			"datasource_permissions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of datasource configurable permissions.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datasource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The datasource type.`,
						},
						"permission_types": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The list of datasource operation permission types.`,
						},
						"permission_actions": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The list of supported datasource permission actions.`,
						},
					},
				},
			},
		},
	}
}

func listSecurityDatasourceConfigurations(client *golangsdk.ServiceClient, workspaceId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/permission-sets/datasource/configurations"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("configurations", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenSecurityDatasourcePermissions(configurations []interface{}) []map[string]interface{} {
	if len(configurations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(configurations))
	for _, configuration := range configurations {
		result = append(result, map[string]interface{}{
			"datasource_type": utils.PathSearch("datasource_type", configuration, nil),
			"permission_types": utils.PathSearch("permission_types", configuration,
				make([]interface{}, 0)),
			"permission_actions": utils.PathSearch("permission_actions", configuration,
				make([]interface{}, 0)),
		})
	}

	return result
}

func dataSourceSecurityDatasourcePermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	configurations, err := listSecurityDatasourceConfigurations(client, workspaceId)
	if err != nil {
		return diag.Errorf("error querying DataArts Security datasource permissions: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("datasource_permissions", flattenSecurityDatasourcePermissions(configurations)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
