package drs

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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/export/obs-uris
func DataSourceReportObsUris() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReportObsUrisRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"report_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"export_report_obs_files": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"obs_temp_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildReportObsUrisQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?report_type=%s", d.Get("report_type").(string))
}

func dataSourceReportObsUrisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/export/obs-uris"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", d.Get("job_id").(string))
	getPath += buildReportObsUrisQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS report obs uris: %s", err)
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
		d.Set("total_count", utils.PathSearch("count", respBody, nil)),
		d.Set("export_report_obs_files", flattenReportObsFiles(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReportObsFiles(respBody interface{}) []interface{} {
	reportObsFilesRaw := utils.PathSearch("export_report_obs_files", respBody, make([]interface{}, 0)).([]interface{})
	if len(reportObsFilesRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(reportObsFilesRaw))
	for _, item := range reportObsFilesRaw {
		result = append(result, map[string]interface{}{
			"file_name":    utils.PathSearch("file_name", item, nil),
			"obs_temp_uri": utils.PathSearch("obs_temp_uri", item, nil),
		})
	}
	return result
}
