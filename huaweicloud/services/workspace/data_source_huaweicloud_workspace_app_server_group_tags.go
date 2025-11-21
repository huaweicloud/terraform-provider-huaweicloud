package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/server-group/tags
// @API Workspace GET /v1/{project_id}/server-group/{server_group_id}/tags
func DataSourceAppServerGroupTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppServerGroupTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the server group tags are located.`,
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the server group to which the tags belong.`,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the tag.`,
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The value list of the tag.`,
						},
					},
				},
				Description: `The tag list of the server group.`,
			},
		},
	}
}

func queryServerGroupTags(client *golangsdk.ServiceClient, serverGroupId string) ([]interface{}, error) {
	var httpUrl string
	if serverGroupId != "" {
		httpUrl = "v1/{project_id}/server-group/{server_group_id}/tags"
	} else {
		httpUrl = "v1/{project_id}/server-group/tags"
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	if serverGroupId != "" {
		getPath = strings.ReplaceAll(getPath, "{server_group_id}", serverGroupId)
	}

	getOpts := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenServerGroupTags(tags []interface{}) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(tags))
	for _, item := range tags {
		// When querying the tags for all server groups, the return value is a list of strings.
		// When querying a specific server group, the return value is a single string.
		values := utils.PathSearch("values", item, nil)
		if values == nil {
			values = []string{utils.PathSearch("value", item, "").(string)}
		}

		tag := map[string]interface{}{
			"key":    utils.PathSearch("key", item, nil),
			"values": values,
		}
		result = append(result, tag)
	}

	return result
}

func dataSourceAppServerGroupTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating WorkSpace APP client: %s", err)
	}

	tags, err := queryServerGroupTags(client, d.Get("server_group_id").(string))
	if err != nil {
		return diag.Errorf("error querying Workspace APP server group tags: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tags", flattenServerGroupTags(tags)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
