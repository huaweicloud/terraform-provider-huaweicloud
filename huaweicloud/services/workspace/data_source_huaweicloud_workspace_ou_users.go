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

// @API Workspace GET /v2/{project_id}/ou-users
func DataSourceOuUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOuUsersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the OU users are located.`,
			},
			"ou_dn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The distinguished name (DN) of the OU.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the user to which the OU belongs.`,
			},
			"has_existed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the user already exists in the user list.`,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the user.`,
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The expiration time of the user.`,
						},
						"has_existed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the user already exists in the user list.`,
						},
					},
				},
				Description: `The list of users that match the filter parameters.`,
			},
			"enable_create_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of users that can be created.`,
			},
		},
	}
}

func buildListOuUsersQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("&ou_dn=%v", d.Get("ou_dn"))

	if v, ok := d.GetOk("user_name"); ok {
		res = fmt.Sprintf("%s&user_name=%v", res, v)
	}

	// If the 'has_existed' parameter is not specified, query all users.
	rawConfig := d.GetRawConfig()
	if v := utils.GetNestedObjectFromRawConfig(rawConfig, "has_existed"); v != nil {
		res = fmt.Sprintf("%s&has_existed=%v", res, v)
	}

	return res
}

func listOuUsers(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/ou-users"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%v", listPath, limit)
	listPath += buildListOuUsersQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=utf-8"},
	}

	var enableCreateCount interface{}
	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, nil, err
		}

		users := utils.PathSearch("user_infos", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, users...)
		if len(users) < limit {
			enableCreateCount = utils.PathSearch("enable_create_count", respBody, nil)
			break
		}
		offset += len(users)
	}

	return result, enableCreateCount, nil
}

func dataSourceOuUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	users, enableCreateCount, err := listOuUsers(client, d)
	if err != nil {
		return diag.Errorf("error querying users under the OU: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("users", flattenOuUsers(users)),
		d.Set("enable_create_count", enableCreateCount),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOuUsers(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"name":         utils.PathSearch("user_name", item, nil),
			"expired_time": utils.PathSearch("expired_time", item, nil),
			"has_existed":  utils.PathSearch("has_existed", item, nil),
		})
	}

	return result
}
