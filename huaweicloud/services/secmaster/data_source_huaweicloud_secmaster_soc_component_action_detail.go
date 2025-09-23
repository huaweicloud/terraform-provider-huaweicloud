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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/components/{component_id}/action/{action_id}
func DataSourceSocComponentActionDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocComponentActionDetailRead,

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
			"action_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaSocComponentActionData(),
				},
			},
		},
	}
}

func dataSourceSocComponentActionDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/components/{component_id}/action/{action_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{component_id}", d.Get("component_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{action_id}", d.Get("action_id").(string))
	requestPath += fmt.Sprintf("?enabled=%v", d.Get("enabled").(bool))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster soc component action detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
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
		d.Set("data", flattenSocComponentActionDetailData(dataResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocComponentActionDetailData(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":                    utils.PathSearch("id", dataResp, nil),
			"action_name":           utils.PathSearch("action_name", dataResp, nil),
			"action_desc":           utils.PathSearch("action_desc", dataResp, nil),
			"action_type":           utils.PathSearch("action_type", dataResp, nil),
			"create_time":           utils.PathSearch("create_time", dataResp, nil),
			"creator_name":          utils.PathSearch("creator_name", dataResp, nil),
			"can_update":            utils.PathSearch("can_update", dataResp, nil),
			"action_version_id":     utils.PathSearch("action_version_id", dataResp, nil),
			"action_version_name":   utils.PathSearch("action_version_name", dataResp, nil),
			"action_version_number": utils.PathSearch("action_version_number", dataResp, nil),
			"action_enable":         utils.PathSearch("action_enable", dataResp, nil),
		},
	}
}
