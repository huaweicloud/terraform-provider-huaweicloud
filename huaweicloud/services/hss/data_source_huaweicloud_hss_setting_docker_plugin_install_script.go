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

// @API HSS GET /v5/{project_id}/setting/docker-plugin-install-script
func DataSourceSettingDockerPluginInstallScript() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSettingDockerPluginInstallScriptRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"plugin": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operate_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"outside_host": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"batch_install": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
						"package_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cmd": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_download_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildSettingDockerPluginInstallScriptQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?plugin=%v&operate_type=%v", d.Get("plugin"), d.Get("operate_type"))

	if d.Get("outside_host").(bool) {
		queryParams = fmt.Sprintf("%s&outside_host=%v", queryParams, d.Get("outside_host"))
	}
	if d.Get("batch_install").(bool) {
		queryParams = fmt.Sprintf("%s&batch_install=%v", queryParams, d.Get("batch_install"))
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceSettingDockerPluginInstallScriptRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/setting/docker-plugin-install-script"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildSettingDockerPluginInstallScriptQueryParams(d, epsId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving docker plugin installation script: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataList := utils.PathSearch("install_script_list", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenSettingDockerPluginInstallScript(dataList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSettingDockerPluginInstallScript(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"package_type":         utils.PathSearch("package_type", v, nil),
			"cmd":                  utils.PathSearch("cmd", v, nil),
			"package_download_url": utils.PathSearch("package_download_url", v, nil),
		})
	}

	return rst
}
