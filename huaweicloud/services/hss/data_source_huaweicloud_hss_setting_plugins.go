package hss

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

// @API HSS GET /v5/{project_id}/setting/plugins
func DataSourceSettingPlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSettingPluginsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"detect_result": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_addr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"protect_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"policy_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"policy_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"refresh": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// This parameter does not take effect
			"above_version": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// This parameter does not take effect
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"outside_host": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upgrade_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildSettingPluginsQueryParamsA(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?name=%v&limit=200", d.Get("name"))

	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		queryParams = fmt.Sprintf("%s&public_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("agent_status"); ok {
		queryParams = fmt.Sprintf("%s&agent_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_status"); ok {
		queryParams = fmt.Sprintf("%s&host_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("os_type"); ok {
		queryParams = fmt.Sprintf("%s&os_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("ip_addr"); ok {
		queryParams = fmt.Sprintf("%s&ip_addr=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func buildSettingPluginsQueryParamsB(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}
	if v, ok := d.GetOk("detect_result"); ok {
		queryParams = fmt.Sprintf("%s&detect_result=%v", queryParams, v)
	}
	if v, ok := d.GetOk("protect_status"); ok {
		queryParams = fmt.Sprintf("%s&protect_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("group_name"); ok {
		queryParams = fmt.Sprintf("%s&group_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("policy_group_id"); ok {
		queryParams = fmt.Sprintf("%s&policy_group_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("policy_group_name"); ok {
		queryParams = fmt.Sprintf("%s&policy_group_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("label"); ok {
		queryParams = fmt.Sprintf("%s&label=%v", queryParams, v)
	}
	if v, ok := d.GetOk("charging_mode"); ok {
		queryParams = fmt.Sprintf("%s&charging_mode=%v", queryParams, v)
	}
	if d.Get("refresh").(bool) {
		queryParams = fmt.Sprintf("%s&refresh=%v", queryParams, d.Get("refresh"))
	}
	if d.Get("above_version").(bool) {
		queryParams = fmt.Sprintf("%s&above_version=%v", queryParams, d.Get("above_version"))
	}
	if v, ok := d.GetOk("version"); ok {
		queryParams = fmt.Sprintf("%s&version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("plugin"); ok {
		queryParams = fmt.Sprintf("%s&plugin=%v", queryParams, v)
	}
	if d.Get("outside_host").(bool) {
		queryParams = fmt.Sprintf("%s&outside_host=%v", queryParams, d.Get("outside_host"))
	}

	return queryParams
}

func dataSourceSettingPluginsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/setting/plugins"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildSettingPluginsQueryParamsA(d, epsId)
	getPath += buildSettingPluginsQueryParamsB(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving plug-ins: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenSettingPlugins(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSettingPlugins(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"host_name":      utils.PathSearch("host_name", v, nil),
			"host_id":        utils.PathSearch("host_id", v, nil),
			"private_ip":     utils.PathSearch("private_ip", v, nil),
			"public_ip":      utils.PathSearch("public_ip", v, nil),
			"os_type":        utils.PathSearch("os_type", v, nil),
			"plugin_name":    utils.PathSearch("plugin_name", v, nil),
			"plugin_version": utils.PathSearch("plugin_version", v, nil),
			"plugin_status":  utils.PathSearch("plugin_status", v, nil),
			"upgrade_status": utils.PathSearch("upgrade_status", v, nil),
		})
	}

	return result
}
