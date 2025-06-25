package workspace

import (
	"context"
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

// @API Workspace POST /v2/{project_id}/desktops/resource_instances/action
func DataSourceDesktopTagsFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDesktopTagsFilterRead,

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
				Description: `Whether to query resources without any tag.`,
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        desktopQueriedMatchFilterSchema(),
				Description: `The list of matching rules to filter desktops.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        desktopQueriedTagsFilterSchema(),
				Description: `The list of tags to filter desktops. Resources must contain all specified tags.`,
			},
			"tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        desktopQueriedTagsFilterSchema(),
				Description: `The list of tags to filter desktops. Resources must contain at least one of specified tags.`,
			},
			"not_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        desktopQueriedTagsFilterSchema(),
				Description: `The list of tags to filter desktops. Resources must not contain specified tags.`,
			},
			"not_tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        desktopQueriedTagsFilterSchema(),
				Description: `The list of tags to filter desktops. Resources must not contain any of specified tags.`,
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
							Elem:        desktopQueriedTagResponseSchema(),
							Description: `The list of tags attached to the desktop.`,
						},
					},
				},
				Description: `The list of desktops that match the filter parameters.`,
			},
		},
	}
}

func desktopQueriedTagsFilterSchema() *schema.Resource {
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
				Description: "The list of tag values that matched corresponding key.",
			},
		},
	}
}

func desktopQueriedMatchFilterSchema() *schema.Resource {
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

func desktopQueriedTagResponseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key of tag.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of tag.",
			},
		},
	}
}

func buildDesktopTagsFilter(d *schema.ResourceData) map[string]interface{} {
	filter := map[string]interface{}{
		"action": "filter",
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

	return filter
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

func flattenFilteredDesktop(resources []interface{}) []interface{} {
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

func filterDesktopsByTags(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/desktops/resource_instances/action"
		offset  = 0
		limit   = 200
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	filter := buildDesktopTagsFilter(d)

	for {
		filter["offset"] = strconv.Itoa(offset)
		filter["limit"] = strconv.Itoa(limit)

		requestOpts := &golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         filter,
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
		result = append(result, desktops...)
		if len(desktops) < limit {
			break
		}
		offset += len(desktops)
	}
	return result, nil
}

func dataSourceDesktopTagsFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating WorkSpace client: %s", err)
	}

	desktops, err := filterDesktopsByTags(client, d)
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
		d.Set("desktops", flattenFilteredDesktop(desktops)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
