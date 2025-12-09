package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Secmaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/layouts/search
func DataSourceLayouts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLayoutsRead,

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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"used_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"binding_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_built_in": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_template": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"layout_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search_txt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_pack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_pack_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_pack_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_built_in": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_template": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_name": {
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"en_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layout_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"thumbnail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layout_cfg": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layout_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"binding_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"binding_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"binding_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fields_sum": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"wizards_sum": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sections_sum": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"modules_sum": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tabs_sum": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boa_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildReadLayoutsBodyParams(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		"name":         utils.ValueIgnoreEmpty(d.Get("name")),
		"used_by":      utils.ValueIgnoreEmpty(d.Get("used_by")),
		"binding_code": utils.ValueIgnoreEmpty(d.Get("binding_code")),
		"layout_type":  utils.ValueIgnoreEmpty(d.Get("layout_type")),
		"sort_key":     utils.ValueIgnoreEmpty(d.Get("sort_key")),
		"sort_dir":     utils.ValueIgnoreEmpty(d.Get("sort_dir")),
		"search_txt":   utils.ValueIgnoreEmpty(d.Get("search_txt")),
		"from_date":    utils.ValueIgnoreEmpty(d.Get("from_date")),
		"to_date":      utils.ValueIgnoreEmpty(d.Get("to_date")),
	}

	if d.Get("is_built_in").(bool) {
		rst["is_built_in"] = true
	}

	if d.Get("is_template").(bool) {
		rst["is_template"] = true
	}

	if d.Get("is_default").(bool) {
		rst["is_default"] = true
	}

	return rst
}

func dataSourceLayoutsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/search"
		offset      = 0
		result      = make([]interface{}, 0)
		requestBody = utils.RemoveNil(buildReadLayoutsBodyParams(d))
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		// The API requires an offset to be passed in; otherwise, the API will throw an error.
		requestBody["offset"] = offset
		requestOpt.JSONBody = requestBody
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster layouts: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenLayoutsData(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLayoutsData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"cloud_pack_id":      utils.PathSearch("cloud_pack_id", v, nil),
			"cloud_pack_name":    utils.PathSearch("cloud_pack_name", v, nil),
			"cloud_pack_version": utils.PathSearch("cloud_pack_version", v, nil),
			"is_built_in":        utils.PathSearch("is_built_in", v, nil),
			"is_template":        utils.PathSearch("is_template", v, nil),
			"create_time":        utils.PathSearch("create_time", v, nil),
			"creator_id":         utils.PathSearch("creator_id", v, nil),
			"parent_id":          utils.PathSearch("parent_id", v, nil),
			"creator_name":       utils.PathSearch("creator_name", v, nil),
			"description":        utils.PathSearch("description", v, nil),
			"en_description":     utils.PathSearch("en_description", v, nil),
			"id":                 utils.PathSearch("id", v, nil),
			"name":               utils.PathSearch("name", v, nil),
			"en_name":            utils.PathSearch("en_name", v, nil),
			"layout_json":        utils.PathSearch("layout_json", v, nil),
			"project_id":         utils.PathSearch("project_id", v, nil),
			"update_time":        utils.PathSearch("update_time", v, nil),
			"workspace_id":       utils.PathSearch("workspace_id", v, nil),
			"region_id":          utils.PathSearch("region_id", v, nil),
			"domain_id":          utils.PathSearch("domain_id", v, nil),
			"thumbnail":          utils.PathSearch("thumbnail", v, nil),
			"used_by":            utils.PathSearch("used_by", v, nil),
			"layout_cfg":         utils.PathSearch("layout_cfg", v, nil),
			"layout_type":        utils.PathSearch("layout_type", v, nil),
			"binding_id":         utils.PathSearch("binding_id", v, nil),
			"binding_name":       utils.PathSearch("binding_name", v, nil),
			"binding_code":       utils.PathSearch("binding_code", v, nil),
			"fields_sum":         utils.PathSearch("fields_sum", v, nil),
			"wizards_sum":        utils.PathSearch("wizards_sum", v, nil),
			"sections_sum":       utils.PathSearch("sections_sum", v, nil),
			"modules_sum":        utils.PathSearch("modules_sum", v, nil),
			"tabs_sum":           utils.PathSearch("tabs_sum", v, nil),
			"version":            utils.PathSearch("version", v, nil),
			"boa_version":        utils.PathSearch("boa_version", v, nil),
		})
	}

	return rst
}
