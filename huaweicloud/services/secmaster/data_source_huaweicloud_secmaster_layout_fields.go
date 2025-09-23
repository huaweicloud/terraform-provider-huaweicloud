package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields
func DataSourceLayoutFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLayoutFieldsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"business_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_built_in": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"layout_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fields": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_pack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_pack_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataclass_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_pack_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"en_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"en_default_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extra_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field_tooltip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"en_field_tooltip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"json_schema": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_built_in": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"searchable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"visible": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"maintainable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"editable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creatable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wizard_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aopworkflow_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aopworkflow_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"playbook_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"playbook_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boa_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildLayoutFieldsQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?business_code=%v", d.Get("business_code"))

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("is_built_in"); ok {
		queryParams = fmt.Sprintf("%s&is_built_in=%v", queryParams, v)
	}
	if v, ok := d.GetOk("layout_id"); ok {
		queryParams = fmt.Sprintf("%s&layout_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceLayoutFieldsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v2/{project_id}/workspaces/{workspace_id}/soc/layouts/fields"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath += buildLayoutFieldsQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving layout fields: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("fields", flattenLayoutFields(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLayoutFields(fieldsResp interface{}) []interface{} {
	if fieldsArray, ok := fieldsResp.([]interface{}); ok {
		rst := make([]interface{}, 0, len(fieldsArray))
		for _, v := range fieldsArray {
			rst = append(rst, map[string]interface{}{
				"id":                     utils.PathSearch("id", v, nil),
				"cloud_pack_id":          utils.PathSearch("cloud_pack_id", v, nil),
				"cloud_pack_name":        utils.PathSearch("cloud_pack_name", v, nil),
				"dataclass_id":           utils.PathSearch("dataclass_id", v, nil),
				"cloud_pack_version":     utils.PathSearch("cloud_pack_version", v, nil),
				"field_key":              utils.PathSearch("field_key", v, nil),
				"name":                   utils.PathSearch("name", v, nil),
				"description":            utils.PathSearch("description", v, nil),
				"en_description":         utils.PathSearch("en_description", v, nil),
				"default_value":          utils.PathSearch("default_value", v, nil),
				"en_default_value":       utils.PathSearch("en_default_value", v, nil),
				"field_type":             utils.PathSearch("field_type", v, nil),
				"extra_json":             utils.PathSearch("extra_json", v, nil),
				"field_tooltip":          utils.PathSearch("field_tooltip", v, nil),
				"en_field_tooltip":       utils.PathSearch("en_field_tooltip", v, nil),
				"json_schema":            utils.PathSearch("json_schema", v, nil),
				"is_built_in":            utils.PathSearch("is_built_in", v, nil),
				"read_only":              utils.PathSearch("read_only", v, nil),
				"required":               utils.PathSearch("required", v, nil),
				"searchable":             utils.PathSearch("searchable", v, nil),
				"visible":                utils.PathSearch("visible", v, nil),
				"maintainable":           utils.PathSearch("maintainable", v, nil),
				"editable":               utils.PathSearch("editable", v, nil),
				"creatable":              utils.PathSearch("creatable", v, nil),
				"creator_id":             utils.PathSearch("creator_id", v, nil),
				"creator_name":           utils.PathSearch("creator_name", v, nil),
				"modifier_id":            utils.PathSearch("modifier_id", v, nil),
				"modifier_name":          utils.PathSearch("modifier_name", v, nil),
				"create_time":            utils.PathSearch("create_time", v, nil),
				"update_time":            utils.PathSearch("update_time", v, nil),
				"wizard_id":              utils.PathSearch("wizard_id", v, nil),
				"aopworkflow_id":         utils.PathSearch("aopworkflow_id", v, nil),
				"aopworkflow_version_id": utils.PathSearch("aopworkflow_version_id", v, nil),
				"playbook_id":            utils.PathSearch("playbook_id", v, nil),
				"playbook_version_id":    utils.PathSearch("playbook_version_id", v, nil),
				"boa_version":            utils.PathSearch("boa_version", v, nil),
				"version":                utils.PathSearch("version", v, nil),
			})
		}
		return rst
	}

	return make([]interface{}, 0)
}
