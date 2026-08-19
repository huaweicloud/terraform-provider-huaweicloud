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

// @API IAM GET /v3/groups
func DataSourceV3Groups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3GroupsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the user group to be queried.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the account to which the user groups belong.`,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the user group.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the account to which the user group belongs.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the user group, in RFC3339 format.`,
						},
					},
					Description: `The list of user groups that match the filter parameters.`,
				},
			},
		},
	}
}

func buildV3GroupsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("domain_id"); ok {
		res = fmt.Sprintf("%s&domain_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

// The current interface does not support pagination parameters.
func listV3Groups(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var httpUrl = "v3/groups"

	listPath := client.Endpoint + httpUrl + buildV3GroupsQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	groups := utils.PathSearch("groups", respBody, make([]interface{}, 0)).([]interface{})
	return groups, nil
}

func flattenV3Groups(groups []interface{}) []interface{} {
	if len(groups) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(groups))
	for _, group := range groups {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", group, nil),
			"name":        utils.PathSearch("name", group, nil),
			"description": utils.PathSearch("description", group, nil),
			"domain_id":   utils.PathSearch("domain_id", group, nil),
			"created_at":  utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", group, float64(0)).(float64))/1000, false),
		})
	}
	return result
}

func dataSourceV3GroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groups, err := listV3Groups(client, d)
	if err != nil {
		return diag.Errorf("error retrieving groups: %s", err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomId)

	return diag.FromErr(d.Set("groups", flattenV3Groups(groups)))
}
