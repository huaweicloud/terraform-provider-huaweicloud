package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/cloud-logs/platform
func DataSourceCloudLogPlatforms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudLogPlatformsRead,

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
			"platforms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     cloudLogPlatformSchema(),
			},
		},
	}
}

func cloudLogPlatformSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tenant_managed_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform_managed_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dw_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"publish_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"white_list": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudLogPlatformsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/cloud-logs/platform"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster cloud log platforms: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("platforms", flattenCloudLogPlatforms(utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCloudLogPlatforms(result []interface{}) []interface{} {
	if len(result) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(result))
	for _, platform := range result {
		rst = append(rst, map[string]interface{}{
			"tenant_managed_domain_id":   utils.PathSearch("tenant_managed_domain_id", platform, nil),
			"platform_managed_domain_id": utils.PathSearch("platform_managed_domain_id", platform, nil),
			"dw_region":                  utils.PathSearch("dw_region", platform, nil),
			"create_time":                utils.PathSearch("create_time", platform, nil),
			"update_time":                utils.PathSearch("update_time", platform, nil),
			"publish_status":             utils.PathSearch("publish_status", platform, nil),
			"white_list":                 utils.PathSearch("whitelist", platform, nil),
		})
	}

	return rst
}
