package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/users/role
func DataSourceStudioWorkspaceUserRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStudioWorkspaceUserRolesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the workspace user roles are located.`,
			},

			// Optional parameters.
			"instance_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"instance_id", "workspace_id"},
				Description:  `The instance ID to query user roles.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID to query user roles.`,
			},

			// Attributes.
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The role ID.`,
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The role code.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The role name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The role description.`,
						},
					},
				},
				Description: `The list of user roles that matched filter parameters.`,
			},
		},
	}
}

func buildStudioWorkspaceUserRolesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("workspace_id"); ok {
		res = fmt.Sprintf("%s&workspace_id=%v", res, v)
	}

	if len(res) > 0 {
		return fmt.Sprintf("?%s", res[1:])
	}
	return res
}

func getStudioWorkspaceUserRoles(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/users/role"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += buildStudioWorkspaceUserRolesQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", path, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("[]", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenStudioWorkspaceUserRoles(roles []interface{}) []map[string]interface{} {
	if len(roles) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(roles))
	for _, role := range roles {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("role_id", role, nil),
			"code":        utils.PathSearch("role_code", role, nil),
			"name":        utils.PathSearch("role_name", role, nil),
			"description": utils.PathSearch("description", role, nil),
		})
	}
	return result
}

func dataSourceStudioWorkspaceUserRolesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	roles, err := getStudioWorkspaceUserRoles(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Studio workspace user roles: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("roles", flattenStudioWorkspaceUserRoles(roles)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
