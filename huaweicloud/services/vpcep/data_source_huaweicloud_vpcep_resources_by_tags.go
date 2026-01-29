package vpcep

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

// @API VPCEP POST /v1/{project_id}/{resource_type}/resource_instances/action
func DataSourceVpcepResourcesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcepResourcesByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"not_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"not_tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func flattenVpcepTags(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}
	return rst
}

func flattenVpcepResourcesByTagsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	resources := make([]interface{}, len(resp))
	for i, v := range resp {
		resources[i] = map[string]interface{}{
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"tags":          flattenVpcepTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return resources
}

func buildFilterTagsBodyParams(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(tags))
	for i, tagRaw := range tags {
		if tag, ok := tagRaw.(map[string]interface{}); ok {
			res[i] = map[string]interface{}{
				"key":    utils.PathSearch("key", tag, nil),
				"values": utils.PathSearch("values", tag, nil),
			}
		}
	}

	return res
}

func buildFilterMatchesBodyParams(matches []interface{}) []map[string]interface{} {
	if len(matches) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(matches))
	for i, matchRaw := range matches {
		if match, ok := matchRaw.(map[string]interface{}); ok {
			res[i] = map[string]interface{}{
				"key":   utils.PathSearch("key", match, nil),
				"value": utils.PathSearch("value", match, nil),
			}
		}
	}
	return res
}

func buildVpcepResourcesByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":       d.Get("action"),
		"tags":         buildFilterTagsBodyParams(d.Get("tags").([]interface{})),
		"tags_any":     buildFilterTagsBodyParams(d.Get("tags_any").([]interface{})),
		"not_tags":     buildFilterTagsBodyParams(d.Get("not_tags").([]interface{})),
		"not_tags_any": buildFilterTagsBodyParams(d.Get("not_tags_any").([]interface{})),
		"matches":      buildFilterMatchesBodyParams(d.Get("matches").([]interface{})),
	}

	if d.Get("without_any_tag").(bool) {
		bodyParams["without_any_tag"] = true
	}

	return bodyParams
}

func dataSourceVpcepResourcesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/{resource_type}/resource_instances/action"
		resourceType = d.Get("resource_type").(string)
		tagsAction   = d.Get("action").(string)
		offset       = 0
		allResources = make([]interface{}, 0)
		totalCount   int
	)

	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	requestPath := vpcepClient.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", vpcepClient.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{resource_type}", resourceType)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestBody := buildVpcepResourcesByTagsBodyParams(d)
	if tagsAction == "filter" {
		requestBody["limit"] = "1000"
	}

	for {
		if tagsAction == "filter" {
			requestBody["offset"] = offset
		}

		requestOpt.JSONBody = utils.RemoveNil(requestBody)
		resp, err := vpcepClient.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving VPCEP resources by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		if tagsAction == "count" {
			break
		}

		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})

		if len(resources) == 0 {
			break
		}

		allResources = append(allResources, resources...)
		offset += len(resources)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("resources", flattenVpcepResourcesByTagsResponseBody(allResources)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
