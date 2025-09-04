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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/components
func DataSourceComponents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComponentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workspace ID.",
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The components list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The component ID.",
						},
						"component_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The component name.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The SecMaster version.",
						},
						"history_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The history version.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The creation time (timestamp in milliseconds).",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update time (timestamp in milliseconds).",
						},
						"time_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time zone.",
						},
						"upgrade": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The upgrade status.",
						},
						"upgrade_fail_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The upgrade failure message.",
						},
						"maintainer": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The maintainer.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The component description.",
						},
					},
				},
			},
		},
	}
}

func dataSourceComponentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/components"
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

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster components: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	recordsResp := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenComponentsRecords(recordsResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenComponentsRecords(recordsResp []interface{}) []interface{} {
	if len(recordsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(recordsResp))
	for _, v := range recordsResp {
		rst = append(rst, map[string]interface{}{
			"component_id":         utils.PathSearch("component_id", v, nil),
			"component_name":       utils.PathSearch("component_name", v, nil),
			"version":              utils.PathSearch("version", v, nil),
			"history_version":      utils.PathSearch("history_version", v, nil),
			"create_time":          utils.PathSearch("create_time", v, nil),
			"update_time":          utils.PathSearch("update_time", v, nil),
			"time_zone":            utils.PathSearch("time_zone", v, nil),
			"upgrade":              utils.PathSearch("upgrade", v, nil),
			"upgrade_fail_message": utils.PathSearch("upgrade_fail_message", v, nil),
			"maintainer":           utils.PathSearch("maintainer", v, nil),
			"description":          utils.PathSearch("description", v, nil),
		})
	}

	return rst
}
