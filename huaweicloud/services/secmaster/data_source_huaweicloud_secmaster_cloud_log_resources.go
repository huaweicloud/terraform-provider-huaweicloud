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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/cloud-logs/resource
func DataSourceCloudLogResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocCloudLogResourcesRead,

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
			// This parameter does not take effect
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datasets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"allow_alert": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"allow_lts": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"success": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"exist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"workspaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildSocCloudLogResourcesQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("region_id"); ok {
		queryParams = fmt.Sprintf("%s&region_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceSocCloudLogResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/cloud-logs/resource"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath += buildSocCloudLogResourcesQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving cloud log resources: %s", err)
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
		d.Set("datasets", flattenSocCloudLogResources(
			utils.PathSearch("datasets", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("exist", utils.PathSearch("exist", getRespBody, nil)),
		d.Set("workspaces", utils.PathSearch("workspaces", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocCloudLogResources(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"alert":        utils.PathSearch("alert", v, nil),
			"allow_alert":  utils.PathSearch("allow_alert", v, nil),
			"allow_lts":    utils.PathSearch("allow_lts", v, nil),
			"create_time":  utils.PathSearch("create_time", v, nil),
			"domain_id":    utils.PathSearch("domain_id", v, nil),
			"enable":       utils.PathSearch("enable", v, nil),
			"project_id":   utils.PathSearch("project_id", v, nil),
			"region":       utils.PathSearch("region", v, nil),
			"region_id":    utils.PathSearch("region_id", v, nil),
			"success":      utils.PathSearch("success", v, nil),
			"total":        utils.PathSearch("total", v, nil),
			"update_time":  utils.PathSearch("update_time", v, nil),
			"workspace_id": utils.PathSearch("workspace_id", v, nil),
		})
	}

	return rst
}
