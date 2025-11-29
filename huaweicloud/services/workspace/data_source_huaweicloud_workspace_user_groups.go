package workspace

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

// @API Workspace GET /v2/{project_id}/groups
func DataSourceUserGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the user groups are located.`,
			},

			// Attributes.
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        userGroupSchema(),
				Description: `The list of user groups.`,
			},
		},
	}
}

func userGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the user group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the user group.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the user group, in RFC3339 format.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the user group.`,
			},
			"user_quantity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of users in the user list.`,
			},
			"parent": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        userGroupParentSchema(),
				Description: `The parent user group of the user group.`,
			},
			"realm_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain ID of the user group.`,
			},
			"platform_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the user group.`,
			},
			"group_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The distinguished name of the user group.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name of the user group.`,
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SID of the user group.`,
			},
			"total_desktops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of users in the user list.`,
			},
		},
	}
}

func userGroupParentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the parent user group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the parent user group.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the parent user group, in RFC3339 format.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the parent user group.`,
			},
			"user_quantity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of users in the parent user group.`,
			},
			"realm_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain ID of the parent user group.`,
			},
			"platform_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the parent user group.`,
			},
			"group_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The distinguished name of the parent user group.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name of the parent user group.`,
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SID of the parent user group.`,
			},
			"total_desktops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of users in the parent user group.`,
			},
		},
	}
}

func listUserGroups(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/groups?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		userGroups := utils.PathSearch("user_groups", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, userGroups...)
		if len(userGroups) < limit {
			break
		}
		offset += len(userGroups)
	}

	return result, nil
}

func flattenUserGroupParent(item map[string]interface{}) []map[string]interface{} {
	if len(item) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":   utils.PathSearch("id", item, nil),
			"name": utils.PathSearch("name", item, nil),
			"create_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
				item, "").(string))/1000, false),
			"description":    utils.PathSearch("description", item, nil),
			"user_quantity":  utils.PathSearch("user_quantity", item, nil),
			"realm_id":       utils.PathSearch("realm_id", item, nil),
			"platform_type":  utils.PathSearch("platform_type", item, nil),
			"group_dn":       utils.PathSearch("group_dn", item, nil),
			"domain":         utils.PathSearch("domain", item, nil),
			"sid":            utils.PathSearch("sid", item, nil),
			"total_desktops": utils.PathSearch("total_desktops", item, nil),
		},
	}
}

func flattenUserGroups(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":   utils.PathSearch("id", item, nil),
			"name": utils.PathSearch("name", item, nil),
			"create_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
				item, "").(string))/1000, false),
			"description":   utils.PathSearch("description", item, nil),
			"user_quantity": utils.PathSearch("user_quantity", item, nil),
			"parent": flattenUserGroupParent(utils.PathSearch("parent", item,
				make(map[string]interface{})).(map[string]interface{})),
			"realm_id":       utils.PathSearch("realm_id", item, nil),
			"platform_type":  utils.PathSearch("platform_type", item, nil),
			"group_dn":       utils.PathSearch("group_dn", item, nil),
			"domain":         utils.PathSearch("domain", item, nil),
			"sid":            utils.PathSearch("sid", item, nil),
			"total_desktops": utils.PathSearch("total_desktops", item, nil),
		}
	}

	return result
}

func dataSourceUserGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	userGroups, err := listUserGroups(client)
	if err != nil {
		return diag.Errorf("error querying Workspace user groups: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenUserGroups(userGroups)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
