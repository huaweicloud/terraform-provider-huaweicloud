package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/groups
func DataSourceV5Groups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5GroupsRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the user.`,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the user group.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the user group.`,
						},
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The uniform resource name of the user group.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the user group.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the user group.`,
						},
					},
					Description: `The list of user groups that match the filter parameters.`,
				},
			},
		},
	}
}

func buildV5GroupsQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("user_id"); ok {
		return fmt.Sprintf("&user_id=%v", v)
	}

	return ""
}

func listV5Groups(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v5/groups"
		limit   = 100
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	listPath += buildV5GroupsQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker += fmt.Sprintf("&marker=%s", marker)
		}
		resp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		groups := utils.PathSearch("groups", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, groups...)
		if len(groups) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}
	return result, nil
}

func dataSourceV5GroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groups, err := listV5Groups(client, d)
	if err != nil {
		return diag.Errorf("error retrieving groups: %s", err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomId)

	return diag.FromErr(d.Set("groups", flattenV5Groups(groups)))
}

func flattenV5Groups(groups []interface{}) []interface{} {
	if len(groups) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(groups))
	for _, group := range groups {
		result = append(result, map[string]interface{}{
			"group_name":  utils.PathSearch("group_name", group, nil),
			"urn":         utils.PathSearch("urn", group, nil),
			"created_at":  utils.PathSearch("created_at", group, nil),
			"description": utils.PathSearch("description", group, nil),
			"group_id":    utils.PathSearch("group_id", group, nil),
		})
	}
	return result
}
