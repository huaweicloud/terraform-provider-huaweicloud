package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/groups
func DataSourceIdentityV5Groups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5GroupsRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5GroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allGroups []interface{}
	var marker string
	var path string

	for {
		path = client.Endpoint + "v5/groups" + buildListGroupsV5Params(d, marker)
		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving groups: %s", err)
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}

		groups := flattenListGroupsV5Response(resp)
		allGroups = append(allGroups, groups...)

		marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
		if marker == "" {
			break
		}
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("groups", allGroups),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListGroupsV5Params(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	if v, ok := d.GetOk("user_id"); ok {
		res = fmt.Sprintf("%s&user_id=%v", res, v)
	}
	return res
}

func flattenListGroupsV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	groups := utils.PathSearch("groups", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(groups))
	for i, group := range groups {
		result[i] = map[string]interface{}{
			"group_name":  utils.PathSearch("group_name", group, nil),
			"urn":         utils.PathSearch("urn", group, nil),
			"created_at":  utils.PathSearch("created_at", group, nil),
			"description": utils.PathSearch("description", group, nil),
			"group_id":    utils.PathSearch("group_id", group, nil),
		}
	}
	return result
}
