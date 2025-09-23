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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceRmsTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRmsTagsRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Tag key name`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Tag list`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Tag key`,
						},
						"value": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Tag value list`,
						},
					},
				},
			},
		},
	}
}

type TagsV2DSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newTagsV2DSWrapper(d *schema.ResourceData, meta interface{}) *TagsV2DSWrapper {
	return &TagsV2DSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceRmsTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newTagsV2DSWrapper(d, meta)
	listAllTagsRst, err := wrapper.ListAllTags()

	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAllTagsToSchema(listAllTagsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/all-resources/tags
func (w *TagsV2DSWrapper) ListAllTags() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rms")
	if err != nil {
		return nil, err
	}

	uri := "/v1/resource-manager/domains/{domain_id}/all-resources/tags"
	uri = strings.ReplaceAll(uri, "{domain_id}", w.Config.DomainID)
	params := map[string]any{
		"key": w.Get("key"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("tags", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *TagsV2DSWrapper) listAllTagsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("tags", schemas.SliceToList(body.Get("tags"),
			func(tags gjson.Result) any {
				return map[string]any{
					"key":   tags.Get("key").Value(),
					"value": schemas.SliceToStrList(tags.Get("value")),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
