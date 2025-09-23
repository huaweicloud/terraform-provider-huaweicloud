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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/components/{component_id}
func DataSourceSocComponentDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocComponentDetailRead,

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
			"component_id": {
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

func dataSourceSocComponentDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/components/{component_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{component_id}", d.Get("component_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster soc component detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataResp := utils.PathSearch("data", respBody, nil)

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenSocComponentDetailData(dataResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocComponentDetailData(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":                   utils.PathSearch("id", dataResp, nil),
			"name":                 utils.PathSearch("name", dataResp, nil),
			"dev_language":         utils.PathSearch("dev_language", dataResp, nil),
			"dev_language_version": utils.PathSearch("dev_language_version", dataResp, nil),
			"alliance_id":          utils.PathSearch("alliance_id", dataResp, nil),
			"alliance_name":        utils.PathSearch("alliance_name", dataResp, nil),
			"description":          utils.PathSearch("description", dataResp, nil),
			"logo":                 utils.PathSearch("logo", dataResp, nil),
			"label":                utils.PathSearch("label", dataResp, nil),
			"create_time":          utils.PathSearch("create_time", dataResp, nil),
			"update_time":          utils.PathSearch("update_time", dataResp, nil),
			"creator_name":         utils.PathSearch("creator_name", dataResp, nil),
			"operate_history": flattenSocComponentsOperateHistory(
				utils.PathSearch("operate_history", dataResp, make([]interface{}, 0)).([]interface{})),
			"component_versions": flattenSocComponentsComponentVersions(
				utils.PathSearch("component_versions", dataResp, make([]interface{}, 0)).([]interface{})),
			"component_type": utils.PathSearch("component_type", dataResp, nil),
		},
	}
}
