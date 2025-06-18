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

// @API HSS GET /v5/{project_id}/setting/agent-install-script
func DataSourceAgentInstallScript() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAgentInstallScriptRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"os_arch": {
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
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"install_script_list": {
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

func buildAgentInstallScriptQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?os_type=%v", d.Get("os_type"))
	queryParams = fmt.Sprintf("%s&os_arch=%v", queryParams, d.Get("os_arch"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if d.Get("outside_host").(bool) {
		queryParams = fmt.Sprintf("%s&outside_host=%v", queryParams, true)
	}
	if d.Get("batch_install").(bool) {
		queryParams = fmt.Sprintf("%s&batch_install=%v", queryParams, true)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceAgentInstallScriptRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/setting/agent-install-script"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildAgentInstallScriptQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS agent install script: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	scriptResp := utils.PathSearch("install_script_list", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("install_script_list", flattenAgentInstallScript(scriptResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAgentInstallScript(scriptResp []interface{}) []interface{} {
	result := make([]interface{}, 0, len(scriptResp))
	for _, v := range scriptResp {
		result = append(result, map[string]interface{}{
			"package_type":         utils.PathSearch("package_type", v, nil),
			"cmd":                  utils.PathSearch("cmd", v, nil),
			"package_download_url": utils.PathSearch("package_download_url", v, nil),
		})
	}

	return result
}
