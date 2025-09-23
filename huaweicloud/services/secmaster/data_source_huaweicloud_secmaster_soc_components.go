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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/components
func DataSourceSocComponents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocComponentsRead,

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
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaSocComponentData(),
				},
			},
		},
	}
}

func schemaSocComponentData() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dev_language": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dev_language_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"alliance_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"alliance_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"logo": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"label": {
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
		"creator_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"operate_history": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"operate_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"operate_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"component_versions": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"version_num": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"version_desc": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"package_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"component_attachment_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"sub_version_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"connection_config": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"default_value": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"description": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"key": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"readonly": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"required": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"value": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"component_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceSocComponentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/components"
		offset  = 0
		result  = make([]interface{}, 0)
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
	}

	for {
		currentPath := fmt.Sprintf("%s?limit=200&offset=%v", requestPath, offset)
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster soc components: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
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
		d.Set("data", flattenSocComponentsData(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocComponentsData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"name":                 utils.PathSearch("name", v, nil),
			"dev_language":         utils.PathSearch("dev_language", v, nil),
			"dev_language_version": utils.PathSearch("dev_language_version", v, nil),
			"alliance_id":          utils.PathSearch("alliance_id", v, nil),
			"alliance_name":        utils.PathSearch("alliance_name", v, nil),
			"description":          utils.PathSearch("description", v, nil),
			"logo":                 utils.PathSearch("logo", v, nil),
			"label":                utils.PathSearch("label", v, nil),
			"create_time":          utils.PathSearch("create_time", v, nil),
			"update_time":          utils.PathSearch("update_time", v, nil),
			"creator_name":         utils.PathSearch("creator_name", v, nil),
			"operate_history": flattenSocComponentsOperateHistory(
				utils.PathSearch("operate_history", v, make([]interface{}, 0)).([]interface{})),
			"component_versions": flattenSocComponentsComponentVersions(
				utils.PathSearch("component_versions", v, make([]interface{}, 0)).([]interface{})),
			"component_type": utils.PathSearch("component_type", v, nil),
		})
	}

	return rst
}

func flattenSocComponentsOperateHistory(operateHistoryResp []interface{}) []interface{} {
	if len(operateHistoryResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(operateHistoryResp))
	for _, v := range operateHistoryResp {
		rst = append(rst, map[string]interface{}{
			"operator_name": utils.PathSearch("operator_name", v, nil),
			"operate_time":  utils.PathSearch("operate_time", v, nil),
		})
	}

	return rst
}

func flattenSocComponentsComponentVersions(componentVersionsResp []interface{}) []interface{} {
	if len(componentVersionsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(componentVersionsResp))
	for _, v := range componentVersionsResp {
		rst = append(rst, map[string]interface{}{
			"id":                      utils.PathSearch("id", v, nil),
			"version_num":             utils.PathSearch("version_num", v, nil),
			"version_desc":            utils.PathSearch("version_desc", v, nil),
			"status":                  utils.PathSearch("status", v, nil),
			"package_name":            utils.PathSearch("package_name", v, nil),
			"component_attachment_id": utils.PathSearch("component_attachment_id", v, nil),
			"sub_version_id":          utils.PathSearch("sub_version_id", v, nil),
			"connection_config": flattenConnectionConfig(
				utils.PathSearch("connection_config", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenConnectionConfig(connectionConfigResp []interface{}) []interface{} {
	if len(connectionConfigResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(connectionConfigResp))
	for _, v := range connectionConfigResp {
		rst = append(rst, map[string]interface{}{
			"default_value": utils.PathSearch("default_value", v, nil),
			"description":   utils.PathSearch("description", v, nil),
			"key":           utils.PathSearch("key", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"readonly":      utils.PathSearch("readonly", v, nil),
			"required":      utils.PathSearch("required", v, nil),
			"type":          utils.PathSearch("type", v, nil),
			"value":         utils.PathSearch("value", v, nil),
		})
	}

	return rst
}
