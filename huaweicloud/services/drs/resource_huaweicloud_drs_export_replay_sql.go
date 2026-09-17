package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var drsExportReplaySqlNonUpdatableParams = []string{
	"job_id",
	"file_type",
	"field_names",
	"start_time",
	"end_time",
}

// @API DRS POST /v5/{project_id}/jobs/{job_id}/export-replay-sql
func ResourceDrsExportReplaySql() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsExportReplaySqlCreate,
		ReadContext:   resourceDrsExportReplaySqlRead,
		UpdateContext: resourceDrsExportReplaySqlUpdate,
		DeleteContext: resourceDrsExportReplaySqlDelete,

		CustomizeDiff: config.FlexibleForceNew(drsExportReplaySqlNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the DRS replay job to export SQL report.",
			},
			"file_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"abnormal_sql", "abnormal_sql_detail", "slow_sql", "slow_sql_detail", "sql_result_inconsistent",
				}, false),
				Description: "Specifies the type of the SQL file to export.",
			},
			"field_names": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the list of field names to export. Required when `file_type` is **abnormal_sql_detail** or **slow_sql_detail**.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the start time of the export range.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the end time of the export range.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceDrsExportReplaySqlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/export-replay-sql"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{job_id}", d.Get("job_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDrsExportReplaySqlBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error exporting DRS replay SQL: %s", err)
	}

	_, err = utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The API returns an empty response body, use job_id as the resource ID
	d.SetId(d.Get("job_id").(string))

	return resourceDrsExportReplaySqlRead(ctx, d, meta)
}

func buildCreateDrsExportReplaySqlBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"file_type":   d.Get("file_type"),
		"field_names": utils.ValueIgnoreEmpty(d.Get("field_names").([]interface{})),
		"start_time":  utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_time":    utils.ValueIgnoreEmpty(d.Get("end_time")),
	}
	return bodyParams
}

func resourceDrsExportReplaySqlRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsExportReplaySqlUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsExportReplaySqlDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DRS export replay SQL resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
