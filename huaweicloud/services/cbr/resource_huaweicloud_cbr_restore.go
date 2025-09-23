package cbr

import (
	"context"
	"errors"
	"fmt"
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

var nonUpdatableRestoreParams = []string{
	"backup_id",
	"mappings",
	"mappings.*.backup_id",
	"mappings.*.volume_id",
	"power_on",
	"server_id",
	"volume_id",
	"resource_id",
	"details",
	"details.*.destination_path",
}

// @API CBR POST /v3/{project_id}/backups/{backup_id}/restore
// @API CBR GET /v3/{project_id}/operation-logs/{operation_log_id}
func ResourceRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRestoreCreate,
		ReadContext:   resourceRestoreRead,
		UpdateContext: resourceRestoreUpdate,
		DeleteContext: resourceRestoreDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableRestoreParams),

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
				Description: `Specifies the backup ID.`,
			},
			"mappings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the restored mapping relationship.`,
				Elem:        mappingsSchema(),
			},
			"power_on": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the server is powered on after restoration.`,
			},
			"server_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the target VM to be restored.`,
			},
			"volume_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the target disk to be restored.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the resource to be restored.`,
			},
			"details": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Specifies the restoration details.`,
				Elem:        detailsSchema(),
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

func mappingsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the disk backup ID.`,
			},
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the disk to which data is restored.`,
			},
		},
	}
}

func detailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"destination_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the destination path.`,
			},
		},
	}
}

func buildRestoreCreateMappingsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawArray := d.Get("mappings").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rst = append(rst, map[string]interface{}{
			"backup_id": utils.PathSearch("backup_id", v, nil),
			"volume_id": utils.PathSearch("volume_id", v, nil),
		})
	}

	return rst
}

func buildRestoreCreateDetailsBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("details").([]interface{})
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"destination_path": rawMap["destination_path"],
	}
}

func buildRestoreCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	restoreMap := map[string]interface{}{
		"mappings":    buildRestoreCreateMappingsBodyParams(d),
		"power_on":    d.Get("power_on"),
		"server_id":   utils.ValueIgnoreEmpty(d.Get("server_id")),
		"volume_id":   utils.ValueIgnoreEmpty(d.Get("volume_id")),
		"resource_id": utils.ValueIgnoreEmpty(d.Get("resource_id")),
		"details":     buildRestoreCreateDetailsBodyParams(d),
	}

	return map[string]interface{}{
		"restore": restoreMap,
	}
}

func queryRestoreTask(client *golangsdk.ServiceClient, jobID string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{project_id}/operation-logs/{operation_log_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{operation_log_id}", jobID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBR restore task: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func waitingForRestoreTaskSuccess(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	jobID string) error {
	unexpectedStatus := []string{"failed", "timeout"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := queryRestoreTask(client, jobID)
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

func resourceRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v3/{project_id}/backups/{backup_id}/restore"
		product  = "cbr"
		backupID = d.Get("backup_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{backup_id}", backupID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRestoreCreateBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error recovering CBR backup: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("job_id", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error recovering CBR backup: Job ID is not found in API response")
	}

	d.SetId(backupID)
	if err := waitingForRestoreTaskSuccess(ctx, client, d.Timeout(schema.TimeoutCreate), jobID); err != nil {
		return diag.Errorf("error waiting for CBR restore backup (%s) to success: %s", d.Id(), err)
	}

	return resourceRestoreRead(ctx, d, meta)
}

func resourceRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceRestoreUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to restore a CBR backup. Deleting this 
resource will not change the current restore backup result, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
