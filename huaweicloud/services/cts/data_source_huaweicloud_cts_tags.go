package cts

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCtsTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCtsTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},

			// Attribute
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of tags.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The tag key.`,
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of tag values.`,
						},
					},
				},
			},
		},
	}
}

type TagsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newTagsDSWrapper(d *schema.ResourceData, meta interface{}) *TagsDSWrapper {
	return &TagsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCtsTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newTagsDSWrapper(d, meta)
	listTagsRst, err := wrapper.ListTags()
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(randUUID)

	err = wrapper.listTagsToSchema(listTagsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CTS GET /v3/{project_id}/{resource_type}/tags
func (w *TagsDSWrapper) ListTags() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cts")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/{resource_type}/tags"
	uri = strings.ReplaceAll(uri, "{resource_type}", "cts-tracker")
	var allTags []interface{}
	marker := ""
	limit := 200 // API maximum limit

	for {
		params := map[string]any{
			"limit": limit,
		}
		if marker != "" {
			params["marker"] = marker
		}
		params = utils.RemoveNil(params)

		result, err := httphelper.New(client).
			Method("GET").
			URI(uri).
			Query(params).
			OkCode(200).
			Request().
			Result()
		if err != nil {
			return nil, err
		}

		pageTags := result.Get("tags").Array()
		for _, tag := range pageTags {
			allTags = append(allTags, tag.Value())
		}

		nextMarker := result.Get("page_info.next_marker").String()
		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}

	combinedResponse := map[string]interface{}{
		"tags": allTags,
	}

	jsonBytes, err := json.Marshal(combinedResponse)
	if err != nil {
		return nil, err
	}

	combinedResult := gjson.ParseBytes(jsonBytes)
	return &combinedResult, nil
}

func (w *TagsDSWrapper) listTagsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("tags", schemas.SliceToList(body.Get("tags"),
			func(tags gjson.Result) any {
				return map[string]any{
					"key":    tags.Get("key").Value(),
					"values": schemas.SliceToStrList(tags.Get("values")),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
