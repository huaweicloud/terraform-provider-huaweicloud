package geminidb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var geminiDBBackupStopNonUpdatableParams = []string{
	"backup_id",
}

// @API GeminiDB PUT /v3/{project_id}/backups/{backup_id}
// @API GeminiDB GET /v4/{project_id}/backups
// @API GeminiDB GET /v3/{project_id}/instances
func ResourceGeminiDBBackupStop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBBackupStopCreate,
		UpdateContext: resourceGeminiDBBackupStopUpdate,
		ReadContext:   resourceGeminiDBBackupStopRead,
		DeleteContext: resourceGeminiDBBackupStopDelete,

		CustomizeDiff: config.FlexibleForceNew(geminiDBBackupStopNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceGeminiDBBackupStopCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/backups/{backup_id}"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	backupID := d.Get("backup_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{backup_id}", backupID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         map[string]interface{}{"action": "stop"},
	}

	createResp, err := client.Request("PUT", createPath, &createOpt)

	if err != nil {
		return diag.Errorf("error stopping GeminiDB backup: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error stoping GeminiDB backup: job_id is not found in API response")
	}

	d.SetId(backupID)

	backup, err := GetBackup(client, backupID)
	if err != nil {
		return diag.FromErr(err)
	}

	if backup == nil {
		return diag.FromErr(golangsdk.ErrDefault404{})
	}

	instanceId := utils.PathSearch("instance_id", backup, "").(string)

	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Refresh:      geminiDbInstanceStatusRefreshFunc(client, instanceId),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        1 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for GeminiDB instance to ready: %s", err)
	}

	return nil
}

func resourceGeminiDBBackupStopRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBBackupStopUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBBackupStopDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GeminiDB backup stop is not supported. The backup stop resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
