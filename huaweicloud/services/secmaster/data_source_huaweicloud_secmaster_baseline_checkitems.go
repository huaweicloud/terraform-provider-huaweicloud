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

// @API SecMaster POST /v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems/search
func DataSourceBaselineCheckitems() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineCheckitemsRead,

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
			"catalog_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compliance_package_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"suggestion": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"source_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"condition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"data": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"logics": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"query_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"builtin_checkitem_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"checkitem_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"customized_checkitem_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"checkitems": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaBaselineCheckitemsData(),
				},
			},
		},
	}
}

func schemaBaselineCheckitemsData() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aggregation_handle_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"audit_procedure": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"impact": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cloud_server": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"level": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"method": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"source": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"workflow_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"spec_checkitem_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Named as `checkitemUuid` in the API documentation.
					"checkitem_uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					// Named as `createTime` in the API documentation.
					"create_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"language": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					// Named as `removeTime` in the API documentation.
					"remove_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					// Named as `specificationUuid` in the API documentation.
					"specification_uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func buildListCheckitemsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"catalog_uuid":          utils.ValueIgnoreEmpty(d.Get("catalog_uuid")),
		"compliance_package_id": utils.ValueIgnoreEmpty(d.Get("compliance_package_id")),
		"sort_by":               utils.ValueIgnoreEmpty(d.Get("sort_by")),
		"order":                 utils.ValueIgnoreEmpty(d.Get("order")),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"suggestion":            utils.ValueIgnoreEmpty(d.Get("suggestion")),
		"type":                  utils.ValueIgnoreEmpty(d.Get("type")),
		"source_list":           utils.ValueIgnoreEmpty(d.Get("source_list")),
		"query_mode":            utils.ValueIgnoreEmpty(d.Get("query_mode")),
		"severity":              utils.ValueIgnoreEmpty(d.Get("severity")),
	}

	if v, ok := d.GetOk("condition"); ok {
		conditionList := v.([]interface{})
		if len(conditionList) > 0 && conditionList[0] != nil {
			conditionMap := conditionList[0].(map[string]interface{})
			transformedCondition := make(map[string]interface{})

			if logicsRaw, ok := conditionMap["logics"]; ok {
				transformedCondition["logics"] = expandToStringList(logicsRaw.([]interface{}))
			}

			if conditionsRaw, ok := conditionMap["conditions"]; ok {
				var conditionsResult []map[string]interface{}
				for _, item := range conditionsRaw.([]interface{}) {
					iMap := item.(map[string]interface{})
					newMap := make(map[string]interface{})

					if name, ok := iMap["name"]; ok {
						newMap["name"] = name.(string)
					}

					if dataRaw, ok := iMap["data"]; ok {
						newMap["data"] = expandToStringList(dataRaw.([]interface{}))
					}
					conditionsResult = append(conditionsResult, newMap)
				}
				transformedCondition["conditions"] = conditionsResult
			}

			bodyParams["condition"] = transformedCondition
		}
	}

	return bodyParams
}

func dataSourceBaselineCheckitemsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                    = meta.(*config.Config)
		region                 = cfg.GetRegion(d)
		product                = "secmaster"
		httpUrl                = "v2/{project_id}/workspaces/{workspace_id}/sa/baseline/checkitems/search"
		offset                 = 0
		result                 = make([]interface{}, 0)
		builtinCheckitemNum    float64
		checkitemNum           float64
		customizedCheckitemNum float64
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
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
	}

	currentBodyParams := buildListCheckitemsBodyParams(d)
	currentBodyParams["limit"] = 1000
	for {
		currentBodyParams["offset"] = offset
		requestOpt.JSONBody = utils.RemoveNil(currentBodyParams)

		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster baseline checkitems: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		builtinCheckitemNum = utils.PathSearch("builtin_checkitem_num", respBody, float64(0)).(float64)
		checkitemNum = utils.PathSearch("checkitem_num", respBody, float64(0)).(float64)
		customizedCheckitemNum = utils.PathSearch("customized_checkitem_num", respBody, float64(0)).(float64)

		checkitemsResp := utils.PathSearch("checkitems", respBody, make([]interface{}, 0)).([]interface{})
		if len(checkitemsResp) == 0 {
			break
		}

		result = append(result, checkitemsResp...)
		offset += len(checkitemsResp)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("builtin_checkitem_num", builtinCheckitemNum),
		d.Set("checkitem_num", checkitemNum),
		d.Set("customized_checkitem_num", customizedCheckitemNum),
		d.Set("checkitems", flattenBaselineCheckitemsData(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBaselineCheckitemsData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"aggregation_handle_status": utils.PathSearch("aggregation_handle_status", v, nil),
			"audit_procedure":           utils.PathSearch("audit_procedure", v, nil),
			"impact":                    utils.PathSearch("impact", v, nil),
			"cloud_server":              utils.PathSearch("cloud_server", v, nil),
			"description":               utils.PathSearch("description", v, nil),
			"level":                     utils.PathSearch("level", v, nil),
			"method":                    utils.PathSearch("method", v, nil),
			"name":                      utils.PathSearch("name", v, nil),
			"source":                    utils.PathSearch("source", v, nil),
			"workflow_id":               utils.PathSearch("workflow_id", v, nil),
			"spec_checkitem_list": flattenSpecCheckitemList(
				utils.PathSearch("spec_checkitem_list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenSpecCheckitemList(specCheckitemList []interface{}) []interface{} {
	if len(specCheckitemList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(specCheckitemList))
	for _, v := range specCheckitemList {
		rst = append(rst, map[string]interface{}{
			"checkitem_uuid":     utils.PathSearch("checkitemUuid", v, nil),
			"create_time":        utils.PathSearch("createTime", v, nil),
			"language":           utils.PathSearch("language", v, nil),
			"name":               utils.PathSearch("name", v, nil),
			"remove_time":        utils.PathSearch("removeTime", v, nil),
			"specification_uuid": utils.PathSearch("specificationUuid", v, nil),
			"uuid":               utils.PathSearch("uuid", v, nil),
		})
	}

	return rst
}
