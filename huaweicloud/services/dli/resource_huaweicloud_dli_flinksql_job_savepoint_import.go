package dli

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var flinkJobSavepointImportNonUpdatableParams = []string{
	"job_id",
	"savepoint_path",
}

// @API DLI POST /v1.0/{project_id}/streaming/jobs/{job_id}/import-savepoint
func ResourceFlinkSqlJobSavepointImport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlinkSqlJobSavepointImportCreate,
		ReadContext:   resourceFlinkSqlJobSavepointImportRead,
		UpdateContext: resourceFlinkSqlJobSavepointImportUpdate,
		DeleteContext: resourceFlinkSqlJobSavepointImportDelete,

		CustomizeDiff: config.FlexibleForceNew(flinkJobSavepointImportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the Flink SQL job savepoint is located.`,
			},

			// Required parameters.
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Flink SQL job.`,
			},
			"savepoint_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The OBS bucket path of the savepoint.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildFlinkJobSavepointImportBodyParams(d *schema.ResourceData, jobId string) map[string]interface{} {
	return map[string]interface{}{
		"job_id":         jobId,
		"savepoint_path": d.Get("savepoint_path"),
	}
}

func resourceFlinkSqlJobSavepointImportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	var (
		httpUrl = "v1.0/{project_id}/streaming/jobs/{job_id}/import-savepoint"
		jobId   = d.Get("job_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{job_id}", jobId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildFlinkJobSavepointImportBodyParams(d, jobId)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error importing Flink SQL job savepoint: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	isSuccess := utils.PathSearch("is_success", respBody, "false")
	if isSuccess.(string) == "false" {
		return diag.Errorf("error triggering Flink SQL job savepoint: %s", utils.PathSearch("message", respBody, ""))
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceFlinkSqlJobSavepointImportRead(ctx, d, meta)
}

func resourceFlinkSqlJobSavepointImportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFlinkSqlJobSavepointImportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFlinkSqlJobSavepointImportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for importing Flink SQL job savepoint. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
