package coc

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

var diagnosisTaskRetryNonUpdatableParams = []string{"task_id", "instance_id"}

// @API COC POST /v1/diagnosis/tasks/{task_id}/retry
func ResourceDiagnosisTaskRetry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDiagnosisTaskRetryCreate,
		ReadContext:   resourceDiagnosisTaskRetryRead,
		UpdateContext: resourceDiagnosisTaskRetryUpdate,
		DeleteContext: resourceDiagnosisTaskRetryDelete,

		CustomizeDiff: config.FlexibleForceNew(diagnosisTaskRetryNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceDiagnosisTaskRetryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/diagnosis/tasks/{task_id}/retry"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	taskID := d.Get("task_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"instance_id": utils.ValueIgnoreEmpty(d.Get("instance_id")),
		}),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error retrying the COC diagnosis task (%s): %s", taskID, err)
	}

	d.SetId(taskID)

	return resourceDiagnosisTaskRetryRead(ctx, d, meta)
}

func resourceDiagnosisTaskRetryRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDiagnosisTaskRetryUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDiagnosisTaskRetryDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting diagnosis task retry resource is not supported. The diagnosis task retry resource is only" +
		" removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
