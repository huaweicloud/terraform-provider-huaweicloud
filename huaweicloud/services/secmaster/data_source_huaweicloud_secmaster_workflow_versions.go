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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}/versions
func DataSourceWorkflowVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowVersionsRead,

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
			"workflow_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"taskconfig": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"taskflow": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"taskflow_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildWorkflowVersionsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s?status=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceWorkflowVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		workflowId  = d.Get("workflow_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}/versions"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{workflow_id}", workflowId)
	getPath += buildWorkflowVersionsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving workflow versions: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	versions := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenWorkflowVersions(versions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWorkflowVersions(versionsResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(versionsResp))
	for _, v := range versionsResp {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"description":   utils.PathSearch("description", v, nil),
			"workflow_id":   utils.PathSearch("aopworkflow_id", v, nil),
			"project_id":    utils.PathSearch("project_id", v, nil),
			"owner_id":      utils.PathSearch("owner_id", v, nil),
			"creator_id":    utils.PathSearch("creator_id", v, nil),
			"enabled":       utils.PathSearch("enabled", v, nil),
			"status":        utils.PathSearch("status", v, nil),
			"version":       utils.PathSearch("version", v, nil),
			"taskconfig":    utils.PathSearch("taskconfig", v, nil),
			"taskflow":      utils.PathSearch("taskflow", v, nil),
			"taskflow_type": utils.PathSearch("taskflow_type", v, nil),
			"aop_type":      utils.PathSearch("aop_type", v, nil),
		})
	}

	return rst
}
