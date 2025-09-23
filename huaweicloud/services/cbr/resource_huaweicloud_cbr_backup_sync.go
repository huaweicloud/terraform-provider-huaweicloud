package cbr

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableSyncParams = []string{
	"backup_id",
	"backup_name",
	"bucket_name",
	"image_path",
	"resource_id",
	"resource_name",
	"resource_type",
	"created_at",
}

// This resource uses the API for synchronizing local VMware backups.
// Due to the lack of test scenarios, this code is not tested and is not documented externally.
// Documentation is only stored in docs/incubating.

// @API CBR POST /v3/{project_id}/backups/sync
// @API CBR GET /v3/{project_id}/operation-logs/{operation_log_id}
func ResourceBackupSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupSyncCreate,
		ReadContext:   resourceBackupSyncRead,
		UpdateContext: resourceBackupSyncUpdate,
		DeleteContext: resourceBackupSyncDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableSyncParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level
region will be used.`,
			},
			"backup_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the backup ID to be synchronized.`,
			},
			"backup_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the backup.`,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the bucket where the backup is stored.`,
			},
			"image_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the path of the backup image in the bucket.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the resource to be backed up.`,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the resource to be backed up.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the type of the resource to be backed up.`,
			},
			"created_at": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the timestamp when the backup was created.`,
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

func buildSyncCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	syncMap := map[string]interface{}{
		"backup_id":     d.Get("backup_id").(string),
		"backup_name":   d.Get("backup_name").(string),
		"bucket_name":   d.Get("bucket_name").(string),
		"image_path":    d.Get("image_path").(string),
		"resource_id":   d.Get("resource_id").(string),
		"resource_name": d.Get("resource_name").(string),
		"resource_type": d.Get("resource_type").(string),
		"created_at":    d.Get("created_at").(int),
	}

	return map[string]interface{}{
		"sync": []interface{}{syncMap},
	}
}

func querySyncTask(client *golangsdk.ServiceClient, operationLogID string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{project_id}/operation-logs/{operation_log_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{operation_log_id}", operationLogID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBR backup sync task: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func waitingForSyncTaskSuccess(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, operationLogID string) error {
	unexpectedStatus := []string{"failed", "timeout"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := querySyncTask(client, operationLogID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("operation_log.status", respBody, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("status is not found in API response")
			}

			if status == "success" {
				return respBody, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, nil
			}
			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceBackupSyncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/backups/sync"
		product = "cbr"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildSyncCreateBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CBR backup sync: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	operationLogID := utils.PathSearch("sync[0].operation_log_id", respBody, "").(string)
	if operationLogID == "" {
		return diag.Errorf("error creating CBR backup sync: Operation Log ID is not found in API response")
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID for CBR backup sync: %s", err)
	}
	d.SetId(id)

	if err := waitingForSyncTaskSuccess(ctx, client, d.Timeout(schema.TimeoutCreate), operationLogID); err != nil {
		return diag.Errorf("error waiting for CBR backup sync to complete: %s", err)
	}

	return resourceBackupSyncRead(ctx, d, meta)
}

func resourceBackupSyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBackupSyncUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBackupSyncDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to synchronize a CBR backup. Deleting this 
resource will not change the current backup synchronization result, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
