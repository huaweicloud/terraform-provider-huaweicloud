package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/image-servers/{server_id}/actions/latest-attached-app
func DataSourceLatestAttachedApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLatestAttachedApplicationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region where the Workspace APP attached applications are located.`,
			},
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the image server instance.`,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the attached application.`,
						},
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The record ID of the attached application.`,
						},
					},
				},
				Description: `The list of latest attached applications.`,
			},
		},
	}
}

func dataSourceLatestAttachedApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/image-servers/{server_id}/actions/latest-attached-app"
		serverId = d.Get("server_id").(string)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_id}", serverId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving Workspace APP attached applications: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serverId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("applications", flattenAppLatestAttachedApplications(utils.PathSearch("items",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppLatestAttachedApplications(applications []interface{}) []interface{} {
	if len(applications) == 0 {
		return nil
	}

	result := make([]interface{}, len(applications))
	for i, v := range applications {
		result[i] = map[string]interface{}{
			"app_id":    utils.PathSearch("app_id", v, nil),
			"record_id": utils.PathSearch("id", v, nil),
		}
	}

	return result
}
