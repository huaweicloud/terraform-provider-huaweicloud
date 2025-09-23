package cts

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

// @API CTS POST /v3/{project_id}/{resource_type}/resource-instances/filter
func DataSourceResourceTagsFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceTagsFilterRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource type to be queried.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tag list for filtering resources.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the tag key.`,
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the tag values.`,
						},
					},
				},
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the match conditions for filtering resources.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the match key.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the match value.`,
						},
					},
				},
			},

			// Attributes
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of resources that match the filter conditions.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource.`,
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The detailed information of the resource.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The tags associated with the resource.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The tag key.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The tag value.`,
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

func dataSourceResourceTagsFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	resources, err := queryResourceTagsFilterResources(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenResourceTagsFilterResources(
			utils.PathSearch("resources", resources, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryResourceTagsFilterResources(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	// Use fixed limit for internal pagination
	limit := 1000
	offset := 0
	var allResources []interface{}

	uri := "v3/{project_id}/{resource_type}/resource-instances/filter"
	path := client.Endpoint + uri
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{resource_type}", d.Get("resource_type").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestBody := map[string]interface{}{
		"tags":    buildResourceTagsFilterTags(d.Get("tags").([]interface{})),
		"matches": buildResourceTagsFilterMatches(d.Get("matches").([]interface{})),
		"limit":   limit,
	}

	for {
		requestBody["offset"] = offset
		requestOpt.JSONBody = utils.RemoveNil(requestBody)
		response, err := client.Request("POST", path, &requestOpt)
		if err != nil {
			return nil, fmt.Errorf("error querying CTS resources: %s", err)
		}

		respBody, err := utils.FlattenResponse(response)
		if err != nil {
			return nil, err
		}

		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		allResources = append(allResources, resources...)

		if len(resources) < limit {
			break
		}

		offset += len(resources)
	}

	return allResources, nil
}

func buildResourceTagsFilterTags(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, utils.RemoveNil(
			map[string]interface{}{
				"key":    utils.PathSearch("key", tag, ""),
				"values": utils.ExpandToStringList(utils.PathSearch("values", tag, []interface{}{}).([]interface{})),
			}),
		)
	}

	return result
}

func buildResourceTagsFilterMatches(matches []interface{}) []map[string]interface{} {
	if len(matches) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(matches))
	for _, match := range matches {
		result = append(result, utils.RemoveNil(
			map[string]interface{}{
				"key":   utils.PathSearch("key", match, ""),
				"value": utils.PathSearch("value", match, ""),
			}),
		)
	}

	return result
}

func flattenResourceTagsFilterResources(resources []interface{}) []map[string]interface{} {
	if len(resources) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"id":     utils.PathSearch("resource_id", resource, nil),
			"detail": utils.PathSearch("resource_detail", resource, nil),
			"name":   utils.PathSearch("resource_name", resource, nil),
			"tags": flattenResourceTagsFilterTags(
				utils.PathSearch("tags", resource, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenResourceTagsFilterTags(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}

	return result
}
