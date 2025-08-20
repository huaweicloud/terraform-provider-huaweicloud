package fgs

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

// @API FunctionGraph POST /v2/{project_id}/{resource_type}/resource-instances/{action}
func DataSourceResourcesFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcesFilterRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the target resources are located.`,
			},

			// Required parameter(s).
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the resource used to filter the target resources.`,
			},

			// Optional parameter(s).
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The match key used to filter the target resources.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The match value used to filter the target resources.`,
						},
					},
				},
				Description: `The key-value pairs used to filter the target resources.`,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The key of the resource tag used to filter the target resources.`,
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The values corresponding to the current key used to filter the target resources.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The resource tags used to filter the target resources.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},

			// Internal parameter(s).
			"sys_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The key of the system tag used to filter the target resources.`,
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The values of the system tag used to filter the target resources.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The system tags used to filter the target resources.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},

			// Attributes.
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource.`,
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The detailed information of the resource, in JSON format.`,
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
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of the tag.`,
									},
								},
							},
							Description: `The tags of the resource.`,
						},

						// Internal attribute(s).
						"sys_tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The key of the system tag.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of the system tag.`,
									},
								},
							},
							Description: utils.SchemaDesc(
								`The system tags of the resource.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				Description: `The list of target resources that matched filter parameters.`,
			},
		},
	}
}

func buildResourcesFilterBodyParams(d *schema.ResourceData, actionType string, limit int, offset int) map[string]interface{} {
	result := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
		"action": actionType,
	}

	if v, ok := d.GetOk("matches"); ok {
		result["matches"] = buildResourcesFilterMatches(v.([]interface{}))
	}
	if v, ok := d.GetOk("tags"); ok {
		result["tags"] = buildResourcesFilterTags(v.([]interface{}))
	}

	// Internal parameter(s).
	if v, ok := d.GetOk("sys_tags"); ok {
		result["sys_tags"] = buildResourcesFilterTags(v.([]interface{}))
	}

	return result
}

func buildResourcesFilterMatches(matches []interface{}) []map[string]interface{} {
	if len(matches) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(matches))
	for _, match := range matches {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", match, nil),
			"value": utils.PathSearch("value", match, nil),
		})
	}

	return result
}

func buildResourcesFilterTags(tags []interface{}) []map[string]interface{} {
	if len(tags) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":    utils.PathSearch("key", tag, nil),
			"values": utils.PathSearch("values", tag, nil),
		})
	}

	return result
}

func listResourceInstances(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl          = "v2/{project_id}/{resource_type}/resource-instances/{action}"
		result           = make([]interface{}, 0)
		actionTypeFilter = "filter"
		limit            = 100
		offset           = 0
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{resource_type}", d.Get("resource_type").(string))
	listPath = strings.ReplaceAll(listPath, "{action}", actionTypeFilter)

	for {
		listOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildResourcesFilterBodyParams(d, actionTypeFilter, limit, offset),
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}

		requestResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, resources...)
		if len(resources) < limit {
			break
		}
		offset += limit
	}
	return result, nil
}

func flattenResourcesFilterTags(tags []interface{}) []map[string]interface{} {
	if len(tags) < 1 {
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

func flattenResourcesFilterSysTags(sysTags []interface{}) []map[string]interface{} {
	if len(sysTags) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(sysTags))
	for _, sysTag := range sysTags {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", sysTag, nil),
			"value": utils.PathSearch("value", sysTag, nil),
		})
	}

	return result
}

func flattenResourcesFilterResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"id":     utils.PathSearch("resource_id", resource, nil),
			"name":   utils.PathSearch("resource_name", resource, nil),
			"detail": utils.JsonToString(utils.PathSearch("resource_detail", resource, nil)),
			"tags": flattenResourcesFilterTags(utils.PathSearch("tags", resource,
				make([]interface{}, 0)).([]interface{})),
			// Internal attribute(s).
			"sys_tags": flattenResourcesFilterSysTags(utils.PathSearch("sys_tags", resource,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceResourcesFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	resources, err := listResourceInstances(client, d)
	if err != nil {
		return diag.Errorf("error querying FunctionGraph resource instances: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenResourcesFilterResources(resources)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
