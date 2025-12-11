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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/modules
func DataSourceModules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocModulesRead,

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
			// This parameter does not take effect
			"module_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"sort_dir": {
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
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
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
						"module_json": {
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
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"thumbnail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"module_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_built_in": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"data_query": {
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

func buildSocModulesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=100"

	if v, ok := d.GetOk("module_type"); ok {
		queryParams = fmt.Sprintf("%s&module_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceSocModulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/modules"
		result      = make([]interface{}, 0)
		offset      = 0
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath += buildSocModulesQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving modules: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		moduleData := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(moduleData) == 0 {
			break
		}

		result = append(result, moduleData...)
		offset += len(moduleData)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenSocModules(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocModules(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"cloud_pack_id":      utils.PathSearch("cloud_pack_id", v, nil),
			"cloud_pack_name":    utils.PathSearch("cloud_pack_name", v, nil),
			"cloud_pack_version": utils.PathSearch("cloud_pack_version", v, nil),
			"create_time":        utils.PathSearch("create_time", v, nil),
			"creator_id":         utils.PathSearch("creator_id", v, nil),
			"description":        utils.PathSearch("description", v, nil),
			"en_description":     utils.PathSearch("en_description", v, nil),
			"id":                 utils.PathSearch("id", v, nil),
			"module_json":        utils.PathSearch("module_json", v, nil),
			"name":               utils.PathSearch("name", v, nil),
			"en_name":            utils.PathSearch("en_name", v, nil),
			"project_id":         utils.PathSearch("project_id", v, nil),
			"workspace_id":       utils.PathSearch("workspace_id", v, nil),
			"update_time":        utils.PathSearch("update_time", v, nil),
			"thumbnail":          utils.PathSearch("thumbnail", v, nil),
			"module_type":        utils.PathSearch("module_type", v, nil),
			"tag":                utils.PathSearch("tag", v, nil),
			"is_built_in":        utils.PathSearch("is_built_in", v, nil),
			"data_query":         utils.PathSearch("data_query", v, nil),
			"boa_version":        utils.PathSearch("boa_version", v, nil),
			"version":            utils.PathSearch("version", v, nil),
		})
	}

	return rst
}
