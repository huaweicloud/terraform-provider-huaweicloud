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

func workspaceTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key of tag.",
			},
			"values": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The value of tag.",
			},
		},
	}
}

func workspaceTagSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key of tag.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of tag.",
			},
		},
	}
}

func workspaceMatchSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of desktop property.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of desktop property.",
			},
		},
	}
}

// @API Workspace POST /v2/{project_id}/desktops/resource_instances/action
func DataSourceWorkspaceDesktopTagsFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkspaceDesktopTagsFilterRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the desktop tags are located.`,
			},
			"without_any_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to query resources without any tag.`,
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workspaceMatchSchema(),
				Description: `List of matching rules to filter desktops.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workspaceTagsSchema(),
				Description: `List of tags to filter desktops. Resources must contain all specified tags.`,
			},
			"tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workspaceTagsSchema(),
				Description: `List of tags to filter desktops. Resources must contain at least one of specified tags.`,
			},
			"not_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workspaceTagsSchema(),
				Description: `List of tags to filter desktops. Resources must not contain specified tags.`,
			},
			"not_tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workspaceTagsSchema(),
				Description: `List of tags to filter desktops. Resources must not contain any of specified tags.`,
			},
			"desktops": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the desktop.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the desktop.`,
						},
						"resource_detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The detail of the desktop.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        workspaceTagSchema(),
							Description: `The tags of the desktop.`,
						},
					},
				},
				Description: `List of desktops that match the filter criteria.`,
			},
		},
	}
}

func buildTagsFilter(d *schema.ResourceData) (map[string]interface{}, error) {
	filter := map[string]interface{}{
		"action": "filter",
		"offset": "0",
		"limit":  "1000",
	}

	if v, ok := d.GetOk("without_any_tag"); ok {
		filter["without_any_tag"] = v.(bool)
	}

	if v, ok := d.GetOk("matches"); ok {
		matches := make([]map[string]interface{}, 0)
		for _, item := range v.([]interface{}) {
			match := item.(map[string]interface{})
			matches = append(matches, map[string]interface{}{
				"key":   match["key"].(string),
				"value": match["value"].(string),
			})
		}
		filter["matches"] = matches
	}

	if v, ok := d.GetOk("tags"); ok {
		filter["tags"] = buildTagsCondition(v.([]interface{}))
	}

	if v, ok := d.GetOk("tags_any"); ok {
		filter["tags_any"] = buildTagsCondition(v.([]interface{}))
	}

	if v, ok := d.GetOk("not_tags"); ok {
		filter["not_tags"] = buildTagsCondition(v.([]interface{}))
	}

	if v, ok := d.GetOk("not_tags_any"); ok {
		filter["not_tags_any"] = buildTagsCondition(v.([]interface{}))
	}

	return filter, nil
}

func buildTagsCondition(tagsList []interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, 0)
	for _, item := range tagsList {
		tag := item.(map[string]interface{})
		values := make([]string, 0)
		for _, v := range tag["values"].([]interface{}) {
			values = append(values, v.(string))
		}
		tags = append(tags, map[string]interface{}{
			"key":    tag["key"].(string),
			"values": values,
		})
	}
	return tags
}

func flattenDesktopTags(tags []interface{}) []interface{} {
	if len(tags) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		tagMap := tag.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", tagMap, nil),
			"value": utils.PathSearch("value", tagMap, nil),
		})
	}
	return result
}

func flattenDesktops(resources []interface{}) []interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(resources))
	for _, item := range resources {
		resource := item.(map[string]interface{})
		desktop := map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", resource, nil),
			"resource_name":   utils.PathSearch("resource_name", resource, nil),
			"resource_detail": utils.PathSearch("resource_detail", resource, nil),
			"tags":            flattenDesktopTags(utils.PathSearch("tags", resource, nil).([]interface{})),
		}
		result = append(result, desktop)
	}
	return result
}

func filterDesktopsByTags(client *golangsdk.ServiceClient, filter map[string]interface{}) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/desktops/resource_instances/action"
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	requestOpts := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		JSONBody: filter,
	}

	requestResp, err := client.Request("POST", listPath, requestOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	desktops := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
	return desktops, nil
}

func dataSourceWorkspaceDesktopTagsFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	filter, err := buildTagsFilter(d)
	if err != nil {
		return diag.Errorf("error building filter condition: %s", err)
	}

	desktops, err := filterDesktopsByTags(client, filter)
	if err != nil {
		return diag.Errorf("error querying desktop tags: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("desktops", flattenDesktops(desktops)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
