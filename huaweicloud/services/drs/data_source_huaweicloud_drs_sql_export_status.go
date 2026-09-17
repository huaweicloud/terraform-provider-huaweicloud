package drs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v5/{project_id}/jobs/{job_id}/sql-export-status
func DataSourceDrsSqlExportStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsSqlExportStatusRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DRS replay job to query SQL export status.",
			},
			"file_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"abnormal_sql", "abnormal_sql_detail", "slow_sql", "slow_sql_detail",
				}, false),
				Description: "Specifies the type of the SQL file to query export status.",
			},
			"export_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The export status of the SQL file.",
			},
			"failed_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The reason of the export failure.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total amount of exported data.",
			},
			"current_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of currently processed data.",
			},
			"progress_percentage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The progress percentage of the export task.",
			},
			"uploaded_file_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The names of the files uploaded to OBS.",
			},
		},
	}
}

func dataSourceDrsSqlExportStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/sql-export-status"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", d.Get("job_id").(string))
	getPath += buildListSqlExportStatusQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS SQL export status: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("export_status", utils.PathSearch("export_status", getRespBody, nil)),
		d.Set("failed_reason", utils.PathSearch("failed_reason", getRespBody, nil)),
		d.Set("total_count", utils.PathSearch("total_count", getRespBody, nil)),
		d.Set("current_count", utils.PathSearch("current_count", getRespBody, nil)),
		d.Set("progress_percentage", utils.PathSearch("progress_percentage", getRespBody, nil)),
		d.Set("uploaded_file_names", utils.PathSearch("uploaded_file_names", getRespBody, make([]interface{}, 0))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListSqlExportStatusQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?file_type=%v", d.Get("file_type"))

	return res
}
