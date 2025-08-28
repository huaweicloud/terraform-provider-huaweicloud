package coc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC GET /v1/application-model/next
func DataSourceCocApplicationNextModels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocApplicationNextModelsRead,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_applications": {
				Type:     schema.TypeList,
				Elem:     applicationNextModelSubApplicationsDataSchema(),
				Computed: true,
			},
			"components": {
				Type:     schema.TypeList,
				Elem:     applicationNextModelComponentsDataSchema(),
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Elem:     applicationNextModelGroupsDataSchema(),
				Computed: true,
			},
		},
	}
}

func applicationNextModelSubApplicationsDataSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
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
		},
	}
	return &sc
}

func applicationNextModelComponentsDataSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func applicationNextModelGroupsDataSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_tags": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"relation_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func dataSourceCocApplicationNextModelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/application-model/next"
		product = "coc"
		pageNo  = 0
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	basePath := client.Endpoint + httpUrl
	basePath += buildGetApplicationNextModelsParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var subApplicationsRes []map[string]interface{}
	componentsRes := make([]map[string]interface{}, 0)
	groupsRes := make([]map[string]interface{}, 0)
	for {
		getPath := basePath + fmt.Sprintf("&page_no=%d", pageNo)
		getResp, err := client.Request("GET", getPath, &getOpt)

		if err != nil {
			return diag.Errorf("error retrieving COC application next models: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		subApplicationsRes = flattenCocGetApplicationNextModelSubApplications(getRespBody)
		if len(subApplicationsRes) > 0 {
			break
		}

		components := flattenCocGetApplicationNextModelComponents(getRespBody)
		groups, err := flattenCocGetApplicationNextModelGroups(getRespBody)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(components) < 1 && len(groups) < 1 {
			break
		}
		componentsRes = append(componentsRes, components...)
		groupsRes = append(groupsRes, groups...)
		pageNo++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("sub_applications", subApplicationsRes),
		d.Set("components", componentsRes),
		d.Set("groups", groupsRes),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetApplicationNextModelsParams(d *schema.ResourceData) string {
	res := "?limit=100"
	if v, ok := d.GetOk("application_id"); ok {
		res = fmt.Sprintf("%s&application_id=%v", res, v)
	}
	if v, ok := d.GetOk("component_id"); ok {
		res = fmt.Sprintf("%s&component_id=%v", res, v)
	}

	return res
}

func flattenCocGetApplicationNextModelSubApplications(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	dataJson := utils.PathSearch("data.sub_applications", resp, make([]interface{}, 0))
	dataArray := dataJson.([]interface{})
	if len(dataArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataArray))
	for _, data := range dataArray {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", data, nil),
			"name":        utils.PathSearch("name", data, nil),
			"code":        utils.PathSearch("code", data, nil),
			"description": utils.PathSearch("description", data, nil),
			"domain_id":   utils.PathSearch("domain_id", data, nil),
			"parent_id":   utils.PathSearch("parent_id", data, nil),
			"path":        utils.PathSearch("path", data, nil),
			"create_time": utils.PathSearch("create_time", data, nil),
			"update_time": utils.PathSearch("update_time", data, nil),
		})
	}
	return result
}

func flattenCocGetApplicationNextModelComponents(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	dataJson := utils.PathSearch("data.components", resp, make([]interface{}, 0))
	dataArray := dataJson.([]interface{})
	if len(dataArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataArray))
	for _, data := range dataArray {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", data, nil),
			"name":           utils.PathSearch("name", data, nil),
			"code":           utils.PathSearch("code", data, nil),
			"domain_id":      utils.PathSearch("domain_id", data, nil),
			"application_id": utils.PathSearch("application_id", data, nil),
			"path":           utils.PathSearch("path", data, nil),
		})
	}
	return result
}

func flattenCocGetApplicationNextModelGroups(resp interface{}) ([]map[string]interface{}, error) {
	if resp == nil {
		return nil, nil
	}
	dataJson := utils.PathSearch("data.groups", resp, make([]interface{}, 0))
	dataArray := dataJson.([]interface{})
	if len(dataArray) == 0 {
		return nil, nil
	}

	result := make([]map[string]interface{}, 0, len(dataArray))
	for _, data := range dataArray {
		relationConfigurations := utils.PathSearch("relation_configurations", data,
			make([]interface{}, 0)).([]interface{})
		relationConfigurationList := make([]string, len(relationConfigurations))
		for j, relationConfiguration := range relationConfigurations {
			relationConfigurationJson, err := json.Marshal(relationConfiguration)
			if err != nil {
				return nil, err
			}
			relationConfigurationList[j] = string(relationConfigurationJson)
		}
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", data, nil),
			"name":           utils.PathSearch("name", data, nil),
			"code":           utils.PathSearch("code", data, nil),
			"domain_id":      utils.PathSearch("domain_id", data, nil),
			"region_id":      utils.PathSearch("region_id", data, nil),
			"application_id": utils.PathSearch("application_id", data, nil),
			"component_id":   utils.PathSearch("component_id", data, nil),
			"sync_mode":      utils.PathSearch("sync_mode", data, nil),
			"vendor":         utils.PathSearch("vendor", data, nil),
			"sync_rules": flattenCocGetApplicationNextModelGroupsSyncRules(
				utils.PathSearch("sync_rules", data, nil)),
			"relation_configurations": relationConfigurationList,
		})
	}
	return result, nil
}

func flattenCocGetApplicationNextModelGroupsSyncRules(rawParams interface{}) []map[string]interface{} {
	if rawParams == nil {
		return nil
	}
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) < 1 {
			return nil
		}
		res := make([]map[string]interface{}, len(paramsList))
		for i, params := range paramsList {
			raw := params.(map[string]interface{})
			res[i] = map[string]interface{}{
				"enterprise_project_id": utils.PathSearch("ep_id", raw, nil),
				"rule_tags":             utils.PathSearch("rule_tags", raw, nil),
			}
		}
		return res
	}

	return nil
}
