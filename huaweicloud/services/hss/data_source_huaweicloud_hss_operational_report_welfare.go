package hss

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

// @API HSS GET /v5/{project_id}/operational-report/welfare
func DataSourceOperationalReportWelfare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOperationalReportWelfareRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hot_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"version_update_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"activities_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOperationalReportWelfareRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/operational-report/welfare"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the information in the news and promotions area of a monthly operations report: %s", err)
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
		d.Set("hot_info", flattenWelfareHoTInfo(utils.PathSearch("hot_info", getRespBody, nil))),
		d.Set("version_update_info", flattenWelfareVersionInfo(utils.PathSearch("version_update_info", getRespBody, nil))),
		d.Set("activities_info", flattenWelfareActivityInfo(utils.PathSearch("activities_info", getRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWelfareHoTInfo(hotInfo interface{}) []map[string]interface{} {
	if hotInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"title":    utils.PathSearch("title", hotInfo, nil),
		"url_json": utils.PathSearch("url_json", hotInfo, nil),
	}

	return []map[string]interface{}{result}
}

func flattenWelfareVersionInfo(versionInfo interface{}) []map[string]interface{} {
	if versionInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"title":    utils.PathSearch("title", versionInfo, nil),
		"url_json": utils.PathSearch("url_json", versionInfo, nil),
	}

	return []map[string]interface{}{result}
}

func flattenWelfareActivityInfo(activityInfo interface{}) []map[string]interface{} {
	if activityInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"title":    utils.PathSearch("title", activityInfo, nil),
		"url_json": utils.PathSearch("url_json", activityInfo, nil),
	}

	return []map[string]interface{}{result}
}
