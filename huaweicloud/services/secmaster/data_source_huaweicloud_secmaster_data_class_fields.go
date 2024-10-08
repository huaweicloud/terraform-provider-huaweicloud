// Generated by PMS #325
package secmaster

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

func DataSourceSecmasterDataClassFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterDataClassFieldsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID.`,
			},
			"data_class_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the data class ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the field name.`,
			},
			"is_built_in": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether it is built in SecMaster. The value can be **true** or **false**.`,
			},
			"mapping": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether to display in other places other the classification and mapping module.`,
			},
			"fields": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The field list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The field ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The field name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The field type, such as **short text**, **radio** and **grid**.`,
						},
						"business_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The business code of the field.`,
						},
						"business_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated service.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The field description.`,
						},
						"data_class_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data class name.`,
						},
						"subscribed_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The subscribed version.`,
						},
						"mapping": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to display in other places other the classification and mapping module.`,
						},
						"io_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The input and output types.`,
						},
						"field_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The field key.`,
						},
						"extra_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The additional JSON.`,
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The default value.`,
						},
						"is_built_in": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the field is built in SecMaster.`,
						},
						"business_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of associated service.`,
						},
						"used_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Which services are used by.`,
						},
						"target_api": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The target API.`,
						},
						"json_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The JSON mode.`,
						},
						"field_tooltip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The tool tip.`,
						},
						"display_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The display type.`,
						},
						"required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the field is required.`,
						},
						"case_sensitive": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the field is case sensitive.`,
						},
						"editabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the field can be edited.`,
						},
						"visible": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the field is visible.`,
						},
						"maintainabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the field can be maintained.`,
						},
						"searchabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the field is searchable mode.`,
						},
						"read_only": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the field is read-only.`,
						},
						"creatabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the field can be created.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The create time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator.`,
						},
						"modifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The modifier.`,
						},
						"creator_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator ID.`,
						},
						"modifier_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The modifier ID.`,
						},
					},
				},
			},
		},
	}
}

type DataClassFieldsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newDataClassFieldsDSWrapper(d *schema.ResourceData, meta interface{}) *DataClassFieldsDSWrapper {
	return &DataClassFieldsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceSecmasterDataClassFieldsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newDataClassFieldsDSWrapper(d, meta)
	lisDatFieRst, err := wrapper.ListDataclassFields()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listDataclassFieldsToSchema(lisDatFieRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/dataclasses/{dataclass_id}/fields
func (w *DataClassFieldsDSWrapper) ListDataclassFields() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "secmaster")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/workspaces/{workspace_id}/soc/dataclasses/{dataclass_id}/fields"
	uri = strings.ReplaceAll(uri, "{workspace_id}", w.Get("workspace_id").(string))
	uri = strings.ReplaceAll(uri, "{dataclass_id}", w.Get("data_class_id").(string))
	params := map[string]any{
		"name":        w.Get("name"),
		"is_built_in": w.Get("is_built_in"),
		"mapping":     w.Get("mapping"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("field_details", "offset", "limit", 100).
		Request().
		Result()
}

func (w *DataClassFieldsDSWrapper) listDataclassFieldsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("fields", schemas.SliceToList(body.Get("field_details"),
			func(fields gjson.Result) any {
				return map[string]any{
					"id":                 fields.Get("id").Value(),
					"name":               fields.Get("name").Value(),
					"type":               fields.Get("field_type").Value(),
					"business_code":      fields.Get("business_code").Value(),
					"business_type":      fields.Get("business_type").Value(),
					"description":        fields.Get("description").Value(),
					"data_class_name":    fields.Get("dataclass_name").Value(),
					"subscribed_version": fields.Get("cloud_pack_version").Value(),
					"mapping":            fields.Get("mapping").Value(),
					"io_type":            fields.Get("iu_type").Value(),
					"field_key":          fields.Get("field_key").Value(),
					"extra_json":         fields.Get("extra_json").Value(),
					"default_value":      fields.Get("default_value").Value(),
					"is_built_in":        fields.Get("is_built_in").Value(),
					"business_id":        fields.Get("business_id").Value(),
					"used_by":            fields.Get("used_by").Value(),
					"target_api":         fields.Get("target_api").Value(),
					"json_schema":        fields.Get("json_schema").Value(),
					"field_tooltip":      fields.Get("field_tooltip").Value(),
					"display_type":       fields.Get("display_type").Value(),
					"required":           fields.Get("required").Value(),
					"case_sensitive":     fields.Get("case_sensitive").Value(),
					"editabled":          fields.Get("editable").Value(),
					"visible":            fields.Get("visible").Value(),
					"maintainabled":      fields.Get("maintainable").Value(),
					"searchabled":        fields.Get("searchable").Value(),
					"read_only":          fields.Get("read_only").Value(),
					"creatabled":         fields.Get("creatable").Value(),
					"created_at":         fields.Get("create_time").Value(),
					"updated_at":         fields.Get("update_time").Value(),
					"creator":            fields.Get("creator_name").Value(),
					"modifier":           fields.Get("modifier_name").Value(),
					"creator_id":         fields.Get("creator_id").Value(),
					"modifier_id":        fields.Get("modifier_id").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
