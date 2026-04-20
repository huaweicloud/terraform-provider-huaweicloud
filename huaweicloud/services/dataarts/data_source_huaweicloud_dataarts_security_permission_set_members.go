package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/permission-sets/{permission_set_id}/members
func DataSourceSecurityPermissionSetMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityPermissionSetMembersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the permission set associated members are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the permission set belongs.`,
			},
			"permission_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the permission set which the members associated.`,
			},

			// Optional parameters.
			"member_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the specified permission set member to be queried.`,
			},
			"member_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the permission set members to be queried.`,
			},

			// Attributes.
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the permission set member.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the permission set member.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the permission set member.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the permission set member.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID to which the permission set member belongs.`,
						},
						"workspace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The workspace ID to which the permission set member belongs.`,
						},
						"permission_set_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the permission set to which the permission set member belongs.`,
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster ID to which the permission set member belongs.`,
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster type to which the permission set member belongs.`,
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster name to which the permission set member belongs.`,
						},
						"create_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the permission set member.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the permission set member, in RFC3339 format.`,
						},
						"deadline": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The deadline time of the permission set member, in RFC3339 format.`,
						},
					},
				},
				Description: `The list of permission set associated members that match the filter parameters.`,
			},
		},
	}
}

func buildSecurityPermissionSetMembersQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("member_name"); ok {
		res = fmt.Sprintf("%s&member_name=%v", res, v)
	}
	if v, ok := d.GetOk("member_type"); ok {
		res = fmt.Sprintf("%s&member_type=%v", res, v)
	}

	return res
}

func listSecurityPermissionSetMembers(client *golangsdk.ServiceClient, workspaceId, permissionSetId string,
	queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/permission-sets/{permission_set_id}/members?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{permission_set_id}", permissionSetId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		members := utils.PathSearch("permission_set_members", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, members...)
		if len(members) < limit {
			break
		}
		offset += len(members)
	}

	return result, nil
}

func flattenSecurityPermissionSetMembers(members []interface{}) []map[string]interface{} {
	if len(members) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(members))
	for _, member := range members {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("member_id", member, nil),
			"type":              utils.PathSearch("member_type", member, nil),
			"name":              utils.PathSearch("member_name", member, nil),
			"status":            utils.PathSearch("member_status", member, nil),
			"instance_id":       utils.PathSearch("instance_id", member, nil),
			"workspace":         utils.PathSearch("workspace", member, nil),
			"permission_set_id": utils.PathSearch("permission_set_id", member, nil),
			"cluster_id":        utils.PathSearch("cluster_id", member, nil),
			"cluster_type":      utils.PathSearch("cluster_type", member, nil),
			"cluster_name":      utils.PathSearch("cluster_name", member, nil),
			"create_user":       utils.PathSearch("create_user", member, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				member, float64(0)).(float64))/1000, false),
			"deadline": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("deadline",
				member, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceSecurityPermissionSetMembersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		workspaceId     = d.Get("workspace_id").(string)
		permissionSetId = d.Get("permission_set_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	members, err := listSecurityPermissionSetMembers(client, workspaceId, permissionSetId, buildSecurityPermissionSetMembersQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying permission set members: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("members", flattenSecurityPermissionSetMembers(members)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
