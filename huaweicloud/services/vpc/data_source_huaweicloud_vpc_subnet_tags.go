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

func DataSourceVpcSubnetTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcSubnetTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
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

type SubnetTagsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newSubnetTagsDSWrapper(d *schema.ResourceData, meta interface{}) *SubnetTagsDSWrapper {
	return &SubnetTagsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceVpcSubnetTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newSubnetTagsDSWrapper(d, meta)
	listSubnetTagRst, err := wrapper.ListSubnetTags()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listSubnetTagsToSchema(listSubnetTagRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API VPC GET /v2.0/{project_id}/subnets/tags
func (w *SubnetTagsDSWrapper) ListSubnetTags() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "vpc")
	if err != nil {
		return nil, err
	}

	uri := "/v2.0/{project_id}/subnets/tags"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *SubnetTagsDSWrapper) listSubnetTagsToSchema(body *gjson.Result) error {
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
