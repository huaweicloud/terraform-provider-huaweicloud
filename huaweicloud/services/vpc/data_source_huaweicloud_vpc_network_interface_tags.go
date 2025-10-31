package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceVpcNetworkInterfaceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcNetworkInterfaceTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource. If omitted, the provider-level region will be used.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of tags`,
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

type NetworkInterfaceTagsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newNetworkInterfaceTagsDSWrapper(d *schema.ResourceData, meta interface{}) *NetworkInterfaceTagsDSWrapper {
	return &NetworkInterfaceTagsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceVpcNetworkInterfaceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newNetworkInterfaceTagsDSWrapper(d, meta)
	listPortTagsRst, err := wrapper.ListPortTags()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.ListPortTagsToSchema(listPortTagsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API VPC GET /v3/{project_id}/ports/tags
func (w *NetworkInterfaceTagsDSWrapper) ListPortTags() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "vpc")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/ports/tags"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("tags", "offset", "limit", 0).
		Request().
		Result()
}

func (w *NetworkInterfaceTagsDSWrapper) ListPortTagsToSchema(body *gjson.Result) error {
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
