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

func DataSourceTrackedResourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrackedResourceTagsRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Tag key name`,
			},
			"resource_deleted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to query deleted resources`,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
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
				Description: `Tag list`,
			},
		},
	}
}

type TrackedResourceTagsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newTrackedResourceTagsDSWrapper(d *schema.ResourceData, meta interface{}) *TrackedResourceTagsDSWrapper {
	return &TrackedResourceTagsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceTrackedResourceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newTrackedResourceTagsDSWrapper(d, meta)
	lisTraResTagRst, err := wrapper.ListTrackedResourceTags()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listTrackedResourceTagsToSchema(lisTraResTagRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/tracked-resources/tags
func (w *TrackedResourceTagsDSWrapper) ListTrackedResourceTags() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rms")
	if err != nil {
		return nil, err
	}

	uri := "/v1/resource-manager/domains/{domain_id}/tracked-resources/tags"
	uri = strings.ReplaceAll(uri, "{domain_id}", w.Config.DomainID)
	params := map[string]any{
		"key":              w.Get("key"),
		"resource_deleted": w.Get("resource_deleted"),
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

func (w *TrackedResourceTagsDSWrapper) listTrackedResourceTagsToSchema(body *gjson.Result) error {
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
