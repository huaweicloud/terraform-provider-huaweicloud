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

// @API Secmaster GET /v2/{project_id}/workspaces/{workspace_id}/siem/directories
func DataSourceSiemDirectories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSiemDirectoriesRead,

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
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspaceid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"directory_i18ns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"directory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_fr": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSiemDirectoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/workspaces/{workspace_id}/siem/directories"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", d.Get("workspace_id").(string))
	listPath = fmt.Sprintf("%s?category=%s", listPath, d.Get("category").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving directories: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspaceid", utils.PathSearch("workspace_id", listRespBody, nil)),
		d.Set("project_id", utils.PathSearch("project_id", listRespBody, nil)),
		d.Set("directories", utils.PathSearch("directories", listRespBody, nil)),
		d.Set("directory_i18ns", flattenSiemDirectories(
			utils.PathSearch("directory_i18ns", listRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSiemDirectories(simeDirectories []interface{}) []interface{} {
	if len(simeDirectories) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(simeDirectories))
	for _, v := range simeDirectories {
		rst = append(rst, map[string]interface{}{
			"directory":    utils.PathSearch("directory", v, nil),
			"directory_en": utils.PathSearch("directory_en", v, nil),
			"directory_fr": utils.PathSearch("directory_fr", v, nil),
		})
	}

	return rst
}
