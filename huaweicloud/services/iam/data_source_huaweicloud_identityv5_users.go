package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/users
func DataSourceIdentityV5Users() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5UserRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the IAM user.",
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_root_user": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5UserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allUsers []interface{}
	var marker string
	var path string
	for {
		path = fmt.Sprintf("%sv5/users", client.Endpoint) + buildListUsersV5Params(d, marker)
		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving users")
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}
		users := flattenListUsersV5Response(resp)
		allUsers = append(allUsers, users...)

		marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
		if marker == "" {
			break
		}
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("users", allUsers),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting users fields: %s", err)
	}
	return nil
}

func buildListUsersV5Params(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("group_id"); ok {
		res = fmt.Sprintf("%s&group_id=%v", res, v)
	}
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func flattenListUsersV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	users := utils.PathSearch("users", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = map[string]interface{}{
			"enabled":      utils.PathSearch("enabled", user, nil),
			"user_name":    utils.PathSearch("user_name", user, nil),
			"description":  utils.PathSearch("description", user, nil),
			"is_root_user": utils.PathSearch("is_root_user", user, nil),
			"created_at":   utils.PathSearch("created_at", user, nil),
			"urn":          utils.PathSearch("urn", user, nil),
			"user_id":      utils.PathSearch("user_id", user, nil),
		}
	}
	return result
}
