package dli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var sqlJobResultExportNonUpdatableParams = []string{
	"job_id",
	"data_path",
	"data_type",
	"compress",
	"queue_name",
	"export_mode",
	"with_column_header",
	"limit_num",
	"encoding_type",
	"quote_char",
	"escape_char",
}

// @API DLI POST /v1.0/{project_id}/jobs/{job_id}/export-result
// @API DLI GET /v1.0/{project_id}/jobs/{job_id}/status
func ResourceSqlJobResultExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSqlJobResultExportCreate,
		ReadContext:   resourceSqlJobResultExportRead,
		UpdateContext: resourceSqlJobResultExportUpdate,
		DeleteContext: resourceSqlJobResultExportDelete,

		CustomizeDiff: config.FlexibleForceNew(sqlJobResultExportNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the SQL job with result to be exported is located.`,
			},

			// Required parameters.
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the SQL job.`,
			},
			"data_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The OBS path for storing the exported data.`,
			},
			"data_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The storage format of the exported data.`,
			},

			// Optional parameters.
			"compress": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The compression format of the exported data.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The queue name used to execute the export task.`,
			},
			"export_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The mode of data export.`,
			},
			"with_column_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to export column names when exporting data.`,
			},
			"limit_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of data records to be exported.`,
			},
			"encoding_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The encoding format of the exported data.`,
			},
			"quote_char": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The custom quote character.`,
			},
			"escape_char": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The custom escape character.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func buildSqlJobResultExportBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"data_path":          d.Get("data_path"),
		"data_type":          d.Get("data_type"),
		"compress":           utils.ValueIgnoreEmpty(d.Get("compress")),
		"queue_name":         utils.ValueIgnoreEmpty(d.Get("queue_name")),
		"export_mode":        utils.ValueIgnoreEmpty(d.Get("export_mode")),
		"with_column_header": utils.ValueIgnoreEmpty(d.Get("with_column_header")),
		"limit_num":          utils.ValueIgnoreEmpty(d.Get("limit_num")),
		"encoding_type":      utils.ValueIgnoreEmpty(d.Get("encoding_type")),
		"quote_char":         utils.ValueIgnoreEmpty(d.Get("quote_char")),
		"escape_char":        utils.ValueIgnoreEmpty(d.Get("escape_char")),
	}
}

func exportSqlJobResult(client *golangsdk.ServiceClient, jobId string, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v1.0/{project_id}/jobs/{job_id}/export-result"
	exportPath := client.Endpoint + httpUrl
	exportPath = strings.ReplaceAll(exportPath, "{project_id}", client.ProjectID)
	exportPath = strings.ReplaceAll(exportPath, "{job_id}", jobId)

	exportOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildSqlJobResultExportBodyParams(d)),
	}

	resp, err := client.Request("POST", exportPath, &exportOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceSqlJobResultExportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		jobId  = d.Get("job_id").(string)
	)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	respBody, err := exportSqlJobResult(client, jobId, d)
	if err != nil {
		return diag.Errorf("error exporting query result of the SQL job (%s) to OBS: %s", jobId, err)
	}

	if !utils.PathSearch("is_success", respBody, false).(bool) {
		return diag.Errorf("unable to export the query result of the SQL job (%s) to OBS: %s", jobId,
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	exportJobId := utils.PathSearch("job_id", respBody, "").(string)
	if exportJobId == "" {
		return diag.Errorf("unable to find the export job ID in the API response")
	}

	if err = waitForSqlJobResultExportCompleted(ctx, client, exportJobId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for SQL job result to be exported: %s", err)
	}

	d.SetId(exportJobId)

	return resourceSqlJobResultExportRead(ctx, d, meta)
}

func getSqlJobStatus(client *golangsdk.ServiceClient, jobId string) (interface{}, error) {
	httpUrl := "v1.0/{project_id}/jobs/{job_id}/status"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitForSqlJobResultExportCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getSqlJobStatus(client, jobId)
			if err != nil {
				return nil, "ERROR", err
			}

			if !utils.PathSearch("is_success", respBody, false).(bool) {
				return respBody, "ERROR", errors.New(utils.PathSearch("message", respBody, "Message Not Found").(string))
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "FAILED" {
				return respBody, "ERROR", fmt.Errorf("unexpected status (%s), error message: %s", status,
					utils.PathSearch("message", respBody, "Message Not Found").(string))
			}

			if status == "FINISHED" {
				return respBody, "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceSqlJobResultExportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSqlJobResultExportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSqlJobResultExportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for exporting query result of the SQL job to OBS. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
