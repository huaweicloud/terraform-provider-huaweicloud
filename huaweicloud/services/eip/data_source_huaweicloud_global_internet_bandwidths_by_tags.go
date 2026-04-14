package eip

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceGlobalInternetBandwidthsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalInternetBandwidthsByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
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
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_detail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
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
		},
	}
}

func globalInternetBandwidthsTagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

type GlobalInternetBandwidthsByTagsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newGlobalInternetBandwidthsByTagsDSWrapper(d *schema.ResourceData, meta interface{}) *GlobalInternetBandwidthsByTagsDSWrapper {
	return &GlobalInternetBandwidthsByTagsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

// @API EIP POST /v3/internet-bandwidth/resource-instances/filter
func dataSourceGlobalInternetBandwidthsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newGlobalInternetBandwidthsByTagsDSWrapper(d, meta)
	result, err := wrapper.ListInternetBandwidthsByTags()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.listInternetBandwidthsByTagsToSchema(result)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

func (w *GlobalInternetBandwidthsByTagsDSWrapper) ListInternetBandwidthsByTags() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "geip")
	if err != nil {
		return nil, err
	}

	var tags []map[string]interface{}
	if v, ok := w.GetOk("tags"); ok {
		tagList, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid tags format")
		}
		tags = make([]map[string]interface{}, 0, len(tagList))
		for _, item := range tagList {
			if tagMap, ok := item.(map[string]interface{}); ok {
				tags = append(tags, map[string]interface{}{
					"key":   tagMap["key"],
					"value": tagMap["value"],
				})
			}
		}
	}

	if tags == nil {
		tags = make([]map[string]interface{}, 0)
	}

	params := map[string]any{
		"tags": tags,
	}

	var allResources []interface{}
	var requestID string
	offset := 0
	limit := 100
	totalCount := 0

	for {
		uri := fmt.Sprintf("/v3/internet-bandwidth/resource-instances/filter?limit=%d&offset=%d", limit, offset)

		result, err := httphelper.New(client).
			Method("POST").
			URI(uri).
			Body(params).
			Request().
			Result()
		if err != nil {
			return nil, err
		}

		if requestID == "" {
			requestID = result.Get("request_id").String()
		}

		resources := result.Get("resources").Array()
		if len(resources) == 0 {
			break
		}

		for _, resource := range resources {
			allResources = append(allResources, resource.Value())
		}

		if len(resources) < limit {
			break
		}

		if totalCount == 0 {
			totalCount = int(result.Get("total_count").Int())
		}

		offset += limit

		if offset >= totalCount || len(resources) < limit {
			break
		}
	}

	mergedResult := map[string]interface{}{
		"request_id": requestID,
		"resources":  allResources,
	}

	jsonBytes, err := json.Marshal(mergedResult)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged result: %s", err)
	}
	parsedResult := gjson.ParseBytes(jsonBytes)
	return &parsedResult, nil
}

func (w *GlobalInternetBandwidthsByTagsDSWrapper) listInternetBandwidthsByTagsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	if !body.Get("resources").Exists() {
		return fmt.Errorf("unable to find resources in API response")
	}
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("request_id", body.Get("request_id").Value()),
		d.Set("resources", schemas.SliceToList(body.Get("resources"),
			func(resources gjson.Result) any {
				return map[string]any{
					"resource_id":     resources.Get("resource_id").Value(),
					"resource_detail": resources.Get("resource_detail").Value(),
					"resource_name":   resources.Get("resource_name").Value(),
					"tags": schemas.SliceToList(resources.Get("tags"),
						func(tag gjson.Result) any {
							return map[string]any{
								"key":   tag.Get("key").Value(),
								"value": tag.Get("value").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
