package workspace

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var appServerBatchMigrateNonUpdatableParams = []string{"server_ids", "host_id"}

// @API Workspace PATCH /v1/{project_id}/app-servers/hosts/batch-migrate
// @API Workspace GET /v2/{project_id}/job/{job_id}
func ResourceAppServerBatchMigrate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppServerBatchMigrateCreate,
		ReadContext:   resourceAppServerBatchMigrateRead,
		UpdateContext: resourceAppServerBatchMigrateUpdate,
		DeleteContext: resourceAppServerBatchMigrateDelete,

		CustomizeDiff: config.FlexibleForceNew(appServerBatchMigrateNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the servers to be migrated are located.`,
			},
			"server_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of server IDs to be migrated.`,
			},
			"host_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the target cloud office host.`,
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

func buildAppServerBatchMigrateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"server_ids": d.Get("server_ids"),
		"host_id":    d.Get("host_id"),
	}
}

func resourceAppServerBatchMigrateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-servers/hosts/batch-migrate"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: buildAppServerBatchMigrateBodyParams(d),
	}

	resp, err := client.Request("PATCH", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error migrating servers to target cloud office host: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find job ID from API response")
	}

	_, err = waitForAppServerJobCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), jobId)
	if err != nil {
		return diag.Errorf("error waiting for the migration job (%s) completed: %s", jobId, err)
	}

	return nil
}

func resourceAppServerBatchMigrateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppServerBatchMigrateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppServerBatchMigrateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch migrate servers to target cloud office host. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
