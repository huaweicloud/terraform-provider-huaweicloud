package nat

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

// @API NAT POST /v3/{project_id}/transit-ips/resource_instances/action
func DataSourcePrivateTransitIpsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateTransitIpsByTagsRead,

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
						// JSON format
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

func flattenPrivateTransitIpsByTags(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	resources := make([]interface{}, len(resp))
	for i, v := range resp {
		resources[i] = map[string]interface{}{
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_detail": utils.JsonToString(utils.PathSearch("resource_detail", v, nil)),
			"tags":            flattenTransitIpsTagsResp(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return resources
}

func flattenTransitIpsTagsResp(tags []interface{}) []interface{} {
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

func buildPrivateTransitIpsByTagsBodyParams(d *schema.ResourceData, action string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":       action,
		"matches":      buildFilterMatchesBodyParams(d.Get("matches").([]interface{})),
		"tags":         buildFilterTagsBodyParams(d.Get("tags").([]interface{})),
		"not_tags":     buildFilterTagsBodyParams(d.Get("not_tags").([]interface{})),
		"tags_any":     buildFilterTagsBodyParams(d.Get("tags_any").([]interface{})),
		"not_tags_any": buildFilterTagsBodyParams(d.Get("not_tags_any").([]interface{})),
	}

	return bodyParams
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

func dataSourcePrivateTransitIpsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v3/{project_id}/transit-ips/resource_instances/action"
		product      = "nat"
		tagsAction   = d.Get("action").(string)
		offset       = 0
		allTrasitIps = make([]interface{}, 0)
		totalCount   int
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestBody := buildPrivateTransitIpsByTagsBodyParams(d, tagsAction)
	if tagsAction == "filter" {
		requestBody["limit"] = "1000"
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
			return diag.Errorf("error retrieving NAT private transit IPs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		if tagsAction == "count" {
			break
		}

		transitIps := utils.PathSearch("resources", respBody, []interface{}{}).([]interface{})
		if len(transitIps) == 0 {
			break
		}

		allTrasitIps = append(allTrasitIps, transitIps...)
		offset += len(transitIps)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenPrivateTransitIpsByTags(allTrasitIps)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
