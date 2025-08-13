package dns

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

// @API APIG POST /v2/{project_id}/{resource_type}/resource_instances/action
func DataSourceDNSTagsFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSTagsFilterRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource type.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the list of the tags to be queried.`,
				Elem:        dnsQueriedTagsFilterSchema(),
			},
			"tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the list of the tags to be queried.`,
				Elem:        dnsQueriedTagsFilterSchema(),
			},
			"not_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the list of the tags to be queried.`,
				Elem:        dnsQueriedTagsFilterSchema(),
			},
			"not_tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the list of the tags to be queried.`,
				Elem:        dnsQueriedTagsFilterSchema(),
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the fields to be queried.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the key to be matched.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the value of the matching field.`,
						},
					},
				},
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates all dedicated resources that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the resource.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the resource.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the tag list associated with the resource.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the key of the resource tag.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the value of the resource tag.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dnsQueriedTagsFilterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the key of tag.",
			},
			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the list of values of the tag.",
			},
		},
	}
}

func dataSourceDNSTagsFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	httpUrl := "v2/{project_id}/{resource_type}/resource_instances/action"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{resource_type}", d.Get("resource_type").(string))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	result := make([]map[string]interface{}, 0)
	for {
		listOpt.JSONBody = utils.RemoveNil(buildGetDNSTagsBodyParams(d, offset))
		requestResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving DNS resources: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, flattenDNSTagsResp(resources)...)
		if len(resources) == 0 {
			break
		}

		offset += len(resources)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", result),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetDNSTagsBodyParams(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"limit":        100,
		"offset":       offset,
		"action":       "filter",
		"tags":         buildGetResourceTagsBodyParams(d.Get("tags").([]interface{})),
		"tags_any":     buildGetResourceTagsBodyParams(d.Get("tags_any").([]interface{})),
		"not_tags":     buildGetResourceTagsBodyParams(d.Get("not_tags").([]interface{})),
		"not_tags_any": buildGetResourceTagsBodyParams(d.Get("not_tags_any").([]interface{})),
		"matches":      buildGetResourceMatchesBodyParams(d.Get("matches").([]interface{})),
	}
	return bodyParams
}

func buildGetResourceTagsBodyParams(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(tags))
	for i, tagRaw := range tags {
		if tag, ok := tagRaw.(map[string]interface{}); ok {
			res[i] = map[string]interface{}{
				// require empty value
				"key":    tag["key"],
				"values": tag["values"],
			}
		}
	}
	return res
}

func buildGetResourceMatchesBodyParams(matches []interface{}) []map[string]interface{} {
	if len(matches) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(matches))
	for i, matchRaw := range matches {
		if match, ok := matchRaw.(map[string]interface{}); ok {
			res[i] = map[string]interface{}{
				// require empty value
				"key":   match["key"],
				"value": match["value"],
			}
		}
	}
	return res
}

func flattenDNSTagsResp(resources []interface{}) []map[string]interface{} {
	if len(resources) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(resources))
	for i, v := range resources {
		rst[i] = map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"tags":          flattenResourceResourceTagsResp(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return rst
}

func flattenResourceResourceTagsResp(tags []interface{}) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	rst := make([]interface{}, len(tags))
	for i, tag := range tags {
		rst[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		}
	}
	return rst
}
