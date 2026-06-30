package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/components/template
func DataSourceComponentAllTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComponentAllTemplatesRead,

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
			"search_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_config": {
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
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildComponentAllTemplatesQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("search_value"); ok {
		return fmt.Sprintf("?search_value=%s", v.(string))
	}

	return ""
}

func dataSourceComponentAllTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/components/template"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath += buildComponentAllTemplatesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// The pagination functionality of this API has issues, so it is not currently supported.
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving component all templates: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	dataList := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenComponentAllTemplatesData(dataList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenComponentAllTemplatesData(allResult []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(allResult))
	for _, v := range allResult {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"component_id":  utils.PathSearch("component_id", v, nil),
			"template_name": utils.PathSearch("template_name", v, nil),
			"task_config":   utils.PathSearch("task_config", v, nil),
			"create_time":   utils.PathSearch("create_time", v, nil),
			"update_time":   utils.PathSearch("update_time", v, nil),
			"project_id":    utils.PathSearch("project_id", v, nil),
		})
	}
	return rst
}
