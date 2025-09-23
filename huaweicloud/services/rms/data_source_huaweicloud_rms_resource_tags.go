package rms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceResourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource type.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource ID.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The tags.`,
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
							Description: `The tag values.`,
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

func dataSourceResourceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newTagsDSWrapper(d, meta)
	if _, ok := d.GetOk("resource_id"); ok {
		lisTagForResRst, err := wrapper.ListTagsForResource()
		if err != nil {
			return diag.FromErr(err)
		}

		err = wrapper.listTagsToSchema(lisTagForResRst, true)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		lisTagForResTypRst, err := wrapper.ListTagsForResourceType()
		if err != nil {
			return diag.FromErr(err)
		}

		err = wrapper.listTagsToSchema(lisTagForResTypRst, false)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return nil
}

// @API CONFIG GET /v1/resource-manager/{resource_type}/{resource_id}/tags
func (w *TagsDSWrapper) ListTagsForResource() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rms")
	if err != nil {
		return nil, err
	}

	uri := "/v1/resource-manager/{resource_type}/{resource_id}/tags"
	uri = strings.ReplaceAll(uri, "{resource_type}", w.Get("resource_type").(string))
	uri = strings.ReplaceAll(uri, "{resource_id}", w.Get("resource_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

// @API CONFIG GET /v1/resource-manager/{resource_type}/tags
func (w *TagsDSWrapper) ListTagsForResourceType() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rms")
	if err != nil {
		return nil, err
	}

	uri := "/v1/resource-manager/{resource_type}/tags"
	uri = strings.ReplaceAll(uri, "{resource_type}", w.Get("resource_type").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("tags", "offset", "limit", 1000).
		Request().
		Result()
}

func (w *TagsDSWrapper) listTagsToSchema(body *gjson.Result, specifiedResourceId bool) error {
	d := w.ResourceData
	parseTagsFunc := parseResourceTypeTags
	if specifiedResourceId {
		parseTagsFunc = parseResourceTags
	}
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("tags", schemas.SliceToList(body.Get("tags"), parseTagsFunc)),
	)

	return mErr.ErrorOrNil()
}

func parseResourceTags(tags gjson.Result) any {
	rawTag := tags.Get("value").Value().(string)
	return map[string]any{
		"key":    tags.Get("key").Value(),
		"values": []string{rawTag},
	}
}

func parseResourceTypeTags(tags gjson.Result) any {
	return map[string]any{
		"key":    tags.Get("key").Value(),
		"values": schemas.SliceToStrList(tags.Get("values")),
	}
}
