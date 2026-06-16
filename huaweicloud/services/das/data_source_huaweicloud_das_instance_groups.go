package das

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/batch-inspection/instance-group
func DataSourceInstanceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the DAS instance groups are located.`,
			},

			// Required parameters.
			"datastore_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The database type.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The instance group name.`,
			},

			// Attributes.
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the instance groups that matched the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The instance group ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance group name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the instance group.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceInstanceGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	queryParams := buildInstanceGroupsQueryParams(d)
	resp, err := ListInstanceGroups(client, queryParams)
	if err != nil {
		return diag.Errorf("error querying DAS instance groups: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenInstanceGroups(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildInstanceGroupsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("&datastore_type=%v", d.Get("datastore_type"))

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&group_name=%v", res, v)
	}

	return res
}

// ListInstanceGroups queries the instance groups with optional query parameters.
func ListInstanceGroups(client *golangsdk.ServiceClient, queryParams string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/batch-inspection/instance-group?limit={limit}"
		limit   = 200
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += queryParams

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
		instanceGroups := utils.PathSearch("instance_group_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, instanceGroups...)
		if len(instanceGroups) < limit {
			break
		}
		offset += len(instanceGroups)
	}

	return result, nil
}

func flattenInstanceGroups(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":          utils.PathSearch("group_id", item, nil),
			"name":        utils.PathSearch("group_name", item, nil),
			"description": utils.PathSearch("description", item, nil),
		}
	}

	return result
}
