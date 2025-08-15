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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/nodes/installation-scripts
func DataSourceSecmasterInstallationScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterInstallationScriptsRead,

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
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commands": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSecmasterInstallationScriptsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/nodes/installation-scripts"
		workspaceID = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster installation scripts: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("records", flattenRecordsAttribute(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRecordsAttribute(respBody interface{}) []map[string]interface{} {
	records := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
	if len(records) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(records))
	for _, v := range records {
		rst = append(rst, map[string]interface{}{
			"os_type":  utils.PathSearch("os_type", v, nil),
			"commands": utils.PathSearch("commands", v, nil),
		})
	}

	return rst
}
