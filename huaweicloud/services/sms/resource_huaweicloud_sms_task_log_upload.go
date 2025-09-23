package sms

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

var taskLogUploadNonUpdatableParams = []string{"task_id", "log_bucket"}

// @API SMS POST /v3/tasks/{task_id}/log
func ResourceTaskLogUpload() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskLogUploadCreate,
		ReadContext:   resourceTaskLogUploadRead,
		UpdateContext: resourceTaskLogUploadUpdate,
		DeleteContext: resourceTaskLogUploadDelete,

		CustomizeDiff: config.FlexibleForceNew(taskLogUploadNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_bucket": {
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

func resourceTaskLogUploadCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	taskID := d.Get("task_id").(string)
	logBucket := d.Get("log_bucket").(string)
	err = taskLogUpload(client, taskID, logBucket)
	if err != nil {
		return diag.Errorf("error uploading SMS task log: %s", err)
	}

	d.SetId(taskID)

	return nil
}

func taskLogUpload(client *golangsdk.ServiceClient, taskID, logBucket string) error {
	httpUrl := "v3/tasks/{task_id}/log"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"log_bucket": logBucket,
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return err
	}
	return nil
}

func resourceTaskLogUploadRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaskLogUploadUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaskLogUploadDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting the task log upload resource is not supported. The task log upload resource is only removed" +
		" from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
