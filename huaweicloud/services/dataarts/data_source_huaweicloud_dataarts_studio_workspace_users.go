package dataarts

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v2/{project_id}/{workspace_id}/users
func DataSourceStudioWorkspaceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStudioWorkspaceUsersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the workspace users are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workspace ID to which the users belong.`,
			},

			// Attributes.
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the IAM user to which the workspace user correspond.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the IAM user to which the workspace user correspond.`,
						},
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
							Description: `The role list of the workspace user.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the workspace user, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the workspace user, in RFC3339 format.`,
						},
					},
				},
				Description: `The list of workspace users.`,
			},
		},
	}
}

func flattenStudioWorkspaceUsers(users []interface{}) []map[string]interface{} {
	if len(users) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"id":    utils.PathSearch("user_id", user, nil),
			"name":  utils.PathSearch("user_name", user, nil),
			"roles": flattenStudioWorkspaceUserRoles(utils.PathSearch("roles", user, make([]interface{}, 0)).([]interface{})),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", user,
				float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", user,
				float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceStudioWorkspaceUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	users, err := listStudioWorkspaceUsers(client, workspaceId)
	if err != nil {
		return diag.Errorf("error querying DataArts Studio workspace users: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("users", flattenStudioWorkspaceUsers(users)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
