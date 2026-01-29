package elb

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

// @API ELB POST /v2.0/{project_id}/loadbalancers/resource_instances/action
func DataSourceElbLoadbalancersByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbLoadbalancersByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
						"super_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_detail": {
							Type:     schema.TypeString,
							Computed: true,
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

func flattenElbTags(tags []interface{}) []map[string]interface{} {
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

func flattenElbLoadbalancersByTagsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	resources := make([]interface{}, len(resp))
	for i, v := range resp {
		resources[i] = map[string]interface{}{
			"resource_name":     utils.PathSearch("resource_name", v, nil),
			"resource_id":       utils.PathSearch("resource_id", v, nil),
			"super_resource_id": utils.PathSearch("super_resource_id", v, nil),
			"resource_detail":   utils.PathSearch("resource_detail", v, nil),
			"tags":              flattenElbTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return resources
}

func buildFilterTagsBodyParams(rawTags []interface{}) []map[string]interface{} {
	if len(rawTags) == 0 {
		return nil
	}

	tags := make([]map[string]interface{}, len(rawTags))
	for i, raw := range rawTags {
		if tag, ok := raw.(map[string]interface{}); ok {
			tags[i] = map[string]interface{}{
				"key":    utils.PathSearch("key", tag, nil),
				"values": utils.PathSearch("values", tag, nil),
			}
		}
	}
	return tags
}

func buildFilterMatchesBodyParams(rawMatches []interface{}) []map[string]interface{} {
	if len(rawMatches) == 0 {
		return nil
	}

	matches := make([]map[string]interface{}, len(rawMatches))
	for i, raw := range rawMatches {
		if match, ok := raw.(map[string]interface{}); ok {
			matches[i] = map[string]interface{}{
				"key":   utils.PathSearch("key", match, nil),
				"value": utils.PathSearch("value", match, nil),
			}
		}
	}
	return matches
}

func buildElbLoadbalancersByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":  d.Get("action"),
		"matches": buildFilterMatchesBodyParams(d.Get("matches").([]interface{})),
		"tags":    buildFilterTagsBodyParams(d.Get("tags").([]interface{})),
	}

	return bodyParams
}

func dataSourceElbLoadbalancersByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2.0/{project_id}/loadbalancers/resource_instances/action"
		product    = "elb"
		tagsAction = d.Get("action").(string)
		offset     = 0
		result     = make([]interface{}, 0)
		totalCount int
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	requestBody := buildElbLoadbalancersByTagsBodyParams(d)
	if tagsAction == "filter" {
		requestBody["limit"] = 1000
	}
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {

		if tagsAction == "filter" {
			requestBody["offset"] = offset
		}

		requestOpt.JSONBody = utils.RemoveNil(requestBody)
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving loadbalancers by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		if tagsAction == "count" {
			break
		}

		loadbalancers := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		if len(loadbalancers) == 0 {
			break
		}

		result = append(result, loadbalancers...)
		offset += len(loadbalancers)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("resources", flattenElbLoadbalancersByTagsResponseBody(result)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
