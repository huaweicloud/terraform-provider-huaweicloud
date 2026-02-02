package tms

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceTmsResourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsResourceTagsRead,

		Schema: map[string]*schema.Schema{
			"resource_types": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource type.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the project ID.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        tmsResourceTagsSchema(),
				Description: `Indicates the tags.`,
			},
			"errors": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        tmsResourceTagsErrorsSchema(),
				Description: `Indicates the tag error.`,
			},
		},
	}
}

func tmsResourceTagsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func tmsResourceTagsErrorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

type TagListDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newTagListDSWrapper(d *schema.ResourceData, meta interface{}) *TagListDSWrapper {
	return &TagListDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceTmsResourceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newTagListDSWrapper(d, meta)
	listTagsRst, err := wrapper.ListTags()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listTagsToSchema(listTagsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API TMS GET /v1.0/tags
func (w *TagListDSWrapper) ListTags() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "tms")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/tags"
	params := map[string]any{
		"resource_types": w.Get("resource_types"),
		"project_id":     w.Get("project_id"),
		"tag_types":      "resource",
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *TagListDSWrapper) listTagsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("tags", schemas.SliceToList(body.Get("tags"),
			func(tags gjson.Result) any {
				return map[string]any{
					"key":    tags.Get("key").Value(),
					"values": tags.Get("values").Value(),
				}
			},
		)),
		d.Set("errors", schemas.SliceToList(body.Get("errors"),
			func(errors gjson.Result) any {
				return map[string]any{
					"project_id":    errors.Get("project_id").Value(),
					"resource_type": errors.Get("resource_type").Value(),
					"error_code":    errors.Get("error_code").Value(),
					"error_msg":     errors.Get("error_msg").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
